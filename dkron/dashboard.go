package dkron

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type commonDashboardData struct {
	Version    string
	LeaderName string
	MemberName string
	Backend    string
}

func newCommonDashboardData(a *AgentCommand, nodeName string) *commonDashboardData {
	l, _ := a.leaderMember()
	return &commonDashboardData{
		Version:    a.Version,
		LeaderName: l.Name,
		MemberName: nodeName,
		Backend:    a.config.Backend,
	}
}

func (a *AgentCommand) dashboardRoutes(r *mux.Router) {
	r.Path("/dashboard").HandlerFunc(a.dashboardIndexHandler).Methods("GET")
	subui := r.PathPrefix("/dashboard").Subrouter()
	subui.HandleFunc("/jobs", a.dashboardJobsHandler).Methods("GET")
	subui.HandleFunc("/jobs/{job}/executions", a.dashboardExecutionsHandler).Methods("GET")
}

func (a *AgentCommand) dashboardIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl := template.Must(template.New("dashboard.html.tmpl").ParseFiles(
		"templates/dashboard.html.tmpl", "templates/index.html.tmpl", "templates/status.html.tmpl"))

	data := struct {
		Common    *commonDashboardData
		StartTime string
	}{
		Common: newCommonDashboardData(a, a.config.NodeName),
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Error(err)
	}
}

func (a *AgentCommand) dashboardJobsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	jobs, _ := a.store.GetJobs()

	funcs := template.FuncMap{
		"executionStatus": func(job *Job) string {
			execs, _ := a.store.GetLastExecutionGroup(job.Name)
			success := 0
			failed := 0
			for _, ex := range execs {
				if ex.Success {
					success = success + 1
				} else {
					failed = failed + 1
				}
			}

			if failed == 0 {
				return "success"
			} else if failed > 0 && success == 0 {
				return "danger"
			} else if failed > 0 && success > 0 {
				return "warning"
			}

			return ""
		},
	}

	tmpl := template.Must(template.New("dashboard.html.tmpl").Funcs(funcs).ParseFiles(
		"templates/dashboard.html.tmpl", "templates/jobs.html.tmpl", "templates/status.html.tmpl"))

	data := struct {
		Common *commonDashboardData
		Jobs   []*Job
	}{
		Common: newCommonDashboardData(a, a.config.NodeName),
		Jobs:   jobs,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Error(err)
	}
}

func (a *AgentCommand) dashboardExecutionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	vars := mux.Vars(r)
	job := vars["job"]

	execs, _ := a.store.GetExecutions(job)

	tmpl := template.Must(template.New("dashboard.html.tmpl").Funcs(template.FuncMap{
		"html": func(value []byte) template.HTML {
			return template.HTML(value)
		},
	}).ParseFiles("templates/dashboard.html.tmpl", "templates/executions.html.tmpl"))

	if len(execs) > 100 {
		execs = execs[len(execs)-100:]
	}

	data := struct {
		Common     *commonDashboardData
		Executions []*Execution
		JobName    string
	}{
		Common:     newCommonDashboardData(a, a.config.NodeName),
		Executions: execs,
		JobName:    job,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Error(err)
	}
}
