package dkron

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"flag"
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
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

const (
	QuerySchedulerRestart = "scheduler:restart"
	QueryRunJob           = "run:job"
	QueryRPCConfig        = "rpc:config"

	// gracefulTimeout controls how long we wait before forcefully terminating
	gracefulTimeout = 3 * time.Second
)

// AgentCommand run server
type AgentCommand struct {
	Ui         cli.Ui
	Version    string
	ShutdownCh <-chan struct{}
	serf       *serf.Serf
	config     *Config
	store      *Store
	eventCh    chan serf.Event
	sched      *Scheduler
}

func (a *AgentCommand) Help() string {
	helpText := `
Usage: dkron agent [options]
	Run dkron agent

Options:

  -bind=0.0.0.0:8946              Address to bind network listeners to.
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
  -keyspace=dkron                 The etcd keyspace to use. A prefix under all data is stored
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

  -debug=false                    Output debug log
`
	return strings.TrimSpace(helpText)
}

// readConfig is responsible for setup of our configuration using
// the command line and any file configs
func (a *AgentCommand) readConfig(args []string) *Config {
	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	cmdFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	cmdFlags.Usage = func() { a.Ui.Output(a.Help()) }
	cmdFlags.String("node", hostname, "node name")
	viper.SetDefault("node_name", cmdFlags.Lookup("node").Value)
	cmdFlags.String("bind", fmt.Sprintf("0.0.0.0:%d", DefaultBindPort), "address to bind listeners to")
	viper.SetDefault("bind_addr", cmdFlags.Lookup("bind").Value)
	cmdFlags.String("http-addr", ":8080", "HTTP address")
	viper.SetDefault("http_addr", cmdFlags.Lookup("http-addr").Value)
	cmdFlags.String("discover", "dkron", "mDNS discovery name")
	viper.SetDefault("discover", cmdFlags.Lookup("discover").Value)
	cmdFlags.String("backend", "etcd", "store backend")
	viper.SetDefault("backend", cmdFlags.Lookup("backend").Value)
	cmdFlags.String("backend-machine", "127.0.0.1:2379", "store backend machines addresses")
	viper.SetDefault("backend_machine", cmdFlags.Lookup("backend-machine").Value)
	cmdFlags.String("profile", "lan", "timing profile to use (lan, wan, local)")
	viper.SetDefault("profile", cmdFlags.Lookup("profile").Value)
	viper.SetDefault("server", cmdFlags.Bool("server", false, "start dkron server"))
	var startJoin []string
	cmdFlags.Var((*AppendSliceValue)(&startJoin), "join", "address of agent to join on startup")
	var tag []string
	cmdFlags.Var((*AppendSliceValue)(&tag), "tag", "tag pair, specified as key=value")
	cmdFlags.String("keyspace", "dkron", "key namespace to use")
	viper.SetDefault("keyspace", cmdFlags.Lookup("keyspace").Value)
	cmdFlags.String("encrypt", "", "encryption key")
	viper.SetDefault("encrypt", cmdFlags.Lookup("encrypt").Value)
	viper.SetDefault("debug", cmdFlags.Bool("debug", false, "output debug log"))
	cmdFlags.String("ui-dir", ".", "directory to serve web UI")
	viper.SetDefault("ui_dir", cmdFlags.Lookup("ui-dir").Value)
	viper.SetDefault("rpc_port", cmdFlags.Int("rpc-port", 6868, "RPC port"))

	// Notifications
	cmdFlags.String("mail-host", "", "notification mail server host")
	viper.SetDefault("mail_host", cmdFlags.Lookup("mail-host").Value)
	cmdFlags.String("mail-port", "", "port to use for the mail server")
	viper.SetDefault("mail_port", cmdFlags.Lookup("mail-port").Value)
	cmdFlags.String("mail-username", "", "username for the mail server")
	viper.SetDefault("mail_username", cmdFlags.Lookup("mail-username").Value)
	cmdFlags.String("mail-password", "", "password of the mail server")
	viper.SetDefault("mail_password", cmdFlags.Lookup("mail-password").Value)
	cmdFlags.String("mail-from", "", "notification emails from address")
	viper.SetDefault("mail_from", cmdFlags.Lookup("mail-from").Value)

	cmdFlags.String("webhook-url", "", "notification webhook url")
	viper.SetDefault("webhook_url", cmdFlags.Lookup("webhook-url").Value)
	cmdFlags.String("webhook-payload", "", "notification webhook payload")
	viper.SetDefault("webhook_payload", cmdFlags.Lookup("webhook-payload").Value)
	webhookHeaders := &AppendSliceValue{}
	cmdFlags.Var(webhookHeaders, "webhook-header", "notification webhook additional header")

	if err := cmdFlags.Parse(args); err != nil {
		log.Fatal(err)
	}

	ut, err := UnmarshalTags(tag)
	if err != nil {
		log.Fatal(err)
	}
	viper.SetDefault("tags", ut)
	viper.SetDefault("join", startJoin)
	viper.SetDefault("webhook_headers", webhookHeaders)

	tags := viper.GetStringMapString("tags")
	server := viper.GetBool("server")
	nodeName := viper.GetString("node_name")

	if server {
		data := []byte(nodeName + fmt.Sprintf("%s", time.Now()))
		tags["key"] = fmt.Sprintf("%x", sha1.Sum(data))
		tags["server"] = "true"
	}

	SetLogLevel(viper.GetBool("debug"))

	return &Config{
		NodeName:        nodeName,
		BindAddr:        viper.GetString("bind_addr"),
		HTTPAddr:        viper.GetString("http_addr"),
		Discover:        viper.GetString("discover"),
		Backend:         viper.GetString("backend"),
		BackendMachines: viper.GetStringSlice("backend_machine"),
		Server:          server,
		Profile:         viper.GetString("profile"),
		StartJoin:       viper.GetStringSlice("join"),
		Tags:            tags,
		Keyspace:        viper.GetString("keyspace"),
		EncryptKey:      viper.GetString("encrypt"),
		UIDir:           viper.GetString("ui_dir"),
		RPCPort:         viper.GetInt("rpc_port"),

		MailHost:     viper.GetString("mail_host"),
		MailPort:     uint16(viper.GetInt("mail_port")),
		MailUsername: viper.GetString("mail_username"),
		MailPassword: viper.GetString("mail_password"),
		MailFrom:     viper.GetString("mail_from"),

		WebhookURL:     viper.GetString("webhook_url"),
		WebhookPayload: viper.GetString("webhook_payload"),
		WebhookHeaders: viper.GetStringSlice("webhook_headers"),
	}
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
	a.config = a.readConfig(args)
	if a.serf = a.setupSerf(); a.serf == nil {
		log.Fatal("agent: Can not setup serf")
	}
	a.join(a.config.StartJoin, true)

	if a.config.Server {
		a.store = NewStore(a.config.Backend, a.config.BackendMachines, a, a.config.Keyspace)
		a.sched = NewScheduler()

		a.ServeHTTP()
		listenRPC(a)

		if a.ElectLeader() {
			a.schedule()
		}
	}
	go a.eventLoop()

	return a.handleSignals()
}

