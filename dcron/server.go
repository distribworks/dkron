package dcron

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/mitchellh/cli"
	"github.com/tylerb/graceful"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ServerCommand run dcron server
type ServerCommand struct {
	Ui cli.Ui
}

func (s *ServerCommand) Help() string {
	helpText := `
Usage: dcron server [options]
	Provides debugging information for operators
Options:
  -format                  If provided, output is returned in the specified
                           format. Valid formats are 'json', and 'text' (default)
`
	return strings.TrimSpace(helpText)
}

func (s *ServerCommand) Run(args []string) int {
	var format string
	cmdFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	cmdFlags.Usage = func() { s.Ui.Output(s.Help()) }
	cmdFlags.StringVar(&format, "format", "text", "output format")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	var wg sync.WaitGroup
	wg.Add(2)

	s.ElectLeader()

	go func() {
		defer func() {
			serf.Terminate()
		}()
		s.ServeHTTP()
		wg.Done()
	}()

	sched.Load()

	go func() {
		defer func() {
			serf.Terminate()
		}()
		initSerf()
		wg.Done()
	}()

	wg.Wait()

	return 0
}

func (s *ServerCommand) Synopsis() string {
	return "Run dcron server"
}

func (s *ServerCommand) ServeHTTP() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", IndexHandler)
	sub := r.PathPrefix("/jobs").Subrouter()
	sub.HandleFunc("/", JobCreateHandler).Methods("POST")
	sub.HandleFunc("/", JobsHandler).Methods("GET")

	middle := interpose.New()
	middle.UseHandler(r)

	srv := &graceful.Server{
		Timeout: 1 * time.Second,
		Server:  &http.Server{Addr: ":8080", Handler: middle},
	}

	log.Infoln("Running HTTP server on 8080")

	certFile := config.GetString("certFile")
	keyFile := config.GetString("keyFile")
	if certFile != "" && keyFile != "" {
		srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		srv.ListenAndServe()
	}
	log.Debug("Exiting HTTP server")
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	stats, err := serf.Stats()
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

func JobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := etcd.GetJobs()
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

func JobCreateHandler(w http.ResponseWriter, r *http.Request) {
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
	if err = etcd.SetJob(&job); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Schedule the new job
	sched.Reload()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := fmt.Fprintf(w, `{"result": "ok"}`); err != nil {
		log.Fatal(err)
	}
}

func ExecutionsHandler(w http.ResponseWriter, r *http.Request) {
	executions, err := etcd.GetExecutions()
	if err != nil {
		log.Error(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(executions); err != nil {
		panic(err)
	}
}

// Get leader key and store it globally
// If it exists check agains the member list if the member is still there
// If it is, do nothing
// If not exists try to CompareAndSwap with the current node_name setting the self node_name
// If successful load the scheduler
// On failure do nothing and listen for member-leave events
func (s *ServerCommand) ElectLeader() {
	leader := etcd.GetLeader()
	if leader != "" {

	}
	log.Debug(leader)
}
