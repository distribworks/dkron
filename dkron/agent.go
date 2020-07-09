package dkron

import (
	"crypto/tls"
	"errors"
	"expvar"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/distribworks/dkron/v3/plugin"
	proto "github.com/distribworks/dkron/v3/plugin/types"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/hashicorp/serf/serf"
	"github.com/juliangruber/go-intersect"
	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	raftTimeout = 30 * time.Second
	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 512
	minRaftProtocol  = 3
)

var (
	expNode = expvar.NewString("node")

	// ErrLeaderNotFound is returned when obtained leader is not found in member list
	ErrLeaderNotFound = errors.New("no member leader found in member list")

	// ErrNoSuitableServer returns an error in case no suitable server to send the request is found.
	ErrNoSuitableServer = errors.New("no suitable server found to send the request, aborting")

	runningExecutions sync.Map
)

// Agent is the main struct that represents a dkron agent
type Agent struct {
	// ProcessorPlugins maps processor plugins
	ProcessorPlugins map[string]plugin.Processor

	//ExecutorPlugins maps executor plugins
	ExecutorPlugins map[string]plugin.Executor

	// HTTPTransport is a swappable interface for the HTTP server interface
	HTTPTransport Transport

	// Store interface to set the storage engine
	Store Storage

	// GRPCServer interface for setting the GRPC server
	GRPCServer DkronGRPCServer

	// GRPCClient interface for setting the GRPC client
	GRPCClient DkronGRPCClient

	// TLSConfig allows setting a TLS config for transport
	TLSConfig *tls.Config

	// Pro features
	GlobalLock         bool
	MemberEventHandler func(serf.Event)
	ProAppliers        LogAppliers

	serf        *serf.Serf
	config      *Config
	eventCh     chan serf.Event
	sched       *Scheduler
	ready       bool
	shutdownCh  chan struct{}
	retryJoinCh chan error

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
	peers      map[string][]*ServerParts
	localPeers map[raft.ServerAddress]*ServerParts
	peerLock   sync.RWMutex

	activeExecutions sync.Map

	listener net.Listener
}

// ProcessorFactory is a function type that creates a new instance
// of a processor.
type ProcessorFactory func() (plugin.Processor, error)

// Plugins struct to store loaded plugins of each type
type Plugins struct {
	Processors map[string]plugin.Processor
	Executors  map[string]plugin.Executor
}

// AgentOption type that defines agent options
type AgentOption func(agent *Agent)

// NewAgent return a new Agent instace capable of starting
// and running a Dkron instance.
func NewAgent(config *Config, options ...AgentOption) *Agent {
	agent := &Agent{
		config:      config,
		retryJoinCh: make(chan error),
	}

	for _, option := range options {
		option(agent)
	}

	return agent
}

// Start the current agent by running all the necessary
// checks and server or client routines.
func (a *Agent) Start() error {
	InitLogger(a.config.LogLevel, a.config.NodeName)

	// Normalize configured addresses
	a.config.normalizeAddrs()

	s, err := a.setupSerf()
	if err != nil {
		return fmt.Errorf("agent: Can not setup serf, %s", err)
	}
	a.serf = s

	// start retry join
	if len(a.config.RetryJoinLAN) > 0 {
		a.retryJoinLAN()
	} else {
		a.join(a.config.StartJoin, true)
	}

	if err := initMetrics(a); err != nil {
		log.Fatal("agent: Can not setup metrics")
	}

	// Expose the node name
	expNode.Set(a.config.NodeName)

	//Use the value of "RPCPort" if AdvertiseRPCPort has not been set
	if a.config.AdvertiseRPCPort <= 0 {
		a.config.AdvertiseRPCPort = a.config.RPCPort
	}

	// Create a listener for RPC subsystem
	addr := a.bindRPCAddr()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	a.listener = l

	if a.config.Server {
		a.StartServer()
	} else {
		opts := []grpc.ServerOption{}
		if a.TLSConfig != nil {
			tc := credentials.NewTLS(a.TLSConfig)
			opts = append(opts, grpc.Creds(tc))
		}

		grpcServer := grpc.NewServer(opts...)
		as := NewAgentServer(a)
		proto.RegisterAgentServer(grpcServer, as)
		go grpcServer.Serve(l)
	}

	if a.GRPCClient == nil {
		a.GRPCClient = NewGRPCClient(nil, a)
	}

	tags := a.serf.LocalMember().Tags
	tags["rpc_addr"] = a.advertiseRPCAddr() // Address that clients will use to RPC to servers
	tags["port"] = strconv.Itoa(a.config.AdvertiseRPCPort)
	a.serf.SetTags(tags)

	go a.eventLoop()
	a.ready = true

	return nil
}

