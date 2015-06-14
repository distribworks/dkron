package dcron

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	// "github.com/tylerb/graceful"
	"io"
	"io/ioutil"
)

func (a *AgentCommand) ServeHTTP() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", a.IndexHandler)
	r.HandleFunc("/members", a.MembersHandler)
	r.HandleFunc("/leader", a.LeaderHandler)

	subui := r.PathPrefix("/ui").Subrouter()
	subui.HandleFunc("/dashboard", a.DashboardHandler)

	sub := r.PathPrefix("/jobs").Subrouter()
	sub.HandleFunc("/", a.JobCreateOrUpdateHandler).Methods("POST", "PUT")
	sub.HandleFunc("/", a.JobsHandler).Methods("GET")
	sub.HandleFunc("/{job}", a.JobDeleteHandler).Methods("DELETE")

	middle := interpose.New()
	middle.UseHandler(r)

	// Path of static files must be last!
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	srv := &http.Server{Addr: a.config.HTTPAddr, Handler: middle}

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

func (a *AgentCommand) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	local := a.serf.LocalMember()
	stats := map[string]map[string]string{
		"agent": map[string]string{
			"name": local.Name,
		},
		"serf": a.serf.Stats(),
		"tags": local.Tags,
	}

	statsJson, _ := json.MarshalIndent(stats, "", "\t")
	if _, err := fmt.Fprintf(w, string(statsJson)); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) JobsHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *AgentCommand) JobCreateOrUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	a.schedulerReloadQuery(a.etcd.GetLeader())

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := fmt.Fprintf(w, `{"result": "ok"}`); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) ExecutionsHandler(w http.ResponseWriter, r *http.Request) {
	executions, err := a.etcd.GetExecutions()
	if err != nil {
		log.Error(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(executions); err != nil {
		panic(err)
	}
}

func (a *AgentCommand) MembersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(a.serf.Members()); err != nil {
		log.Fatal(err)
	}
}

func (a *AgentCommand) JobDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job := vars["job"]
	// Save the new job to etcd
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

func (a *AgentCommand) LeaderHandler(w http.ResponseWriter, r *http.Request) {
	leader := a.etcd.GetLeader()
	for _, member := range a.serf.Members() {
		if key, ok := member.Tags["key"]; ok {
			if key == leader {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(member); err != nil {
					log.Fatal(err)
				}
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

func (a *AgentCommand) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	jobs, _ := a.etcd.GetJobs()

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")
	if err != nil {
		return
	}

	funcs := template.FuncMap{
		"isSuccess": func() bool {
			return true //job.LastSuccess.After(job.LastError)
		},
	}
	tmpl.Funcs(funcs)

	data := struct {
		Jobs []*Job
	}{
		Jobs: jobs,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Error(err)
	}
}
