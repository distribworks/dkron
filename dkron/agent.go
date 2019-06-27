package dkron

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
)

const (
	raftTimeout    = 10 * time.Second
	rescheduleTime = 2 * time.Second
	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512
	minRaftProtocol  = 3
)

var (
	expNode = expvar.NewString("node")

	// ErrLeaderNotFound is returned when obtained leader is not found in member list
	ErrLeaderNotFound = errors.New("No member leader found in member list")

	defaultLeaderTTL = 20 * time.Second

	runningExecutions sync.Map
)

// Agent is the main struct that represents a dkron agent
type Agent struct {
	ProcessorPlugins map[string]ExecutionProcessor
	ExecutorPlugins  map[string]Executor
	HTTPTransport    Transport
	Store            Storage
	GRPCServer       DkronGRPCServer
	GRPCClient       DkronGRPCClient
	TLSConfig        *tls.Config

	// Set a global peer updater func
	PeerUpdaterFunc func(...string)

	serf       *serf.Serf
	config     *Config
	eventCh    chan serf.Event
	sched      *Scheduler
	ready      bool
	shutdownCh chan struct{}

	// The raft instance is used among Dkron nodes within the
	// region to protect operations that require strong consistency
	leaderCh <-chan bool
	raft     *raft.Raft
	// raftLayer provides network layering of the raft RPC along with
	// the Dkron gRPC transport layer.
	raftLayer     *RaftLayer
	raftStore     *raftboltdb.BoltStore
	raftInmem     *raft.InmemStore
	raftTransport *raft.NetworkTransport

	// reconcileCh is used to pass events from the serf handler
	// into the leader manager. Mostly used to handle when servers
	// join/leave from the region.
	reconcileCh chan serf.Member

	// peers is used to track the known Dkron servers. This is
	// used for region forwarding and clustering.
	peers      map[string][]*serverParts
	localPeers map[raft.ServerAddress]*serverParts
	peerLock   sync.RWMutex
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
	a.join(a.config.StartJoin, false)

	if err := initMetrics(a); err != nil {
		log.Fatal("agent: Can not setup metrics")
	}

	// Expose the node name
	expNode.Set(a.config.NodeName)

	if a.config.Server {
		a.StartServer()
	}

	if a.GRPCClient == nil {
		a.GRPCClient = NewGRPCClient(nil, a)
	}

	tags := a.serf.LocalMember().Tags
	if a.config.Server {
		tags["rpc_addr"] = a.getRPCAddr() // Address that clients will use to RPC to servers
		tags["port"] = strconv.Itoa(a.config.AdvertiseRPCPort)
	}
	a.serf.SetTags(tags)

	go a.eventLoop()
	a.ready = true

	return nil
}

// Stop stops an agent, if the agent is a server and is running for election
// stop running for election, if this server was the leader
// this will force the cluster to elect a new leader and start a new scheduler.
// If this is a server and has the scheduler started stop it, ignoring if this server
// was participating in leader election or not (local storage).
// Then actually leave the cluster.
func (a *Agent) Stop() error {
	log.Info("agent: Called member stop, now stopping")

	if a.config.Server {
		a.raft.Shutdown()
		a.Store.Shutdown()
	}

	if a.config.Server && a.sched.Started {
		a.sched.Stop()
	}

	if err := a.serf.Leave(); err != nil {
		return err
	}

	if err := a.serf.Shutdown(); err != nil {
		return err
	}

	return nil
}