// RetryJoinCh is a channel that transports errors
// from the retry join process.
func (a *Agent) RetryJoinCh() <-chan error {
	return a.retryJoinCh
}

// JoinLAN is used to have Dkron join the inner-DC pool
// The target address should be another node inside the DC
// listening on the Serf LAN address
func (a *Agent) JoinLAN(addrs []string) (int, error) {
	return a.serf.Join(addrs, true)
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
		a.sched.ClearCron()
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

	// Raft performance
	raftMultiplier := a.config.RaftMultiplier
	if raftMultiplier < 1 || raftMultiplier > 10 {
		return fmt.Errorf("raft-multiplier cannot be %d. Must be between 1 and 10", raftMultiplier)
	}
	config.HeartbeatTimeout = config.HeartbeatTimeout * time.Duration(raftMultiplier)
	config.ElectionTimeout = config.ElectionTimeout * time.Duration(raftMultiplier)
	config.LeaderLeaseTimeout = config.LeaderLeaseTimeout * time.Duration(a.config.RaftMultiplier)

	config.LogOutput = logger
	config.LocalID = raft.ServerID(a.config.NodeName)

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
			return fmt.Errorf("error creating new raft store: %s", err)
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

		// Check for peers.json file for recovery
		peersFile := filepath.Join(a.config.DataDir, "raft", "peers.json")
		if _, err := os.Stat(peersFile); err == nil {
			log.Info("found peers.json file, recovering Raft configuration...")
			var configuration raft.Configuration
			configuration, err = raft.ReadConfigJSON(peersFile)
			if err != nil {
				return fmt.Errorf("recovery failed to parse peers.json: %v", err)
			}
			tmpFsm := newFSM(nil, nil)
			if err := raft.RecoverCluster(config, tmpFsm,
				logStore, stableStore, snapshots, transport, configuration); err != nil {
				return fmt.Errorf("recovery failed: %v", err)
			}
			if err := os.Remove(peersFile); err != nil {
				return fmt.Errorf("recovery failed to delete peers.json, please delete manually (see peers.info for details): %v", err)
			}
			log.Info("deleted peers.json file after successful recovery")
		}
	}

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
	fsm := newFSM(a.Store, a.ProAppliers)
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
	a.localPeers = make(map[raft.ServerAddress]*ServerParts)
	a.peers = make(map[string][]*ServerParts)

	bindIP, bindPort, err := config.AddrParts(config.BindAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid bind address: %s", err)
	}

	var advertiseIP string
	var advertisePort int
	if config.AdvertiseAddr != "" {
		advertiseIP, advertisePort, err = config.AddrParts(config.AdvertiseAddr)
		if err != nil {
			return nil, fmt.Errorf("invalid advertise address: %s", err)
		}
	}

	encryptKey, err := config.EncryptBytes()
	if err != nil {
		return nil, fmt.Errorf("invalid encryption key: %s", err)
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
		return nil, fmt.Errorf("unknown profile: %s", config.Profile)
	}

	serfConfig.MemberlistConfig.BindAddr = bindIP
	serfConfig.MemberlistConfig.BindPort = bindPort
	serfConfig.MemberlistConfig.AdvertiseAddr = advertiseIP
	serfConfig.MemberlistConfig.AdvertisePort = advertisePort
	serfConfig.MemberlistConfig.SecretKey = encryptKey
	serfConfig.NodeName = config.NodeName
	serfConfig.Tags = config.Tags
	serfConfig.CoalescePeriod = 3 * time.Second
	serfConfig.QuiescentPeriod = time.Second
	serfConfig.UserCoalescePeriod = 3 * time.Second
	serfConfig.UserQuiescentPeriod = time.Second
	serfConfig.ReconnectTimeout, err = time.ParseDuration(config.SerfReconnectTimeout)

	if err != nil {
		log.Fatal(err)
	}

	// Create a channel to listen for events from Serf
	a.eventCh = make(chan serf.Event, 2048)
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
		s, err := NewStore()
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

	// Create a cmux object.
	tcpm := cmux.New(a.listener)
	var grpcl, raftl net.Listener

	// If TLS config present listen to TLS
	if a.TLSConfig != nil {
		// Create a RaftLayer with TLS
		a.raftLayer = NewTLSRaftLayer(a.TLSConfig)

		// Match any connection to the recursive mux
		tlsl := tcpm.Match(cmux.Any())
		tlsl = tls.NewListener(tlsl, a.TLSConfig)

		// Declare sub cMUX for TLS
		tlsm := cmux.New(tlsl)

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

// Members is used to return the members of the serf cluster
func (a *Agent) Members() []serf.Member {
	return a.serf.Members()
}

// LocalMember is used to return the local node
func (a *Agent) LocalMember() serf.Member {
	return a.serf.LocalMember()
}

// Leader is used to return the Raft leader
func (a *Agent) Leader() raft.ServerAddress {
	return a.raft.Leader()
}

// Servers returns a list of known server
func (a *Agent) Servers() (members []*ServerParts) {
	for _, member := range a.serf.Members() {
		ok, parts := isServer(member)
		if !ok || member.Status != serf.StatusAlive {
			continue
		}
		members = append(members, parts)
	}
	return members
}

// LocalServers returns a list of the local known server
func (a *Agent) LocalServers() (members []*ServerParts) {
	for _, member := range a.serf.Members() {
		ok, parts := isServer(member)
		if !ok || member.Status != serf.StatusAlive {
			continue
		}
		if a.config.Region == parts.Region {
			members = append(members, parts)
		}
	}
	return members
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

				if a.MemberEventHandler != nil {
					a.MemberEventHandler(e)
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
			}

		case <-serfShutdownCh:
			log.Warn("agent: Serf shutdown detected, quitting")
			return
		}
	}
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

