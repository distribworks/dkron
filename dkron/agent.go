package dkron

import (
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	metrics "github.com/armon/go-metrics"
	"github.com/docker/leadership"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
	"github.com/mitchellh/cli"
)

const (
	// gracefulTimeout controls how long we wait before forcefully terminating
	gracefulTimeout = 3 * time.Second

	defaultRecoverTime = 10 * time.Second
	defaultLeaderTTL   = 20 * time.Second
)

var (
	expNode = expvar.NewString("node")

	// Error thrown on obtained leader from store is not found in member list
	ErrLeaderNotFound = errors.New("No member leader found in member list")
)

// ProcessorFactory is a function type that creates a new instance
// of a processor.
type ProcessorFactory func() (ExecutionProcessor, error)

// AgentCommand run server
type AgentCommand struct {
	Ui               cli.Ui
	Version          string
	ShutdownCh       <-chan struct{}
	ProcessorPlugins map[string]ExecutionProcessor

	serf      *serf.Serf
	config    *Config
	store     *Store
	eventCh   chan serf.Event
	sched     *Scheduler
	candidate *leadership.Candidate
	ready     bool
}

func (a *AgentCommand) Help() string {
	helpText := `
Usage: dkron agent [options]
	Run dkron agent

Options:

  -bind=0.0.0.0:8946              Address to bind network listeners to.
  -advertise=bind_addr            Address used to advertise to other nodes in the cluster. By default, the bind address is advertised.
  -http-addr=0.0.0.0:8080         Address to bind the UI web server to. Only used when server.
  -discover=cluster               A cluster name used to discovery peers. On
                                  networks that support multicast, this can be used to have
                                  peers join each other without an explicit join.
  -join=addr                      An initial agent to join with. This flag can be
                                  specified multiple times.
  -node=hostname                  Name of this node. Must be unique in the cluster
  -profile=[lan|wan|local]        Profile is used to control the timing profiles used.
                                  The default if not provided is lan.
  -server=false                   This node is running in server mode.
  -tag key=value                  Tag can be specified multiple times to attach multiple
                                  key/value tag pairs to the given node.
  -keyspace=dkron                 The keyspace to use. A prefix under all data is stored
                                  for this instance.
  -backend=[etcd|consul|zk]       Backend storage to use, etcd, consul or zookeeper. The default
                                  is etcd.
  -backend-machine=127.0.0.1:2379 Backend storage servers addresses to connect to. This flag can be
                                  specified multiple times.
  -encrypt                        Key for encrypting network traffic.
                                  Must be a base64-encoded 16-byte key.
  -ui-dir                         Directory from where to serve Web UI
  -rpc-port=6868                  RPC Port used to communicate with clients. Only used when server.
                                  The RPC IP Address will be the same as the bind address.

  -mail-host                      Mail server host address to use for notifications.
  -mail-port                      Mail server port.
  -mail-username                  Mail server username used for authentication.
  -mail-password                  Mail server password to use.
  -mail-from                      From email address to use.

  -webhook-url                    Webhook url to call for notifications.
  -webhook-payload                Body of the POST request to send on webhook call.
  -webhook-header                 Headers to use when calling the webhook URL. Can be specified multiple times.

  -log-level=info                 Log level (debug, info, warn, error, fatal, panic). Default to info.
`
	return strings.TrimSpace(helpText)
}

