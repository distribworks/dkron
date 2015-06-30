package dcron

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *AgentCommand) dashboardRoutes(r *mux.Router) {
	subui := r.PathPrefix("/dashboard").Subrouter()
	subui.HandleFunc("/jobs", a.dashboardJobsHandler).Methods("GET")
	subui.HandleFunc("/jobs/{job}/executions", a.dashboardExecutionsHandler).Methods("GET")
}

func (a *AgentCommand) dashboardJobsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	jobs, _ := a.etcd.GetJobs()

	funcs := template.FuncMap{
		"isSuccess": func(job *Job) bool {
			return job.LastSuccess.After(job.LastError)
		},
	}

	tmpl := template.Must(template.New("dashboard.html.tmpl").Funcs(funcs).ParseFiles(
		"templates/dashboard.html.tmpl", "templates/jobs.html.tmpl"))

	data := struct {
		Jobs []*Job
	}{
		Jobs: jobs,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Error(err)
	}
}

func (a *AgentCommand) dashboardExecutionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	vars := mux.Vars(r)
	job := vars["job"]

	execs, _ := a.etcd.GetExecutions(job)

	tmpl := template.Must(template.New("dashboard.html.tmpl").Funcs(template.FuncMap{
		"html": func(value []byte) template.HTML {
			return template.HTML(value)
		},
	}).ParseFiles("templates/dashboard.html.tmpl", "templates/executions.html.tmpl"))

	data := struct {
		Executions []*Execution
		JobName    string
	}{
		Executions: execs,
		JobName:    job,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Error(err)
	}
}