func (a *Agent) setupRaft() error {
	if a.config.BootstrapExpect > 0 {
		if a.config.BootstrapExpect == 1 {
			a.config.Bootstrap = true
		}
	}

	logger := ioutil.Discard
	if log.Logger.Level == logrus.DebugLevel {
		logger = log.Logger.Writer()
	}

	transport := raft.NewNetworkTransport(a.raftLayer, 3, raftTimeout, logger)
	a.raftTransport = transport

	config := raft.DefaultConfig()
	config.LogOutput = logger

	// Build an all in-memory setup for dev mode, otherwise prepare a full
	// disk-based setup.
	var logStore raft.LogStore
	var stableStore raft.StableStore
	var snapshots raft.SnapshotStore
	if a.config.DevMode {
		store := raft.NewInmemStore()
		a.raftInmem = store
		stableStore = store
		logStore = store
		snapshots = raft.NewDiscardSnapshotStore()
	} else {
		var err error
		// Create the snapshot store. This allows the Raft to truncate the log to
		// mitigate the issue of having an unbounded replicated log.
		snapshots, err = raft.NewFileSnapshotStore(filepath.Join(a.config.DataDir, "raft"), 3, logger)
		if err != nil {
			return fmt.Errorf("file snapshot store: %s", err)
		}

		// Create the BoltDB backend
		s, err := raftboltdb.NewBoltStore(filepath.Join(a.config.DataDir, "raft", "raft.db"))
		if err != nil {
			return fmt.Errorf("error creating new badger store: %s", err)
		}
		a.raftStore = s
		stableStore = s

		// Wrap the store in a LogCache to improve performance
		cacheStore, err := raft.NewLogCache(raftLogCacheSize, s)
		if err != nil {
			s.Close()
			return err
		}
		logStore = cacheStore
	}

	config.LocalID = raft.ServerID(a.config.NodeName)

	// If we are in bootstrap or dev mode and the state is clean then we can
	// bootstrap now.
	if a.config.Bootstrap || a.config.DevMode {
		hasState, err := raft.HasExistingState(logStore, stableStore, snapshots)
		if err != nil {
			return err
		}
		if !hasState {
			configuration := raft.Configuration{
				Servers: []raft.Server{
					{
						ID:      config.LocalID,
						Address: transport.LocalAddr(),
					},
				},
			}
			if err := raft.BootstrapCluster(config, logStore, stableStore, snapshots, transport, configuration); err != nil {
				return err
			}
		}
	}

	// Instantiate the Raft systems. The second parameter is a finite state machine
	// which stores the actual kv pairs and is operated upon through Apply().
	fsm := NewFSM(a.Store)
	rft, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	a.leaderCh = rft.LeaderCh()
	a.raft = rft
	return nil
}

// setupSerf is used to create the agent we use
func (a *Agent) setupSerf() (*serf.Serf, error) {
	config := a.config

	// Init peer list
	a.localPeers = make(map[raft.ServerAddress]*serverParts)
	a.peers = make(map[string][]*serverParts)

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
	//Use the value of "RPCPort" if AdvertiseRPCPort has not been set
	if config.AdvertiseRPCPort <= 0 {
		config.AdvertiseRPCPort = config.RPCPort
	}

	encryptKey, err := config.EncryptBytes()
	if err != nil {
		return nil, fmt.Errorf("Invalid encryption key: %s", err)
	}

	serfConfig := serf.DefaultConfig()
	serfConfig.Init()

	serfConfig.Tags = a.config.Tags
	serfConfig.Tags["role"] = "dkron"
	serfConfig.Tags["dc"] = a.config.Datacenter
	serfConfig.Tags["region"] = a.config.Region
	serfConfig.Tags["version"] = Version
	if a.config.Server {
		serfConfig.Tags["server"] = strconv.FormatBool(a.config.Server)
	}
	if a.config.Bootstrap {
		serfConfig.Tags["bootstrap"] = "1"
	}
	if a.config.BootstrapExpect != 0 {
		serfConfig.Tags["expect"] = fmt.Sprintf("%d", a.config.BootstrapExpect)
	}

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

	if log.Logger.Level == logrus.DebugLevel {
		serfConfig.LogOutput = log.Logger.Writer()
		serfConfig.MemberlistConfig.LogOutput = log.Logger.Writer()
	} else {
		serfConfig.LogOutput = ioutil.Discard
		serfConfig.MemberlistConfig.LogOutput = ioutil.Discard
	}

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

// SetConfig sets the agent's config.
func (a *Agent) SetConfig(c *Config) {
	a.config = c
}

// StartServer launch a new dkron server process
func (a *Agent) StartServer() {
	if a.Store == nil {
		dirExists, err := exists(a.config.DataDir)
		if err != nil {
			log.WithError(err).WithField("dir", a.config.DataDir).Fatal("Invalid Dir")
		}
		if !dirExists {
			// Try to create the directory
			err := os.Mkdir(a.config.DataDir, 0700)
			if err != nil {
				log.WithError(err).WithField("dir", a.config.DataDir).Fatal("Error Creating Dir")
			}
		}
		s, err := NewStore(a, filepath.Join(a.config.DataDir, a.config.NodeName))
		if err != nil {
			log.WithError(err).Fatal("dkron: Error initializing store")
		}
		a.Store = s
	}

	a.sched = NewScheduler()

	if a.HTTPTransport == nil {
		a.HTTPTransport = NewTransport(a)
	}
	a.HTTPTransport.ServeHTTP()

	// Create a listener at the desired port.
	// TODO Fix get address
	addr := fmt.Sprintf("%s:%d", a.serf.LocalMember().Addr, a.config.RPCPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux object.
	tcpm := cmux.New(l)
	var grpcl, raftl net.Listener
	var tlsm cmux.CMux

	// If RaftLayer brings TLS config listen to TLS
	if a.TLSConfig != nil {
		// Create a RaftLayer with TLS
		a.raftLayer = NewTLSRaftLayer(a.TLSConfig)

		// Match any connection to the recursive mux
		tlsl := tcpm.Match(cmux.Any())
		tlsl = tls.NewListener(tlsl, a.TLSConfig)

		// Declare sub cMUX for TLS
		tlsm = cmux.New(tlsl)

		// Declare the match for TLS gRPC
		grpcl = tlsm.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))

		// Declare the match for TLS raft RPC
		raftl = tlsm.Match(cmux.Any())

		go func() {
			if err := tlsm.Serve(); err != nil {
				log.Fatal(err)
			}
		}()
	} else {
		// Declare a plain RaftLayer
		a.raftLayer = NewRaftLayer()

		// Declare the match for gRPC
		grpcl = tcpm.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))

		// Declare the match for raft RPC
		raftl = tcpm.Match(cmux.Any())
	}

	if a.GRPCServer == nil {
		a.GRPCServer = NewGRPCServer(a)
	}

	if err := a.GRPCServer.Serve(grpcl); err != nil {
		log.WithError(err).Fatal("agent: RPC server failed to start")
	}

	if err := a.raftLayer.Open(raftl); err != nil {
		log.Fatal(err)
	}

	if err := a.setupRaft(); err != nil {
		log.WithError(err).Fatal("agent: Raft layer failed to start")
	}

	// Start serving everything
	go func() {
		if err := tcpm.Serve(); err != nil {
			log.Fatal(err)
		}
	}()
	go a.monitorLeadership()
}

