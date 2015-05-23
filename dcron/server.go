package dcron

import (
	"flag"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	serfs "github.com/hashicorp/serf/serf"
	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

// ServerCommand run dcron server
type ServerCommand struct {
	Ui     cli.Ui
	Leader string
	config *Config
	serf   *serfManager
	etcd   *etcdClient
}

func (s *ServerCommand) Help() string {
	helpText := `
Usage: dcron server [options]
	Run dcron server
Options:
`
	return strings.TrimSpace(helpText)
}

// readConfig is responsible for setup of our configuration using
// the command line and any file configs
func (s *ServerCommand) readConfig(args []string) *Config {
	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	cmdFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	cmdFlags.Usage = func() { s.Ui.Output(s.Help()) }
	cmdFlags.String("node", hostname, "node name")
	viper.SetDefault("node_name", cmdFlags.Lookup("node").Value)
	cmdFlags.String("bind", "0.0.0.0:7946", "address to bind listeners to")
	viper.SetDefault("bind_addr", cmdFlags.Lookup("bind").Value)
	cmdFlags.String("rpc-addr", "127.0.0.1:7373", "RPC address")
	viper.SetDefault("rpc_addr", cmdFlags.Lookup("rpc-addr").Value)
	cmdFlags.String("http-addr", ":8080", "HTTP address")
	viper.SetDefault("http_addr", cmdFlags.Lookup("http-addr").Value)
	cmdFlags.String("discover", "dcron", "mDNS discovery name")
	viper.SetDefault("discover", cmdFlags.Lookup("discover").Value)
	cmdFlags.String("etcd-machines", "127.0.0.1:2379", "etcd machines addresses")
	viper.SetDefault("etcd_machines", cmdFlags.Lookup("etcd-machines").Value)

	// cmdFlags.Var((*AppendSliceValue)(&configFiles), "config-file",
	// 	"json file to read config from")
	// cmdFlags.Var((*AppendSliceValue)(&configFiles), "config-dir",
	// 	"directory of json files to read")

	if err := cmdFlags.Parse(args); err != nil {
		log.Fatal(err)
	}

	return &Config{
		NodeName:     viper.GetString("node_name"),
		BindAddr:     viper.GetString("bind_addr"),
		RPCAddr:      viper.GetString("rpc_addr"),
		HTTPAddr:     viper.GetString("http_addr"),
		Discover:     viper.GetString("discover"),
		EtcdMachines: viper.GetStringSlice("etcd_machines"),
	}
}

func (s *ServerCommand) Run(args []string) int {
	s.config = s.readConfig(args)
	s.serf = NewSerfManager(s.config)
	s.etcd = NewEtcdClient(s.config.EtcdMachines)

	go func() {
		defer func() {
			s.serf.Terminate()
		}()
		s.ServeHTTP()
	}()

	s.serf.Start()
	defer func() {
		s.serf.Terminate()
	}()
	if s.ElectLeader() {
		jobs, err := s.etcd.GetJobs()
		if err != nil {
			log.Fatal(err)
		}
		sched.Start(jobs)
	}
	s.ListenEvents()

	return 0
}

func (s *ServerCommand) Synopsis() string {
	return "Run dcron server"
}

// dcron leader init routine
func (s *ServerCommand) ElectLeader() bool {
	leader := s.etcd.GetLeader()

	if leader != "" {
		if leader != s.config.NodeName {
			if !s.isMember(leader) {
				log.Debug("Trying to set itself as leader")
				res, err := s.etcd.Client.CompareAndSwap(keyspace+"/leader", s.config.NodeName, 0, leader, 0)
				if err != nil {
					log.Error(err)
					return false
				}
				log.WithFields(logrus.Fields{
					"old_leader": res.PrevNode.Value,
					"new_leader": res.Node.Value,
				}).Debug("Leader Swap")
			} else {
				log.Printf("The current leader [%s] is active", leader)
				return false
			}
		} else {
			log.Printf("This node is already the leader")
		}
	} else {
		log.Debug("Trying to set itself as leader")
		res, err := s.etcd.Client.Create(keyspace+"/leader", s.config.NodeName, 0)
		if err != nil {
			log.Error(res, err)
			return false
		}
		s.Leader = s.config.NodeName
		log.Printf("Successfully set [%s] as dcron leader", s.Leader)
	}

	return true
}

func (s *ServerCommand) isMember(memberName string) bool {
	members, err := s.serf.Members()
	if err != nil {
		log.Fatal("Error listing cluster members")
	}

	for _, member := range members {
		if member.Name == memberName {
			return true
		}
	}
	return false
}

// dcron leader init routine
func (s *ServerCommand) ListenEvents() {
	ch := make(chan map[string]interface{}, 100)

	sh, err := s.serf.Stream("*", ch)
	if err != nil {
		log.Error(err)
	}
	defer s.serf.Stop(sh)

	for {
		select {
		case event := <-ch:
			for key, val := range event {
				switch ev := val.(type) {
				case serfs.MemberEvent:
					if ev.Type == serfs.EventMemberLeave {
						for _, member := range ev.Members {
							if member.Name == s.Leader {
								s.ElectLeader()
							}
						}
					}

					log.Debug(ev)
				default:
					log.Debugf("Receiving event: %s => %v of type %T", key, val, val)
				}
			}
			if event["Event"] == "query" {
				if event["Payload"] != nil {
					log.Debug(string(event["Payload"].([]byte)))
					s.serf.Respond(uint64(event["ID"].(int64)), []byte("Peetttee"))
				}
			}
		}
	}
}
