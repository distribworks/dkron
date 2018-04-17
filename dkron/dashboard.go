package dkron

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

const (
	tmplPath            = "templates"
	dashboardPathPrefix = "dashboard"
	assetsPrefix        = "static"
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

func newCommonDashboardData(a *Agent, nodeName, path string) *commonDashboardData {
	leaderName := ""
	l, err := a.leaderMember()
	if err == nil {
		leaderName = l.Name
	}

	return &commonDashboardData{
		Version:    Version,
		LeaderName: leaderName,
		MemberName: nodeName,
		Backend:    a.config.Backend,
		AssetsPath: fmt.Sprintf("%s%s", path, assetsPrefix),
		Path:       fmt.Sprintf("%s%s", path, dashboardPathPrefix),
		APIPath:    fmt.Sprintf("%s%s", path, apiPathPrefix),
		Keyspace:   a.config.Keyspace,
	}
}

// dashboardRoutes registers dashboard specific routes on the gin RouterGroup.
func (a *Agent) DashboardRoutes(r *gin.RouterGroup) {
	r.GET("/static/*asset", servePublic)

	dashboard := r.Group("/" + dashboardPathPrefix)
	dashboard.GET("/", a.dashboardIndexHandler)
	dashboard.GET("/jobs", a.dashboardJobsHandler)
	dashboard.GET("/jobs/:job/executions", a.dashboardExecutionsHandler)
}

func (a *Agent) dashboardIndexHandler(c *gin.Context) {
	data := struct {
		Common *commonDashboardData
	}{
		Common: newCommonDashboardData(a, a.config.NodeName, "../"),
	}
	c.HTML(http.StatusOK, "index", data)
}

func (a *Agent) dashboardJobsHandler(c *gin.Context) {
	jobs, _ := a.Store.GetJobs()

	data := struct {
		Common *commonDashboardData
		Jobs   []*Job
	}{
		Common: newCommonDashboardData(a, a.config.NodeName, "../../"),
		Jobs:   jobs,
	}

	c.HTML(http.StatusOK, "jobs", data)
}

func (a *Agent) dashboardExecutionsHandler(c *gin.Context) {
	job := c.Param("job")

	groups, byGroup, err := a.Store.GetGroupedExecutions(job)
	if err != nil {
		log.Error(err)
	}

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

	c.HTML(http.StatusOK, "executions", data)
}

func mustLoadTemplate(path string) []byte {
	tmpl, err := Asset(path)
	if err != nil {
		log.Error(err)
		return nil
	}

	return tmpl
}

func CreateMyRender() multitemplate.Render {
	r := multitemplate.New()

	status := mustLoadTemplate(tmplPath + "/status.html.tmpl")
	dash := mustLoadTemplate(tmplPath + "/dashboard.html.tmpl")

	r.AddFromStringsFuncs("index", funcMap(),
		string(dash),
		string(status),
		string(mustLoadTemplate(tmplPath+"/index.html.tmpl")))

	r.AddFromStringsFuncs("jobs", funcMap(),
		string(dash),
		string(status),
		string(mustLoadTemplate(tmplPath+"/jobs.html.tmpl")))

	r.AddFromStringsFuncs("executions", funcMap(),
		string(dash),
		string(status),
		string(mustLoadTemplate(tmplPath+"/executions.html.tmpl")))

	return r
}

//go:generate go-bindata -prefix "../" -pkg dkron -ignore=scss -ignore=.*\.md -ignore=\.?bower\.json -ignore=\.gitignore -ignore=Makefile -ignore=examples -ignore=tutorial -ignore=tests -ignore=rickshaw\/src -o bindata.go ../static/... ../templates
func servePublic(c *gin.Context) {
	path := c.Request.URL.Path

	path = strings.Replace(path, "/", "", 1)
	split := strings.Split(path, ".")
	suffix := split[len(split)-1]

	res, err := Asset(path)
	if err != nil {
		c.Next()
		return
	}

	contentType := "text/plain"
	switch suffix {
	case "png":
		contentType = "image/png"
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "gif":
		contentType = "image/gif"
	case "js":
		contentType = "application/javascript"
	case "css":
		contentType = "text/css"
	case "woff":
		contentType = "application/x-font-woff"
	case "ttf":
		contentType = "application/x-font-ttf"
	case "otf":
		contentType = "application/x-font-otf"
	case "html":
		contentType = "text/html"
	}

	c.Writer.Header().Set("content-type", contentType)
	c.String(200, string(res))
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
		"toString": func(value []byte) string {
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
