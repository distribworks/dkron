package dkron

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
)

func (a *AgentCommand) ServeHTTP() {
	r := mux.NewRouter().StrictSlash(true)

	a.apiRoutes(r)
	a.dashboardRoutes(r)

	middle := interpose.New()
	middle.UseHandler(r)

	srv := &http.Server{Addr: a.config.HTTPAddr, Handler: middle}

	log.WithFields(logrus.Fields{
		"address": a.config.HTTPAddr,
	}).Info("api: Running HTTP server")

	certFile := "" //config.GetString("certFile")
	keyFile := ""  //config.GetString("keyFile")
	if certFile != "" && keyFile != "" {
		srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		srv.ListenAndServe()
	}
	log.Info("api: Exiting HTTP server")
}

func (a *AgentCommand) apiRoutes(r *mux.Router) {
	r.Path("/v1").HandlerFunc(a.indexHandler)
	subver := r.PathPrefix("/v1").Subrouter()
	subver.HandleFunc("/members", a.membersHandler)
	subver.HandleFunc("/leader", a.leaderHandler)

	subver.Path("/jobs").HandlerFunc(a.jobCreateOrUpdateHandler).Methods("POST", "PATCH")
	subver.Path("/jobs").HandlerFunc(a.jobsHandler).Methods("GET")
	sub := subver.PathPrefix("/jobs").Subrouter()
	sub.HandleFunc("/{job}", a.jobGetHandler).Methods("GET")
	sub.HandleFunc("/{job}", a.jobDeleteHandler).Methods("DELETE")
	sub.HandleFunc("/{job}", a.jobRunHandler).Methods("POST")

	subex := subver.PathPrefix("/executions").Subrouter()
	subex.HandleFunc("/{job}", a.executionsHandler).Methods("GET")
}

func printJson(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if _, ok := r.URL.Query()["pretty"]; ok {
		j, _ := json.MarshalIndent(v, "", "\t")
		if _, err := fmt.Fprintf(w, string(j)); err != nil {
			return err
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(v); err != nil {
			return err
		}
	}

	return nil
}

func (a *AgentCommand) indexHandler(w http.ResponseWriter, r *http.Request) {
	local := a.serf.LocalMember()
	stats := map[string]map[string]string{
		"agent": {
			"name":    local.Name,
			"version": a.Version,
			"backend": a.config.Backend,
		},
		"serf": a.serf.Stats(),
		"tags": local.Tags,
	}

	if err := printJson(w, r, stats); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) jobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := a.store.GetJobs()
	if err != nil {
		log.Fatal(err)
	}

	if err := printJson(w, r, jobs); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) jobGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["job"]

	job, err := a.store.GetJob(jobName)
	if err != nil {
		log.Error(err)
	}

	if err := printJson(w, r, job); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) jobCreateOrUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	// Save the new job to the store
	if err = a.store.SetJob(&job); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	a.schedulerRestartQuery(a.store.GetLeader())

	if err := printJson(w, r, job); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) jobDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["job"]

	job, err := a.store.DeleteJob(jobName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := printJson(w, r, job); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) jobRunHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["job"]

	job, err := a.store.GetJob(jobName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	a.RunQuery(job)

	if err := printJson(w, r, job); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) executionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["job"]

	executions, err := a.store.GetExecutions(jobName)
	if err != nil {
		log.Error(err)
	}

	if err := printJson(w, r, executions); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) membersHandler(w http.ResponseWriter, r *http.Request) {
	if err := printJson(w, r, a.serf.Members()); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) leaderHandler(w http.ResponseWriter, r *http.Request) {
	member, err := a.leaderMember()
	if err == nil {
		if err := printJson(w, r, member); err != nil {
			log.Fatal(err)
		}
	}
}
