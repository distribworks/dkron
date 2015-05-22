package dcron

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	serfs "github.com/hashicorp/serf/serf"
	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"
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
		EtcdMachines: viper.GetStringSlice("etcd_machines"),
	}
}

func (s *ServerCommand) Run(args []string) int {
	var wg sync.WaitGroup

	s.config = s.readConfig(args)
	s.serf = NewSerfManager(s.config)
	s.etcd = NewEtcdClient(s.config.EtcdMachines)

	wg.Add(1)
	go func() {
		defer func() {
			s.serf.Terminate()
		}()
		s.ServeHTTP()
		wg.Done()
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

	wg.Wait()

	return 0
}

func (s *ServerCommand) Synopsis() string {
	return "Run dcron server"
}

func (s *ServerCommand) ServeHTTP() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", s.IndexHandler)
	sub := r.PathPrefix("/jobs").Subrouter()
	sub.HandleFunc("/", s.JobCreateHandler).Methods("POST")
	sub.HandleFunc("/", s.JobsHandler).Methods("GET")

	middle := interpose.New()
	middle.UseHandler(r)

	srv := &graceful.Server{
		Timeout: 1 * time.Second,
		Server:  &http.Server{Addr: s.config.HTTPAddr, Handler: middle},
	}

	log.Infoln("Running HTTP server on 8080")

	certFile := "" //config.GetString("certFile")
	keyFile := ""  //config.GetString("keyFile")
	if certFile != "" && keyFile != "" {
		srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		srv.ListenAndServe()
	}
	log.Debug("Exiting HTTP server")
}

func (s *ServerCommand) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	stats, err := s.serf.Stats()
	if err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	statsJson, _ := json.MarshalIndent(stats, "", "\t")
	if _, err := fmt.Fprintf(w, string(statsJson)); err != nil {
		log.Fatal(err)
	}
}

func (s *ServerCommand) JobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := s.etcd.GetJobs()
	if err != nil {
		log.Error(err)
	}
	log.Debug(jobs)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		log.Fatal(err)
	}
}

func (s *ServerCommand) JobCreateHandler(w http.ResponseWriter, r *http.Request) {
	var job Job
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(body, &job); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Save the new job to etcd
	if err = s.etcd.SetJob(&job); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Schedule the new job
	jobs, err := s.etcd.GetJobs()
	if err != nil {
		log.Fatal(err)
	}
	sched.Restart(jobs)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := fmt.Fprintf(w, `{"result": "ok"}`); err != nil {
		log.Fatal(err)
	}
}

func (s *ServerCommand) ExecutionsHandler(w http.ResponseWriter, r *http.Request) {
	executions, err := s.etcd.GetExecutions()
	if err != nil {
		log.Error(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(executions); err != nil {
		panic(err)
	}
}

// dcron leader init routine
func (s *ServerCommand) ElectLeader() bool {
	stats, err := s.serf.Stats()
	if err != nil {
		log.Fatal(err)
	}

	memberName := stats["agent"]["name"]
	leader := s.etcd.GetLeader()

	if leader != "" {
		if leader != memberName && !s.isMember(leader) {
			log.Debug("Trying to set itself as leader")
			res, err := s.etcd.Client.CompareAndSwap(keyspace+"/leader", memberName, 0, leader, 0)
			if err != nil {
				log.Error(err)
				return false
			}
			log.WithFields(logrus.Fields{
				"old_leader": res.PrevNode.Value,
				"new_leader": res.Node.Value,
			}).Debug("Leader Swap")
		} else {
			log.Printf("This node is already the leader")
		}
	} else {
		log.Debug("Trying to set itself as leader")
		res, err := s.etcd.Client.Create(keyspace+"/leader", memberName, 0)
		if err != nil {
			log.Error(res, err)
			return false
		}
		s.Leader = memberName
		log.Printf("Successfully set %s as dcron leader", memberName)
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
	ch := make(chan map[string]interface{}, 1)

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
