package dkron

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/distribworks/dkron/v3/dkron/assets"
	"github.com/distribworks/dkron/v3/dkron/templates"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

const (
	dashboardPathPrefix = "dashboard"
	assetsPrefix        = "static"
	apiPathPrefix       = "v1"
)

type commonDashboardData struct {
	Version    string
	LeaderName string
	MemberName string
	AssetsPath string
	Path       string
	APIPath    string
	Name       string
	Time       string
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
		AssetsPath: fmt.Sprintf("%s%s", path, assetsPrefix),
		Path:       fmt.Sprintf("%s%s", path, dashboardPathPrefix),
		APIPath:    fmt.Sprintf("%s%s", path, apiPathPrefix),
		Name:       Name,
		Time:       time.Now().UTC().Format("15:04:05 UTC"),
	}
}

// DashboardRoutes registers dashboard specific routes on the gin RouterGroup.
func (a *Agent) DashboardRoutes(r *gin.RouterGroup) {
	// If we are visiting from a browser redirect to the dashboard
	r.GET("/", func(c *gin.Context) {
		switch c.NegotiateFormat(gin.MIMEHTML) {
		case gin.MIMEHTML:
			c.Redirect(http.StatusMovedPermanently, "/"+dashboardPathPrefix+"/")
		default:
			c.AbortWithStatus(http.StatusNotFound)
		}
	})

	r.StaticFS("static", assets.Assets)

	dashboard := r.Group("/" + dashboardPathPrefix)
	dashboard.GET("/", a.dashboardIndexHandler)
	dashboard.GET("/jobs", a.dashboardJobsHandler)
	dashboard.GET("/jobs/:job/executions", a.dashboardExecutionsHandler)
	dashboard.GET("/busy", a.dashboardBusyHandler)
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
	data := struct {
		Common *commonDashboardData
	}{
		Common: newCommonDashboardData(a, a.config.NodeName, "../"),
	}

	c.HTML(http.StatusOK, "jobs", data)
}

func (a *Agent) dashboardExecutionsHandler(c *gin.Context) {
	job := c.Param("job")

	groups, byGroup, err := a.Store.GetGroupedExecutions(job)
	if err != nil {
		log.Error(err)
	}

	var count int
	for _, v := range groups {
		count = count + len(v)
	}

	data := struct {
		Common  *commonDashboardData
		Groups  map[int64][]*Execution
		JobName string
		ByGroup int64arr
		Count   int
	}{
		Common:  newCommonDashboardData(a, a.config.NodeName, "../../../"),
		Groups:  groups,
		JobName: job,
		ByGroup: byGroup,
		Count:   count,
	}

	c.HTML(http.StatusOK, "executions", data)
}

func (a *Agent) dashboardBusyHandler(c *gin.Context) {
	data := struct {
		Common *commonDashboardData
	}{
		Common: newCommonDashboardData(a, a.config.NodeName, "../"),
	}
	c.HTML(http.StatusOK, "busy", data)
}

func mustLoadTemplate(path string) []byte {
	f, err := templates.Templates.Open(path)
	if err != nil {
		log.Error(err)
		return nil
	}

	tmpl, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error(err)
		return nil
	}

	return tmpl
}

// CreateMyRender returns a new custom multitemplate renderer
// to use with Gin.
func CreateMyRender() multitemplate.Render {
	r := multitemplate.New()

	status := mustLoadTemplate("/status.html.tmpl")
	dash := mustLoadTemplate("/dashboard.html.tmpl")

	r.AddFromStringsFuncs("index", funcMap(),
		string(dash),
		string(status),
		string(mustLoadTemplate("/index.html.tmpl")))

	r.AddFromStringsFuncs("jobs", funcMap(),
		string(dash),
		string(status),
		string(mustLoadTemplate("/jobs.html.tmpl")))

	r.AddFromStringsFuncs("executions", funcMap(),
		string(dash),
		string(status),
		string(mustLoadTemplate("/executions.html.tmpl")))

	r.AddFromStringsFuncs("busy", funcMap(),
		string(dash),
		string(status),
		string(mustLoadTemplate("/busy.html.tmpl")))

	return r
}

func funcMap() template.FuncMap {
	funcs := template.FuncMap{
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
