package dkron

import (
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"io/ioutil"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/abronan/leadership"
	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
)

const (
	defaultRecoverTime = 10 * time.Second
)

var (
	expNode = expvar.NewString("node")

	// ErrLeaderNotFound is returned when obtained leader from store is not found in member list
	ErrLeaderNotFound = errors.New("No member leader found in member list")
	ErrNoRPCAddress   = errors.New("No RPC address tag found in server")

	defaultLeaderTTL = 20 * time.Second
)

type Agent struct {
	ProcessorPlugins map[string]ExecutionProcessor
	ExecutorPlugins  map[string]Executor
	HTTPTransport    Transport
	Store            *Store

	serf      *serf.Serf
	config    *Config
	eventCh   chan serf.Event
	sched     *Scheduler
	candidate *leadership.Candidate
	ready     bool
}

// ProcessorFactory is a function type that creates a new instance
// of a processor.
type ProcessorFactory func() (ExecutionProcessor, error)

type Plugins struct {
	Processors map[string]ExecutionProcessor
	Executors  map[string]Executor
}

func NewAgent(config *Config, plugins *Plugins) *Agent {
	a := &Agent{config: config}

	if plugins != nil {
		a.ProcessorPlugins = plugins.Processors
		a.ExecutorPlugins = plugins.Executors
	}

	return a
}

func (a *Agent) Start() error {
	s, err := a.setupSerf()
	if err != nil {
		return fmt.Errorf("agent: Can not setup serf, %s", err)
	}
	a.serf = s
	a.join(a.config.StartJoin, true)

	if err := initMetrics(a); err != nil {
		log.Fatal("agent: Can not setup metrics")
	}

	// Expose the node name
	expNode.Set(a.config.NodeName)

	if a.config.Server {
		a.StartServer()
	}
	go a.eventLoop()
	a.ready = true

	return nil
}

func (a *Agent) Stop() error {
	a.candidate.Stop()

	if err := a.serf.Leave(); err != nil {
		return err
	}

	return nil
}

// setupSerf is used to create the agent we use
func (a *Agent) setupSerf() (*serf.Serf, error) {
	config := a.config

	bindIP, bindPort, err := config.AddrParts(config.BindAddr)
	if err != nil {
		return nil, fmt.Errorf("Invalid bind address: %s", err)
	}

	// Check if we have an interface
	if iface, _ := config.NetworkInterface(); iface != nil {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("Failed to get interface addresses: %s", err)
		}
		if len(addrs) == 0 {
			return nil, fmt.Errorf("Interface '%s' has no addresses", config.Interface)
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
				log.Infof("Using interface '%s' address '%s'", config.Interface, bindIP)

				// Update the configuration
				bindAddr := &net.TCPAddr{
					IP:   net.ParseIP(bindIP),
					Port: bindPort,
				}
				config.BindAddr = bindAddr.String()
				break
			}
			if !found {
				return nil, fmt.Errorf("Failed to find usable address for interface '%s'", config.Interface)
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
				return nil, fmt.Errorf("Interface '%s' has no '%s' address",
					config.Interface, bindIP)
			}
		}
	}

	var advertiseIP string
	var advertisePort int
	if config.AdvertiseAddr != "" {
		advertiseIP, advertisePort, err = config.AddrParts(config.AdvertiseAddr)
		if err != nil {
			return nil, fmt.Errorf("Invalid advertise address: %s", err)
		}
	}
	//Ues the value of "RPCPort" if AdvertiseRPCPort has not been set
	if config.AdvertiseRPCPort <= 0 {
		config.AdvertiseRPCPort = config.RPCPort
	}

	encryptKey, err := config.EncryptBytes()
	if err != nil {
		return nil, fmt.Errorf("Invalid encryption key: %s", err)
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
		return nil, fmt.Errorf("Unknown profile: %s", config.Profile)
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
	log.Info("agent: Dkron agent starting")

	serfConfig.LogOutput = ioutil.Discard
	serfConfig.MemberlistConfig.LogOutput = ioutil.Discard

	// Create serf first
	serf, err := serf.Create(serfConfig)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return serf, nil
}

