package dkron

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const uiPathPrefix = "ui/"

//go:embed ui-dist
var uiDist embed.FS

// UI registers UI specific routes on the gin RouterGroup.
func (h *HTTPTransport) UI(r *gin.RouterGroup) {
	// If we are visiting from a browser redirect to the dashboard
	r.GET("/", func(c *gin.Context) {
		switch c.NegotiateFormat(gin.MIMEHTML) {
		case gin.MIMEHTML:
			c.Redirect(http.StatusSeeOther, "/ui/")
		default:
			c.AbortWithStatus(http.StatusNotFound)
		}
	})

	ui := r.Group("/" + uiPathPrefix)

	assets, err := fs.Sub(uiDist, "ui-dist")
	if err != nil {
		h.logger.Fatal(err)
	}
	a, err := assets.Open("index.html")
	if err != nil {
		h.logger.Fatal(err)
	}
	b, err := ioutil.ReadAll(a)
	if err != nil {
		h.logger.Fatal(err)
	}
	t, err := template.New("index.html").Parse(string(b))
	if err != nil {
		h.logger.Fatal(err)
	}
	h.Engine.SetHTMLTemplate(t)

	ui.GET("/*filepath", func(ctx *gin.Context) {
		p := ctx.Param("filepath")
		f := strings.TrimPrefix(p, "/")
		_, err := assets.Open(f)
		if err == nil && p != "/" && p != "/index.html" {
			ctx.FileFromFS(p, http.FS(assets))
		} else {
			jobs, err := h.agent.Store.GetJobs(nil)
			if err != nil {
				h.logger.Error(err)
			}
			var (
				totalJobs                                   = len(jobs)
				successfulJobs, failedJobs, untriggeredJobs int
			)
			for _, j := range jobs {
				if j.Status == "success" {
					successfulJobs++
				} else if j.Status == "failed" {
					failedJobs++
				} else if j.Status == "" {
					untriggeredJobs++
				}
			}
			l, err := h.agent.leaderMember()
			ln := "no leader"
			if err != nil {
				h.logger.Error(err)
			} else {
				ln = l.Name
			}
			ctx.HTML(http.StatusOK, "index.html", gin.H{
				"DKRON_API_URL":          fmt.Sprintf("/%s", apiPathPrefix),
				"DKRON_LEADER":           ln,
				"DKRON_TOTAL_JOBS":       totalJobs,
				"DKRON_FAILED_JOBS":      failedJobs,
				"DKRON_UNTRIGGERED_JOBS": untriggeredJobs,
				"DKRON_SUCCESSFUL_JOBS":  successfulJobs,
			})
		}
	})
}
