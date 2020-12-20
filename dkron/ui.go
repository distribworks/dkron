package dkron

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/distribworks/dkron/v3/dkron/assets_ui"
	"github.com/gin-gonic/gin"
)

const uiPathPrefix = "ui/"

// DashboardRoutes registers dashboard specific routes on the gin RouterGroup.
func (h *HTTPTransport) UI(r *gin.RouterGroup) {
	ui := r.Group("/" + uiPathPrefix)

	a, err := assets_ui.Assets.Open("index.html")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.New("index.html").Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	h.Engine.SetHTMLTemplate(t)

	ui.GET("/*filepath", func(ctx *gin.Context) {
		p := ctx.Param("filepath")
		_, err := assets_ui.Assets.Open(p)
		if err == nil && p != "/" && p != "/index.html" {
			ctx.FileFromFS(p, assets_ui.Assets)
		} else {
			jobs, err := h.agent.Store.GetJobs(nil)
			if err != nil {
				log.Error(err)
			}
			totalJobs := len(jobs)
			successfulJobs := 0
			failedJobs := 0
			for _, j := range jobs {
				if j.Status == "success" {
					successfulJobs++
				} else {
					failedJobs++
				}
			}
			l, err := h.agent.leaderMember()
			ln := "no leader"
			if err != nil {
				log.Error(err)
			} else {
				ln = l.Name
			}
			ctx.HTML(http.StatusOK, "index.html", gin.H{
				"DKRON_API_URL":         fmt.Sprintf("/%s", apiPathPrefix),
				"DKRON_LEADER":          ln,
				"DKRON_TOTAL_JOBS":      totalJobs,
				"DKRON_FAILED_JOBS":     failedJobs,
				"DKRON_SUCCESSFUL_JOBS": successfulJobs,
			})
		}
	})
}
