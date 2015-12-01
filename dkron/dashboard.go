package dkron

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

const (
	tmplPath            = "templates"
	dashboardPathPrefix = "dashboard"
	apiPathPrefix       = "v1"
)

type commonDashboardData struct {
	Version    string
	LeaderName string
	MemberName string
	Backend    string
	Path       string
	APIPath    string
}

func newCommonDashboardData(a *AgentCommand, nodeName, path string) *commonDashboardData {
	l, _ := a.leaderMember()
	return &commonDashboardData{
		Version:    a.Version,
		LeaderName: l.Name,
		MemberName: nodeName,
		Backend:    a.config.Backend,
		Path:       fmt.Sprintf("%s%s", path, dashboardPathPrefix),
		APIPath:    fmt.Sprintf("%s%s", path, apiPathPrefix),
	}
}

func (a *AgentCommand) dashboardRoutes(r *mux.Router) {
	r.Path("/" + dashboardPathPrefix).HandlerFunc(a.dashboardIndexHandler).Methods("GET")
	subui := r.PathPrefix("/" + dashboardPathPrefix).Subrouter()
	subui.HandleFunc("/jobs", a.dashboardJobsHandler).Methods("GET")
	subui.HandleFunc("/jobs/{job}/executions", a.dashboardExecutionsHandler).Methods("GET")

	// Path of static files must be last!
	r.PathPrefix("/dashboard").Handler(
		http.StripPrefix("/dashboard", http.FileServer(
			http.Dir(filepath.Join(a.config.UIDir, "static")))))
	r.PathPrefix("/").Handler(http.RedirectHandler("dashboard", 301))
}

func templateSet(uiDir string, template string) []string {
	return []string{
		filepath.Join(uiDir, tmplPath, "dashboard.html.tmpl"),
		filepath.Join(uiDir, tmplPath, "status.html.tmpl"),
		filepath.Join(uiDir, tmplPath, template+".html.tmpl"),
	}
}

func (a *AgentCommand) dashboardIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl := template.Must(template.New("dashboard.html.tmpl").ParseFiles(
		templateSet(a.config.UIDir, "index")...))

	data := struct {
		Common *commonDashboardData
	}{
		Common: newCommonDashboardData(a, a.config.NodeName, ""),
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
		"jobJson": func(job *Job) string {
			j, _ := json.MarshalIndent(job, "", "<br>")
			return string(j)
		},
	}

	tmpl := template.Must(template.New("dashboard.html.tmpl").Funcs(funcs).ParseFiles(
		templateSet(a.config.UIDir, "jobs")...))

	data := struct {
		Common *commonDashboardData
		Jobs   []*Job
	}{
		Common: newCommonDashboardData(a, a.config.NodeName, "../"),
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
	groups := make(map[string][]*Execution)
	for _, exec := range execs {
		groups[exec.Group.String()] = append(groups[exec.Group.String()], exec)
	}

	tmpl := template.Must(template.New("dashboard.html.tmpl").Funcs(template.FuncMap{
		"b2s": func(value []byte) string {
			return template.JSEscapeString(string(value))
		},
		// Now unicode compliant
		"truncate": func(s string) string {
			var numRunes = 0
			for index, _ := range s {
				numRunes++
				if numRunes > 25 {
					return s[:index]
				}
			}
			return s
		},
	}).ParseFiles(templateSet(a.config.UIDir, "executions")...))

	if len(execs) > 100 {
		execs = execs[len(execs)-100:]
	}

	data := struct {
		Common  *commonDashboardData
		Groups  map[string][]*Execution
		JobName string
	}{
		Common:  newCommonDashboardData(a, a.config.NodeName, "../../../"),
		Groups:  groups,
		JobName: job,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Error(err)
	}
}