// setupSerf is used to create the agent we use
func (a *AgentCommand) setupSerf() *serf.Serf {
	config := a.config

	bindIP, bindPort, err := config.AddrParts(config.BindAddr)
	if err != nil {
		a.Ui.Error(fmt.Sprintf("Invalid bind address: %s", err))
		return nil
	}

	// Check if we have an interface
	if iface, _ := config.NetworkInterface(); iface != nil {
		addrs, err := iface.Addrs()
		if err != nil {
			a.Ui.Error(fmt.Sprintf("Failed to get interface addresses: %s", err))
			return nil
		}
		if len(addrs) == 0 {
			a.Ui.Error(fmt.Sprintf("Interface '%s' has no addresses", config.Interface))
			return nil
		}

		// If there is no bind IP, pick an address
		if bindIP == "0.0.0.0" {
			found := false
			for _, ad := range addrs {
				var addrIP net.IP
				if runtime.GOOS == "windows" {
					// Waiting for https://github.com/golang/go/issues/5395 to use IPNet only
					addr, ok := ad.(*net.IPAddr)
					if !ok {
						continue
					}
					addrIP = addr.IP
				} else {
					addr, ok := ad.(*net.IPNet)
					if !ok {
						continue
					}
					addrIP = addr.IP
				}

				// Skip self-assigned IPs
				if addrIP.IsLinkLocalUnicast() {
					continue
				}

				// Found an IP
				found = true
				bindIP = addrIP.String()
				a.Ui.Output(fmt.Sprintf("Using interface '%s' address '%s'",
					config.Interface, bindIP))

				// Update the configuration
				bindAddr := &net.TCPAddr{
					IP:   net.ParseIP(bindIP),
					Port: bindPort,
				}
				config.BindAddr = bindAddr.String()
				break
			}
			if !found {
				a.Ui.Error(fmt.Sprintf("Failed to find usable address for interface '%s'", config.Interface))
				return nil
			}

		} else {
			// If there is a bind IP, ensure it is available
			found := false
			for _, ad := range addrs {
				addr, ok := ad.(*net.IPNet)
				if !ok {
					continue
				}
				if addr.IP.String() == bindIP {
					found = true
					break
				}
			}
			if !found {
				a.Ui.Error(fmt.Sprintf("Interface '%s' has no '%s' address",
					config.Interface, bindIP))
				return nil
			}
		}
	}

	var advertiseIP string
	var advertisePort int
	if config.AdvertiseAddr != "" {
		advertiseIP, advertisePort, err = config.AddrParts(config.AdvertiseAddr)
		if err != nil {
			a.Ui.Error(fmt.Sprintf("Invalid advertise address: %s", err))
			return nil
		}
	}

	encryptKey, err := config.EncryptBytes()
	if err != nil {
		a.Ui.Error(fmt.Sprintf("Invalid encryption key: %s", err))
		return nil
	}

	serfConfig := serf.DefaultConfig()
	switch config.Profile {
	case "lan":
		serfConfig.MemberlistConfig = memberlist.DefaultLANConfig()
	case "wan":
		serfConfig.MemberlistConfig = memberlist.DefaultWANConfig()
	case "local":
		serfConfig.MemberlistConfig = memberlist.DefaultLocalConfig()
	default:
		a.Ui.Error(fmt.Sprintf("Unknown profile: %s", config.Profile))
		return nil
	}

	serfConfig.MemberlistConfig.BindAddr = bindIP
	serfConfig.MemberlistConfig.BindPort = bindPort
	serfConfig.MemberlistConfig.AdvertiseAddr = advertiseIP
	serfConfig.MemberlistConfig.AdvertisePort = advertisePort
	serfConfig.MemberlistConfig.SecretKey = encryptKey
	serfConfig.NodeName = config.NodeName
	serfConfig.Tags = config.Tags
	serfConfig.SnapshotPath = config.SnapshotPath
	serfConfig.CoalescePeriod = 3 * time.Second
	serfConfig.QuiescentPeriod = time.Second
	serfConfig.UserCoalescePeriod = 3 * time.Second
	serfConfig.UserQuiescentPeriod = time.Second
	if config.ReconnectInterval != 0 {
		serfConfig.ReconnectInterval = config.ReconnectInterval
	}
	if config.ReconnectTimeout != 0 {
		serfConfig.ReconnectTimeout = config.ReconnectTimeout
	}
	if config.TombstoneTimeout != 0 {
		serfConfig.TombstoneTimeout = config.TombstoneTimeout
	}
	serfConfig.EnableNameConflictResolution = !config.DisableNameResolution
	if config.KeyringFile != "" {
		serfConfig.KeyringFile = config.KeyringFile
	}
	serfConfig.RejoinAfterLeave = config.RejoinAfterLeave

	// Create a channel to listen for events from Serf
	a.eventCh = make(chan serf.Event, 64)
	serfConfig.EventCh = a.eventCh

	// Start Serf
	a.Ui.Output("Starting Dkron agent...")
	log.Info("agent: Dkron agent starting")

	serfConfig.LogOutput = ioutil.Discard
	serfConfig.MemberlistConfig.LogOutput = ioutil.Discard

	// Create serf first
	serf, err := serf.Create(serfConfig)
	if err != nil {
		a.Ui.Error(err.Error())
		log.Error(err)
		return nil
	}

	return serf
}

