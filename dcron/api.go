package dcron

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/tylerb/graceful"
	"io"
	"io/ioutil"
)

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