func (a *Agent) processFilteredNodes(job *Job) (map[string]string, map[string]string, error) {
	// candidates will contain a set of candidates by tags
	// the final set of nodes will be the intesection of all groups
	tags := make(map[string]string)

	// Actually copy the map
	for key, val := range job.Tags {
		tags[key] = val
	}

	// Always filter by region tag as we currently only target nodes
	// on the same region.
	tags["region"] = a.config.Region

	candidates := [][]string{}
	for jtk, jtv := range tags {
		cans := []string{}
		tc := strings.Split(jtv, ":")

		tv := tc[0]

		// Set original tag to clean tag
		tags[jtk] = tv

		for _, member := range a.serf.Members() {
			if member.Status == serf.StatusAlive {
				for mtk, mtv := range member.Tags {
					if mtk == jtk && mtv == tv {
						cans = append(cans, member.Name)
					}
				}
			}
		}

		// In case there is cardinality in the tag, randomize the order and select the amount of nodes
		// or else just add all nodes to the result.
		if len(tc) == 2 {
			f := []string{}
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(cans), func(i, j int) {
				cans[i], cans[j] = cans[j], cans[i]
			})

			count, err := strconv.Atoi(tc[1])
			if err != nil {
				return nil, nil, err
			}
			for i := 1; i <= count; i++ {
				if len(cans) == 0 {
					break
				}
				f = append(f, cans[0])
				cans = cans[1:]
			}
			cans = f
		}

		candidates = append(candidates, cans)
	}

	// The final result will be the intersection of all candidates.
	nodes := make(map[string]string)
	r := candidates[0]
	for i := 1; i <= len(candidates)-1; i++ {
		isec := intersect.Simple(r, candidates[i]).([]interface{})
		// Empty the slice
		r = []string{}

		// Refill with the intersection
		for _, v := range isec {
			r = append(r, v.(string))
		}
	}

	for _, n := range r {
		for _, m := range a.serf.Members() {
			if n == m.Name {
				// If the server is missing the rpc_addr tag, default to the serf advertise addr
				if addr, ok := m.Tags["rpc_addr"]; ok {
					nodes[n] = addr
				} else {
					nodes[n] = m.Addr.String()
				}
			}
		}
	}

	return nodes, tags, nil
}