// UnmarshalTags is a utility function which takes a slice of strings in
// key=value format and returns them as a tag mapping.
func UnmarshalTags(tags []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, tag := range tags {
		parts := strings.SplitN(tag, "=", 2)
		if len(parts) != 2 || len(parts[0]) == 0 {
			return nil, fmt.Errorf("Invalid tag: '%s'", tag)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}

func (a *AgentCommand) Run(args []string) int {
	a.config = NewConfig(args, a)
	if a.serf = a.setupSerf(); a.serf == nil {
		log.Fatal("agent: Can not setup serf")
	}
	a.join(a.config.StartJoin, true)

	if i := initMetrics(a); i != 0 {
		return i
	}

	// Expose the node name
	expNode.Set(a.config.NodeName)

	if a.config.Server {
		a.store = NewStore(a.config.Backend, a.config.BackendMachines, a, a.config.Keyspace)
		a.sched = NewScheduler()

		a.ServeHTTP()
		listenRPC(a)
		a.participate()
	}
	go a.eventLoop()
	a.ready = true
	return a.handleSignals()
}

// handleSignals blocks until we get an exit-causing signal
func (a *AgentCommand) handleSignals() int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

WAIT:
	// Wait for a signal
	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	case <-a.ShutdownCh:
		sig = os.Interrupt
	}
	a.Ui.Output(fmt.Sprintf("Caught signal: %v", sig))

	// Check if this is a SIGHUP
	if sig == syscall.SIGHUP {
		a.handleReload()
		goto WAIT
	}

	// Check if we should do a graceful leave
	graceful := false
	if sig == syscall.SIGTERM || sig == os.Interrupt {
		graceful = true
	}

	// Bail fast if not doing a graceful leave
	if !graceful {
		return 1
	}

	// Attempt a graceful leave
	gracefulCh := make(chan struct{})
	a.Ui.Output("Gracefully shutting down agent...")
	log.Info("agent: Gracefully shutting down agent...")
	go func() {
		// If we're exiting a server
		if a.config.Server {
			// Stop running for leader election
			a.candidate.Stop()
		}
		if err := a.serf.Leave(); err != nil {
			a.Ui.Error(fmt.Sprintf("Error: %s", err))
			log.Error(fmt.Sprintf("Error: %s", err))
			return
		}
		close(gracefulCh)
	}()

	// Wait for leave or another signal
	select {
	case <-signalCh:
		return 1
	case <-time.After(gracefulTimeout):
		return 1
	case <-gracefulCh:
		return 0
	}
}

// handleReload is invoked when we should reload our configs, e.g. SIGHUP
func (a *AgentCommand) handleReload() {
	a.Ui.Output("Reloading configuration...")
	newConf := ReadConfig(a)
	if newConf == nil {
		a.Ui.Error(fmt.Sprintf("Failed to reload configs"))
		return
	} else {
		a.config = newConf
	}

	// Reset serf tags
	a.serf.SetTags(a.config.Tags)
	//Config reloading will also reload Notification settings
}

func (a *AgentCommand) Synopsis() string {
	return "Run dkron"
}

func (a *AgentCommand) participate() {
	a.candidate = leadership.NewCandidate(a.store.Client, a.store.LeaderKey(), a.config.NodeName, defaultLeaderTTL)

	go func() {
		for {
			a.runForElection()
			// retry
			time.Sleep(defaultRecoverTime)
		}
	}()
}

// Leader election routine
func (a *AgentCommand) runForElection() {
	log.Info("agent: Running for election")
	defer metrics.MeasureSince([]string{"agent", "runForElection"}, time.Now())
	electedCh, errCh := a.candidate.RunForElection()

	for {
		select {
		case isElected := <-electedCh:
			if isElected {
				log.Info("agent: Cluster leadership acquired")
				metrics.IncrCounter([]string{"agent", "leadership_acquired"}, 1)
				// If this server is elected as the leader, start the scheduler
				a.schedule()
			} else {
				log.Info("agent: Cluster leadership lost")
				metrics.IncrCounter([]string{"agent", "leadership_lost"}, 1)
				// Always stop the schedule of this server to prevent multiple servers with the scheduler on
				a.sched.Stop()
			}

		case err := <-errCh:
			log.WithError(err).Debug("Leader election failed, channel is probably closed")
			metrics.IncrCounter([]string{"agent", "election", "failure"}, 1)
			// Always stop the schedule of this server to prevent multiple servers with the scheduler on
			a.sched.Stop()
			return
		}
	}
}

// Utility method to get leader nodename
func (a *AgentCommand) leaderMember() (*serf.Member, error) {
	leaderName := a.store.GetLeader()
	for _, member := range a.serf.Members() {
		if member.Name == string(leaderName) {
			return &member, nil
		}
	}
	return nil, ErrLeaderNotFound
}

func (a *AgentCommand) listServers() []serf.Member {
	members := []serf.Member{}

	for _, member := range a.serf.Members() {
		if key, ok := member.Tags["dkron_server"]; ok {
			if key == "true" && member.Status == serf.StatusAlive {
				members = append(members, member)
			}
		}
	}
	return members
}

