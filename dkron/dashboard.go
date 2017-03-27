package dkron

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"gopkg.in/gin-gonic/gin.v1"
)

const (
	tmplPath            = "templates"
	dashboardPathPrefix = "dashboard"
	assetsPrefix        = "assets"
	apiPathPrefix       = "v1"
)

type commonDashboardData struct {
	Version    string
	LeaderName string
	MemberName string
	Backend    string
	AssetsPath string
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
		AssetsPath: fmt.Sprintf("%s%s", path, assetsPrefix),
		Path:       fmt.Sprintf("%s%s", path, dashboardPathPrefix),
		APIPath:    fmt.Sprintf("%s%s", path, apiPathPrefix),
		Keyspace:   a.config.Keyspace,
	}
}

func (a *AgentCommand) dashboardRoutes(r *gin.Engine) {
	r.HTMLRender = createMyRender(filepath.Join(a.config.UIDir, tmplPath))

	dashboard := r.Group("/" + dashboardPathPrefix)
	dashboard.GET("/", a.dashboardIndexHandler)
	dashboard.GET("/jobs", a.dashboardJobsHandler)
	dashboard.GET("/jobs/:job/executions", a.dashboardExecutionsHandler)
}

func (a *AgentCommand) dashboardIndexHandler(c *gin.Context) {
	data := struct {
		Common *commonDashboardData
	}{
		Common: newCommonDashboardData(a, a.config.NodeName, "/"),
	}
	c.HTML(http.StatusOK, "index", data)
}

func (a *AgentCommand) dashboardJobsHandler(c *gin.Context) {
	jobs, _ := a.store.GetJobs()

	data := struct {
		Common *commonDashboardData
		Jobs   []*Job
	}{
		Common: newCommonDashboardData(a, a.config.NodeName, "/"),
		Jobs:   jobs,
	}

	c.HTML(http.StatusOK, "jobs", data)
}

func (a *AgentCommand) dashboardExecutionsHandler(c *gin.Context) {
	job := c.Param("job")

	groups, byGroup, err := a.store.GetGroupedExecutions(job)
	if err != nil {
		log.Error(err)
	}

	data := struct {
		Common  *commonDashboardData
		Groups  map[int64][]*Execution
		JobName string
		ByGroup int64arr
	}{
		Common:  newCommonDashboardData(a, a.config.NodeName, "/"),
		Groups:  groups,
		JobName: job,
		ByGroup: byGroup,
	}

	c.HTML(http.StatusOK, "executions", data)
}

func createMyRender(path string) multitemplate.Render {
	r := multitemplate.New()

	r.AddFromFilesFuncs("index",
		funcMap(),
		filepath.Join(path, "dashboard.html.tmpl"),
		filepath.Join(path, "status.html.tmpl"),
		filepath.Join(path, "index.html.tmpl"))

	r.AddFromFilesFuncs("jobs",
		funcMap(),
		filepath.Join(path, "dashboard.html.tmpl"),
		filepath.Join(path, "status.html.tmpl"),
		filepath.Join(path, "jobs.html.tmpl"))

	r.AddFromFilesFuncs("executions",
		funcMap(),
		filepath.Join(path, "dashboard.html.tmpl"),
		filepath.Join(path, "status.html.tmpl"),
		filepath.Join(path, "executions.html.tmpl"))

	return r
}

func funcMap() template.FuncMap {
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
	}

	return funcs
}
