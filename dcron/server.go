package dcron

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/client"
	"github.com/hashicorp/serf/serf"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// ServerCommand run dcron server
type ServerCommand struct {
	Ui     cli.Ui
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
	cmdFlags.String("profile", "", "timing profile to use (lan, wan, local)")
	viper.SetDefault("profile", cmdFlags.Lookup("profile").Value)

	// cmdFlags.Var((*AppendSliceValue)(&configFiles), "config-file",
	// 	"json file to read config from")
	// cmdFlags.Var((*AppendSliceValue)(&configFiles), "config-dir",
	// 	"directory of json files to read")

	if err := cmdFlags.Parse(args); err != nil {
		log.Fatal(err)
	}

	config := &Config{
		NodeName:     viper.GetString("node_name"),
		BindAddr:     viper.GetString("bind_addr"),
		RPCAddr:      viper.GetString("rpc_addr"),
		HTTPAddr:     viper.GetString("http_addr"),
		Discover:     viper.GetString("discover"),
		EtcdMachines: viper.GetStringSlice("etcd_machines"),
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
		s.Ui.Error(fmt.Sprintf("Unknown profile: %s", config.Profile))
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
	serfConfig.ProtocolVersion = uint8(config.Protocol)
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

	return config
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
			if !s.isActiveMember(leader) {
				log.Debug("Trying to set itself as leader")
				res, err := s.etcd.Client.CompareAndSwap(keyspace+"/leader", s.config.NodeName, 0, leader, 0)
				if err != nil {
					log.Error(err)
				}
				log.WithFields(logrus.Fields{
					"old_leader": res.PrevNode.Value,
					"new_leader": res.Node.Value,
				}).Debug("Leader Swap")
				return true
			} else {
				log.Printf("The current leader [%s] is active", leader)
			}
		} else {
			log.Printf("This node is already the leader")
			return true
		}
	} else {
		log.Debug("Trying to set itself as leader")
		res, err := s.etcd.Client.Create(keyspace+"/leader", s.config.NodeName, 0)
		if err != nil {
			log.Error(res, err)
		}
		log.Printf("Successfully set [%s] as dcron leader", s.config.NodeName)
		return true
	}

	return false
}

func (s *ServerCommand) isActiveMember(memberName string) bool {
	members, err := s.serf.Members()
	if err != nil {
		log.Fatal("Error listing cluster members")
	}

	for _, member := range members {
		if member.Name == memberName && member.Status == "alive" {
			return true
		}
	}
	return false
}

// dcron leader init routine
func (s *ServerCommand) ListenEvents() {
	ch := make(chan map[string]interface{})

	sh, err := s.serf.Stream("*", ch)
	if err != nil {
		log.Error(err)
	}
	defer s.serf.Stop(sh)

	for {
		select {
		case event := <-ch:
			switch event["Event"] {
			case "member-leave", "member-failed":
				members, err := s.decodeMembers(event)
				if err != nil {
					log.Fatal("Event member: ", err)
				}

				for _, member := range members {
					if member.Name == s.etcd.GetLeader() && member.Status != "alive" {
						if s.ElectLeader() {
							log.Debug("Restarting scheduler")
							s.schedulerRestart()
						}
					}
				}
			case "query":
				if event["Name"] == "scheduler:reload" {
					s.schedulerRestart()
				}
			}
		}
	}
}

func (s *ServerCommand) decodeMembers(event map[string]interface{}) ([]*client.Member, error) {
	var members []*client.Member

	membersMap, ok := event["Members"].([]interface{})
	if !ok {
		return nil, errors.New("Error decoding members.")
	}

	for _, memberItem := range membersMap {
		memberMap, ok := memberItem.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("Error decoding members.")
		}

		m := make(map[string]interface{}, 1)
		for key, val := range memberMap {
			m[key.(string)] = val
		}
		var member client.Member
		mapstructure.Decode(m, &member)
		members = append(members, &member)
		log.Debug(m)
	}

	return members, nil
}

func (s *ServerCommand) schedulerRestart() {
	// Restart scheduler
	jobs, err := s.etcd.GetJobs()
	if err != nil {
		log.Fatal(err)
	}
	sched.Restart(jobs)
}
