package dkron

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"gopkg.in/gin-gonic/gin.v1"
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
	Keyspace   string
}

func newCommonDashboardData(a *AgentCommand, nodeName, path string) *commonDashboardData {
	leaderName := ""
	l, err := a.leaderMember()
	if err == nil {
		leaderName = l.Name
	}

	return &commonDashboardData{
		Version:    a.Version,
		LeaderName: leaderName,
		MemberName: nodeName,
		Backend:    a.config.Backend,
		Path:       fmt.Sprintf("%s%s", path, dashboardPathPrefix),
		APIPath:    fmt.Sprintf("%s%s", path, apiPathPrefix),
		Keyspace:   a.config.Keyspace,
	}
}

func (a *AgentCommand) dashboardRoutes(r *gin.Engine) {
	r.LoadHTMLGlob(filepath.Join(a.config.UIDir, tmplPath, "*"))

	dashboard := r.Group("/" + dashboardPathPrefix)
	dashboard.GET("/", a.dashboardIndexHandler)
	// dashboard.GET("/jobs", a.dashboardJobsHandler)
	// dashboard.GET("/jobs/:job/executions", a.dashboardExecutionsHandler)
}

func templateSet(uiDir string, template string) []string {
	return []string{
		filepath.Join(uiDir, tmplPath, "dashboard.html.tmpl"),
		filepath.Join(uiDir, tmplPath, "status.html.tmpl"),
		filepath.Join(uiDir, tmplPath, template+".html.tmpl"),
	}
}

func (a *AgentCommand) dashboardIndexHandler(c *gin.Context) {
	data := struct {
		Common *commonDashboardData
	}{
		Common: newCommonDashboardData(a, a.config.NodeName, ""),
	}
	c.HTML(http.StatusOK, "dashboard.html.tmpl", data)
}

func (a *AgentCommand) dashboardJobsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	jobs, _ := a.store.GetJobs()

	funcs := template.FuncMap{
		"executionStatus": func(job *Job) string {
			status := job.Status()
			switch status {
			case Success:
				return "success"
			case Failed:
				return "danger"
			case PartialyFailed:
				return "warning"
			case Running:
				return ""
			}

			return ""
		},
		"jobJson": func(job *Job) string {
			j, _ := json.MarshalIndent(job, "", "\t")
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

	groups, byGroup, err := a.store.GetGroupedExecutions(job)
	if err != nil {
		log.Error(err)
	}

	tmpl := template.Must(template.New("dashboard.html.tmpl").Funcs(template.FuncMap{
		"html": func(value []byte) string {
			return string(template.HTML(value))
		},
		// Now unicode compliant
		"truncate": func(s string) string {
			var numRunes = 0
			for index := range s {
				numRunes++
				if numRunes > 25 {
					return s[:index]
				}
			}
			return s
		},
	}).ParseFiles(templateSet(a.config.UIDir, "executions")...))

	data := struct {
		Common  *commonDashboardData
		Groups  map[int64][]*Execution
		JobName string
		ByGroup int64arr
	}{
		Common:  newCommonDashboardData(a, a.config.NodeName, "../../../"),
		Groups:  groups,
		JobName: job,
		ByGroup: byGroup,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Error(err)
	}
}