// This function is called when a client request the RPCAddress
// of the current member.
// in marathon, it would return the host's IP and advertise RPC port
func (a *Agent) advertiseRPCAddr() string {
	bindIP := a.serf.LocalMember().Addr
	return net.JoinHostPort(bindIP.String(), strconv.Itoa(a.config.AdvertiseRPCPort))
}

// Get bind address for RPC
func (a *Agent) bindRPCAddr() string {
	bindIP, _, _ := a.config.AddrParts(a.config.BindAddr)
	return net.JoinHostPort(bindIP, strconv.Itoa(a.config.RPCPort))
}

// applySetJob is a helper method to be called when
// a job property need to be modified from the leader.
func (a *Agent) applySetJob(job *proto.Job) error {
	cmd, err := Encode(SetJobType, job)
	if err != nil {
		return err
	}
	af := a.raft.Apply(cmd, raftTimeout)
	if err := af.Error(); err != nil {
		return err
	}
	res := af.Response()
	switch res {
	case ErrParentJobNotFound:
		return ErrParentJobNotFound
	case ErrSameParent:
		return ErrParentJobNotFound
	}

	return nil
}

// RaftApply applies a command to the Raft log
func (a *Agent) RaftApply(cmd []byte) raft.ApplyFuture {
	return a.raft.Apply(cmd, raftTimeout)
}

// GetRunningJobs returns amount of active jobs of the local agent
func (a *Agent) GetRunningJobs() int {
	job := 0
	runningExecutions.Range(func(k, v interface{}) bool {
		job = job + 1
		return true
	})
	return job
}

// GetActiveExecutions returns running executions globally
func (a *Agent) GetActiveExecutions() ([]*proto.Execution, error) {
	var executions []*proto.Execution

	for _, s := range a.LocalServers() {
		exs, err := a.GRPCClient.GetActiveExecutions(s.RPCAddr.String())
		if err != nil {
			return nil, err
		}

		executions = append(executions, exs...)
	}

	return executions, nil
}

func (a *Agent) recursiveSetJob(jobs []*Job) []string {
	result := make([]string, 0)
	for _, job := range jobs {
		err := a.GRPCClient.SetJob(job)
		if err != nil {
			result = append(result, "fail create "+job.Name)
			continue
		} else {
			result = append(result, "success create "+job.Name)
			if len(job.ChildJobs) > 0 {
				recursiveResult := a.recursiveSetJob(job.ChildJobs)
				result = append(result, recursiveResult...)
			}
		}
	}
	return result
}

// Check if the server is alive and select it
func (a *Agent) checkAndSelectServer() (string, error) {
	var peers []string
	for _, p := range a.LocalServers() {
		peers = append(peers, p.RPCAddr.String())
	}

	for _, peer := range peers {
		log.WithField("peer", peer).Debug("Checking peer")
		conn, err := net.DialTimeout("tcp", peer, 1*time.Second)
		if err == nil {
			conn.Close()
			log.WithField("peer", peer).Debug("Found good peer")
			return peer, nil
		}
	}
	return "", ErrNoSuitableServer
}