// Config returns the agent's config.
func (a *Agent) Config() *Config {
	return a.config
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

func (a *Agent) StartServer() {
	if a.Store == nil {
		a.Store = NewStore(a.config.Backend, a.config.BackendMachines, a, a.config.Keyspace, nil)
	}

	a.sched = NewScheduler()

	if a.HTTPTransport == nil {
		a.HTTPTransport = NewTransport(a)
	}
	a.HTTPTransport.ServeHTTP()
	listenRPC(a)

	tags := a.serf.LocalMember().Tags
	tags["dkron_rpc_addr"] = a.getRPCAddr()
	if err := a.SetTags(tags); err != nil {
		log.WithError(err).Fatal("agent: Error setting RPC config tags")
	}
	a.participate()
}

func (a *Agent) participate() {
	a.candidate = leadership.NewCandidate(a.Store.Client, a.Store.LeaderKey(), a.config.NodeName, defaultLeaderTTL)

	go func() {
		for {
			a.runForElection()
			// retry
			time.Sleep(defaultRecoverTime)
		}
	}()
}

// Leader election routine
func (a *Agent) runForElection() {
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
func (a *Agent) leaderMember() (*serf.Member, error) {
	leaderName := a.Store.GetLeader()
	for _, member := range a.serf.Members() {
		if member.Name == string(leaderName) {
			return &member, nil
		}
	}
	return nil, ErrLeaderNotFound
}

func (a *Agent) listServers() []serf.Member {
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

// GetBindIP returns the IP address that the agent is bound to.
// This could be different than the originally configured address.
func (a *Agent) GetBindIP() (string, error) {
	bindIP, _, err := a.config.AddrParts(a.config.BindAddr)
	return bindIP, err
}

// Listens to events from Serf and handle the event.
func (a *Agent) eventLoop() {
	serfShutdownCh := a.serf.ShutdownCh()
	log.Info("agent: Listen for events")
	for {
		select {
		case e := <-a.eventCh:
			log.WithFields(logrus.Fields{
				"event": e.String(),
			}).Debug("agent: Received event")
			metrics.IncrCounter([]string{"agent", "event_received", e.String()}, 1)

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

					exJSON, _ := json.Marshal(ex)
					query.Respond(exJSON)
				}

				if query.Name == QueryRPCConfig && a.config.Server {
					log.WithFields(logrus.Fields{
						"query":   query.Name,
						"payload": string(query.Payload),
						"at":      query.LTime,
					}).Debug("agent: RPC Config requested")

					err := query.Respond([]byte(a.getRPCAddr()))
					if err != nil {
						log.WithError(err).Error("agent: query.Respond")
					}
				}
			}

		case <-serfShutdownCh:
			log.Warn("agent: Serf shutdown detected, quitting")
			return
		}
	}
}

// Start or restart scheduler
func (a *Agent) schedule() {
	log.Debug("agent: Restarting scheduler")
	jobs, err := a.Store.GetJobs()
	if err != nil {
		log.Fatal(err)
	}
	a.sched.Restart(jobs)
}

// Join asks the Serf instance to join. See the Serf.Join function.
func (a *Agent) join(addrs []string, replay bool) (n int, err error) {
	log.Infof("agent: joining: %v replay: %v", addrs, replay)
	n, err = a.serf.Join(addrs, !replay)
	if n > 0 {
		log.Infof("agent: joined: %d nodes", n)
	}
	if err != nil {
		log.Warnf("agent: error joining: %v", err)
	}
	return
}

func (a *Agent) processFilteredNodes(job *Job) ([]string, map[string]string, error) {
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

func (a *Agent) setExecution(payload []byte) *Execution {
	var ex Execution
	if err := json.Unmarshal(payload, &ex); err != nil {
		log.Fatal(err)
	}

	// Save the new execution to store
	if _, err := a.Store.SetExecution(&ex); err != nil {
		log.Fatal(err)
	}

	return &ex
}

// This function is called when a client request the RPCAddress
// of the current member.
// in marathon, it would return the host's IP and advertise RPC port
func (a *Agent) getRPCAddr() string {
	bindIP := a.serf.LocalMember().Addr

	return fmt.Sprintf("%s:%d", bindIP, a.config.AdvertiseRPCPort)
}

func (a *Agent) Leave() error {
	return a.serf.Leave()
}

func (a *Agent) SetTags(tags map[string]string) error {
	return a.serf.SetTags(tags)
}