// SchedulerRestart Dispatch a SchedulerRestartQuery to the cluster but
// after a timeout to actually throtle subsequent calls
func (a *Agent) SchedulerRestart() {
	if rescheduleThrotle == nil {
		rescheduleThrotle = time.AfterFunc(rescheduleTime, func() {
			a.schedule()
		})
	} else {
		rescheduleThrotle.Reset(rescheduleTime)
	}
}

// Utility method to get leader nodename
func (a *Agent) leaderMember() (*serf.Member, error) {
	l := a.raft.Leader()
	for _, member := range a.serf.Members() {
		if member.Tags["rpc_addr"] == string(l) {
			return &member, nil
		}
	}
	return nil, ErrLeaderNotFound
}

// IsLeader checks if this server is the cluster leader
func (a *Agent) IsLeader() bool {
	return a.raft.State() == raft.Leader
}

// ListServers returns the list of server members
func (a *Agent) ListServers() (members []*serverParts) {
	ok, lm := isServer(a.serf.LocalMember())
	if !ok {
		return nil
	}

	for _, member := range a.serf.Members() {
		ok, parts := isServer(member)
		if !ok || member.Status != serf.StatusAlive {
			continue
		}
		if lm.Region == parts.Region {
			members = append(members, parts)
		}
	}
	return members
}

// LocalMember return the local serf member
func (a *Agent) LocalMember() serf.Member {
	return a.serf.LocalMember()
}

// GetBindIP returns the IP address that the agent is bound to.
// This could be different than the originally configured address.
func (a *Agent) GetBindIP() (string, error) {
	bindIP, _, err := a.config.AddrParts(a.config.BindAddr)
	return bindIP, err
}

// GetPeers returns a list of the current serf servers peers addresses
func (a *Agent) GetPeers() (peers []string) {
	ps := a.ListServers()
	for _, p := range ps {
		peers = append(peers, p.RPCAddr.String())
	}
	return
}