// Listens to events from Serf and handle the event.
func (a *AgentCommand) eventLoop() {
	serfShutdownCh := a.serf.ShutdownCh()
	log.Info("agent: Listen for events")
	for {
		select {
		case e := <-a.eventCh:
			log.WithFields(logrus.Fields{
				"event": e.String(),
			}).Debug("agent: Received event")
			metrics.AddSample([]string{"agent", "event_received", e.String()}, 1)

			// Log all member events
			if failed, ok := e.(serf.MemberEvent); ok {
				for _, member := range failed.Members {
					log.WithFields(logrus.Fields{
						"node":   a.config.NodeName,
						"member": member.Name,
						"event":  e.EventType(),
					}).Debug("agent: Member event")
				}
			}

			if e.EventType() == serf.EventQuery {
				query := e.(*serf.Query)

				if query.Name == QuerySchedulerRestart && a.config.Server {
					a.schedule()
				}

				if query.Name == QueryRunJob {
					log.WithFields(logrus.Fields{
						"query":   query.Name,
						"payload": string(query.Payload),
						"at":      query.LTime,
					}).Debug("agent: Running job")

					var rqp RunQueryParam
					if err := json.Unmarshal(query.Payload, &rqp); err != nil {
						log.WithField("query", QueryRunJob).Fatal("agent: Error unmarshaling query payload")
					}

					log.WithFields(logrus.Fields{
						"job": rqp.Execution.JobName,
					}).Info("agent: Starting job")

					rpcc := RPCClient{ServerAddr: rqp.RPCAddr}
					job, err := rpcc.GetJob(rqp.Execution.JobName)
					if err != nil {
						log.WithError(err).Error("agent: Error on rpc.GetJob call")
					}
					log.WithField("command", job.Command).Debug("agent: GetJob by RPC")

					ex := rqp.Execution
					ex.StartedAt = time.Now()
					ex.NodeName = a.config.NodeName

					go func() {
						if err := a.invokeJob(job, ex); err != nil {
							log.WithError(err).Error("agent: Error invoking job command")
						}
					}()

					exJson, _ := json.Marshal(ex)
					query.Respond(exJson)
				}

				if query.Name == QueryRPCConfig && a.config.Server {
					log.WithFields(logrus.Fields{
						"query":   query.Name,
						"payload": string(query.Payload),
						"at":      query.LTime,
					}).Debug("agent: RPC Config requested")

					query.Respond([]byte(a.getRPCAddr()))
				}
			}

		case <-serfShutdownCh:
			log.Warn("agent: Serf shutdown detected, quitting")
			return
		}
	}
}

// Start or restart scheduler
func (a *AgentCommand) schedule() {
	log.Debug("agent: Restarting scheduler")
	jobs, err := a.store.GetJobs()
	if err != nil {
		log.Fatal(err)
	}
	a.sched.Restart(jobs)
}

// Join asks the Serf instance to join. See the Serf.Join function.
func (a *AgentCommand) join(addrs []string, replay bool) (n int, err error) {
	log.Infof("agent: joining: %v replay: %v", addrs, replay)
	ignoreOld := !replay
	n, err = a.serf.Join(addrs, ignoreOld)
	if n > 0 {
		log.Infof("agent: joined: %d nodes", n)
	}
	if err != nil {
		log.Warnf("agent: error joining: %v", err)
	}
	return
}

func (a *AgentCommand) processFilteredNodes(job *Job) ([]string, map[string]string, error) {
	var nodes []string
	tags := make(map[string]string)

	// Actually copy the map
	for key, val := range job.Tags {
		tags[key] = val
	}

	for jtk, jtv := range tags {
		var tc []string
		if tc = strings.Split(jtv, ":"); len(tc) == 2 {
			tv := tc[0]

			// Set original tag to clean tag
			tags[jtk] = tv

			count, err := strconv.Atoi(tc[1])
			if err != nil {
				return nil, nil, err
			}

			for _, member := range a.serf.Members() {
				if member.Status == serf.StatusAlive {
					for mtk, mtv := range member.Tags {
						if mtk == jtk && mtv == tv {
							if len(nodes) < count {
								nodes = append(nodes, member.Name)
							}
						}
					}
				}
			}
		}
	}

	return nodes, tags, nil
}

func (a *AgentCommand) setExecution(payload []byte) *Execution {
	var ex Execution
	if err := json.Unmarshal(payload, &ex); err != nil {
		log.Fatal(err)
	}

	// Save the new execution to store
	if _, err := a.store.SetExecution(&ex); err != nil {
		log.Fatal(err)
	}

	return &ex
}

// This function is called when a client request the RPCAddress
// of the current member.
func (a *AgentCommand) getRPCAddr() string {
	bindIp := a.serf.LocalMember().Addr

	return fmt.Sprintf("%s:%d", bindIp, a.config.RPCPort)
}
