package dcron

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	// Path of static files must be last!
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	srv := &http.Server{Addr: a.config.HTTPAddr, Handler: middle}

	log.Infof("Running HTTP server on %s", a.config.HTTPAddr)

	certFile := "" //config.GetString("certFile")
	keyFile := ""  //config.GetString("keyFile")
	if certFile != "" && keyFile != "" {
		srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		srv.ListenAndServe()
	}
	log.Debug("Exiting HTTP server")
}

func (a *AgentCommand) apiRoutes(r *mux.Router) {
	r.HandleFunc("/", a.indexHandler)
	r.HandleFunc("/members", a.membersHandler)
	r.HandleFunc("/leader", a.leaderHandler)

	sub := r.PathPrefix("/jobs").Subrouter()
	sub.HandleFunc("/", a.jobCreateOrUpdateHandler).Methods("POST", "PUT")
	sub.HandleFunc("/", a.jobsHandler).Methods("GET")
	sub.HandleFunc("/{job}", a.jobDeleteHandler).Methods("DELETE")
	sub.HandleFunc("/{job}", a.jobRunHandler).Methods("POST", "PUT")

	subex := r.PathPrefix("/executions").Subrouter()
	subex.HandleFunc("/{job}", a.executionsHandler).Methods("GET")
}

func (a *AgentCommand) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	local := a.serf.LocalMember()
	stats := map[string]map[string]string{
		"agent": map[string]string{
			"name":    local.Name,
			"version": a.Version,
		},
		"serf": a.serf.Stats(),
		"tags": local.Tags,
	}

	statsJson, _ := json.MarshalIndent(stats, "", "\t")
	if _, err := fmt.Fprintf(w, string(statsJson)); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) jobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := a.etcd.GetJobs()
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

	// Save the new job to etcd
	if err = a.etcd.SetJob(&job); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	a.schedulerRestartQuery(a.etcd.GetLeader())

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := fmt.Fprintf(w, `{"result": "ok"}`); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) executionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job := vars["job"]

	executions, err := a.etcd.GetExecutions(job)
	if err != nil {
		log.Error(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(executions); err != nil {
		panic(err)
	}
}

func (a *AgentCommand) membersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(a.serf.Members()); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) jobDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job := vars["job"]

	if _, err := a.etcd.Client.Delete(job, false); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, `{"result": "ok"}`); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) leaderHandler(w http.ResponseWriter, r *http.Request) {
	member, err := a.leaderMember()
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(member); err != nil {
			log.Fatal(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) jobRunHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job := vars["job"]

	j, err := a.etcd.GetJob(job)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatal(err)
		}
		return
	}

	a.RunQuery(j)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, `{"result": "ok"}`); err != nil {
		log.Fatal(err)
	}
}