// Listens to events from Serf and handle the event.
func (a *Agent) eventLoop() {
	serfShutdownCh := a.serf.ShutdownCh()
	log.Info("agent: Listen for events")
	for {
		select {
		case e := <-a.eventCh:
			log.WithField("event", e.String()).Info("agent: Received event")
			metrics.IncrCounter([]string{"agent", "event_received", e.String()}, 1)

			// Log all member events
			if me, ok := e.(serf.MemberEvent); ok {
				for _, member := range me.Members {
					log.WithFields(logrus.Fields{
						"node":   a.config.NodeName,
						"member": member.Name,
						"event":  e.EventType(),
					}).Debug("agent: Member event")
				}

				// serfEventHandler is used to handle events from the serf cluster
				switch e.EventType() {
				case serf.EventMemberJoin:
					a.nodeJoin(me)
					a.localMemberEvent(me)
				case serf.EventMemberLeave, serf.EventMemberFailed:
					a.nodeFailed(me)
					a.localMemberEvent(me)
				case serf.EventMemberReap:
					a.localMemberEvent(me)
				case serf.EventMemberUpdate, serf.EventUser, serf.EventQuery: // Ignore
				default:
					log.WithField("event", e.String()).Warn("agent: Unhandled serf event")
				}

				//In case of member event update peer list
				if a.PeerUpdaterFunc != nil {
					a.PeerUpdaterFunc(a.GetPeers()...)
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

					// There are two error types to handle here:
					// Key not found when the job is removed from store
					// Dial tcp error
					// In case of deleted job or other error, we should report and break the flow.
					// On dial error we should retry with a limit.
					i := 0
				RetryGetJob:
					job, err := a.GRPCClient.CallGetJob(rqp.RPCAddr, rqp.Execution.JobName)
					if err != nil {
						if err == ErrRPCDialing {
							if i < 10 {
								i++
								goto RetryGetJob
							}
							log.WithError(err).Fatal("agent: A working RPC connection to a Dkron server must exists.")
						}
						log.WithError(err).Error("agent: Error on rpc.GetJob call")
						continue
					}
					log.WithField("job", job.Name).Debug("agent: GetJob by RPC")

					ex := rqp.Execution
					ex.StartedAt = time.Now()
					ex.NodeName = a.config.NodeName

					go func() {
						if err := a.invokeJob(job, ex); err != nil {
							log.WithError(err).Error("agent: Error invoking job")
						}
					}()

					exJSON, _ := json.Marshal(ex)
					query.Respond(exJSON)
				}

				if query.Name == QueryExecutionDone {
					group := string(query.Payload)

					log.WithFields(logrus.Fields{
						"query":   query.Name,
						"payload": group,
						"at":      query.LTime,
					}).Debug("agent: Execution done requested")

					// Find if the indicated execution is done processing
					var err error
					if _, ok := runningExecutions.Load(group); ok {
						log.WithField("group", group).Debug("agent: Execution is still running")
						err = query.Respond([]byte("false"))
					} else {
						log.WithField("group", group).Debug("agent: Execution is not running")
						err = query.Respond([]byte("true"))
					}
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
	log.Info("agent: Restarting scheduler")
	jobs, err := a.Store.GetJobs(nil)
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

	cmd, err := Encode(SetExecutionType, ex.ToProto())
	if err != nil {
		log.WithError(err).Fatal("agent: encode error in setExecution")
		return nil
	}
	af := a.raft.Apply(cmd, raftTimeout)
	if err := af.Error(); err != nil {
		log.WithError(err).Fatal("agent: error applying SetExecutionType")
		return nil
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

// RefreshJobStatus asks the nodes their progress on an execution
func (a *Agent) RefreshJobStatus(jobName string) {
	var group string

	execs, _ := a.Store.GetLastExecutionGroup(jobName)
	nodes := []string{}

	unfinishedExecutions := []*Execution{}
	for _, ex := range execs {
		if ex.FinishedAt.IsZero() {
			unfinishedExecutions = append(unfinishedExecutions, ex)
		}
	}

	for _, ex := range unfinishedExecutions {
		// Ignore executions that we know are finished
		log.WithFields(logrus.Fields{
			"member":        ex.NodeName,
			"execution_key": ex.Key(),
		}).Info("agent: Asking member for pending execution")

		nodes = append(nodes, ex.NodeName)
		group = strconv.FormatInt(ex.Group, 10)
		log.WithField("group", group).Debug("agent: Pending execution group")
	}

	// If there is pending executions to finish ask if they are really pending.
	if len(nodes) > 0 && group != "" {
		statuses := a.executionDoneQuery(nodes, group)

		log.WithFields(logrus.Fields{
			"statuses": statuses,
		}).Debug("agent: Received pending executions response")

		for _, ex := range unfinishedExecutions {
			if s, ok := statuses[ex.NodeName]; ok {
				done, _ := strconv.ParseBool(s)
				if done {
					ex.FinishedAt = time.Now()
				}
			} else {
				ex.FinishedAt = time.Now()
			}
			ej, err := json.Marshal(ex)
			if err != nil {
				log.WithError(err).Error("agent: Error marshaling execution")
			}
			a.setExecution(ej)
		}
	}
}