// handleSignals blocks until we get an exit-causing signal
func (a *AgentCommand) handleSignals() int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Wait for a signal
	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	case <-a.ShutdownCh:
		sig = os.Interrupt
	}
	a.Ui.Output(fmt.Sprintf("Caught signal: %v", sig))

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
		if err := a.serf.Leave(); err != nil {
			a.Ui.Error(fmt.Sprintf("Error: %s", err))
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

func (a *AgentCommand) Synopsis() string {
	return "Run dkron"
}

// Leader election routine
func (a *AgentCommand) ElectLeader() bool {
	leader := a.store.GetLeader()

	if leader != nil {
		if !a.serverAlive(string(leader.Key)) {
			log.Debug("agent: Trying to set itself as leader")
			success, err := a.store.TryLeaderSwap(a.config.Tags["key"], leader)
			if err != nil || success == false {
				log.Errorln("agent: Error trying to set itself as leader", err)
				return false
			}
			return true
		} else {
			log.WithFields(logrus.Fields{
				"key": string(leader.Key),
			}).Info("agent: The current leader is active")
		}
	} else {
		log.Debug("agent: Trying to set itself as leader")
		err := a.store.SetLeader(a.config.Tags["key"])
		if err != nil {
			log.Error(err)
		}
		log.WithFields(logrus.Fields{
			"key": a.config.Tags["key"],
		}).Info("agent: Successfully set leader")
		return true
	}

	return false
}

// Checks if the server member identified by key, is alive.
func (a *AgentCommand) serverAlive(key string) bool {
	members := a.serf.Members()
	for _, member := range members {
		if member.Tags["key"] == key && member.Status == serf.StatusAlive {
			return true
		}
	}
	return false
}

// Utility method to check if the node calling the method is the leader.
func (a *AgentCommand) isLeader() bool {
	return a.config.Tags["key"] == string(a.store.GetLeader().Key)
}

// Utility method to get leader nodename
func (a *AgentCommand) leaderMember() (*serf.Member, error) {
	leader := string(a.store.GetLeader().Key)
	for _, member := range a.serf.Members() {
		if key, ok := member.Tags["key"]; ok {
			if key == leader {
				return &member, nil
			}
		}
	}
	return nil, errors.New("No member leader found in member list")
}

func (a *AgentCommand) listServers() []serf.Member {
	members := []serf.Member{}

	for _, member := range a.serf.Members() {
		if key, ok := member.Tags["server"]; ok {
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

			if (e.EventType() == serf.EventMemberFailed || e.EventType() == serf.EventMemberLeave) && a.config.Server {
				failed := e.(serf.MemberEvent)
				for _, member := range failed.Members {
					if member.Tags["key"] == string(a.store.GetLeader().Key) && member.Status != serf.StatusAlive {
						if a.ElectLeader() {
							a.schedule()
						}
					}
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

					var ex Execution
					if err := json.Unmarshal(query.Payload, &ex); err != nil {
						log.WithFields(logrus.Fields{
							"query": QueryRunJob,
						}).Fatal("agent: Error unmarshaling job payload")
					}

					log.WithFields(logrus.Fields{
						"job": ex.JobName,
					}).Info("agent: Starting job")

					ex.StartedAt = time.Now()
					ex.Success = false
					ex.NodeName = a.config.NodeName

					go func() {
						a.invokeJob(&ex)
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
	if a.sched.Started {
		a.sched.Restart(jobs)
	} else {
		a.sched.Start(jobs)
	}
}

func (a *AgentCommand) schedulerRestartQuery(leader *Leader) {
	params := &serf.QueryParam{
		FilterTags: map[string]string{"key": string(leader.Key)},
		RequestAck: true,
	}

	qr, err := a.serf.Query(QuerySchedulerRestart, []byte(""), params)
	if err != nil {
		log.Fatal("agent: Error sending the scheduler reload query", err)
	}
	defer qr.Close()

	ackCh := qr.AckCh()
	respCh := qr.ResponseCh()

	for !qr.Finished() {
		select {
		case ack, ok := <-ackCh:
			if ok {
				log.WithFields(logrus.Fields{
					"from": ack,
				}).Debug("agent: Received ack")
			}
		case resp, ok := <-respCh:
			if ok {
				log.WithFields(logrus.Fields{
					"from":    resp.From,
					"payload": string(resp.Payload),
				}).Debug("agent: Received response")
			}
		}
	}
	log.Debug("agent: Done receiving acks and responses from scheduler reload query")
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

func (a *AgentCommand) RunQuery(job *Job) {
	filterNodes, filterTags, err := a.processFilteredNodes(job)
	if err != nil {
		log.WithFields(logrus.Fields{
			"job": job.Name,
			"err": err.Error(),
		}).Fatal("agent: Error processing filtered nodes")
	}
	log.Debug("agent: Filtered nodes to run: ", filterNodes)
	log.Debug("agent: Filtered tags to run: ", job.Tags)

	params := &serf.QueryParam{
		FilterNodes: filterNodes,
		FilterTags:  filterTags,
		RequestAck:  true,
	}

	ex := Execution{
		JobName: job.Name,
		Group:   time.Now().UnixNano(),
		Job:     job,
	}

	exJson, _ := json.Marshal(ex)
	log.WithFields(logrus.Fields{
		"query":    QueryRunJob,
		"job_name": ex.JobName,
		"json":     string(exJson),
	}).Debug("agent: Sending query")

	qr, err := a.serf.Query(QueryRunJob, exJson, params)
	if err != nil {
		log.WithFields(logrus.Fields{
			"query": QueryRunJob,
			"error": err,
		}).Debug("agent: Sending query error")
	}
	defer qr.Close()

	ackCh := qr.AckCh()
	respCh := qr.ResponseCh()

	for !qr.Finished() {
		select {
		case ack, ok := <-ackCh:
			if ok {
				log.WithFields(logrus.Fields{
					"query": QueryRunJob,
					"from":  ack,
				}).Debug("agent: Received ack")
			}
		case resp, ok := <-respCh:
			if ok {
				log.WithFields(logrus.Fields{
					"query":    QueryRunJob,
					"from":     resp.From,
					"response": string(resp.Payload),
				}).Debug("agent: Received response")

				// Save execution to store
				a.setExecution(resp.Payload)
			}
		}
	}
	log.WithFields(logrus.Fields{
		"query": QueryRunJob,
	}).Debug("agent: Done receiving acks and responses")

}

func (a *AgentCommand) processFilteredNodes(job *Job) ([]string, map[string]string, error) {
	var nodes []string
	tags := job.Tags

	for jtk, jtv := range job.Tags {
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
