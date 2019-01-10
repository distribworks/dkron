package dkron

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/victorcoder/dkron/static"
	"github.com/victorcoder/dkron/templates"
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
	Name       string
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
		Backend:    string(a.config.Backend),
		AssetsPath: fmt.Sprintf("%s%s", path, assetsPrefix),
		Path:       fmt.Sprintf("%s%s", path, dashboardPathPrefix),
		APIPath:    fmt.Sprintf("%s%s", path, apiPathPrefix),
		Keyspace:   a.config.Keyspace,
		Name:       Name,
	}
}

// dashboardRoutes registers dashboard specific routes on the gin RouterGroup.
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

	r.StaticFS("static", static.Assets)
	//r.GET("/static/*asset", servePublic)

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

//var Assets http.FileSystem = http.Dir("../static")
//var Templates http.FileSystem = http.Dir("../templates")

//go:generate vfsgendev -source="github.com/victorcoder/dkron/static".Assets
//go:generate vfsgendev -source="github.com/victorcoder/dkron/templates".Templates
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
