package dkron

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/abronan/valkeyrie/store"
	"github.com/gin-contrib/expvar"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/victorcoder/dkron/plugintypes"
)

const (
	pretty = "pretty"
)

// Transport is the interface that wraps the ServeHTTP method.
type Transport interface {
	ServeHTTP()
}

// HTTPTransport stores pointers to an agent and a gin Engine.
type HTTPTransport struct {
	Engine *gin.Engine

	agent *Agent
}

// NewTransport creates an HTTPTransport with a bound agent.
func NewTransport(a *Agent) *HTTPTransport {
	return &HTTPTransport{
		agent: a,
	}
}

func (h *HTTPTransport) ServeHTTP() {
	h.Engine = gin.Default()
	h.Engine.HTMLRender = CreateMyRender()
	rootPath := h.Engine.Group("/")

	h.APIRoutes(rootPath)
	h.agent.DashboardRoutes(rootPath)

	h.Engine.Use(h.MetaMiddleware())

	log.WithFields(logrus.Fields{
		"address": h.agent.config.HTTPAddr,
	}).Info("api: Running HTTP server")

	go h.Engine.Run(h.agent.config.HTTPAddr)
}

// APIRoutes registers the api routes on the gin RouterGroup.
func (h *HTTPTransport) APIRoutes(r *gin.RouterGroup) {
	r.GET("/debug/vars", expvar.Handler())

	h.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	r.GET("/v1", h.indexHandler)
	v1 := r.Group("/v1")
	v1.GET("/", h.indexHandler)
	v1.GET("/members", h.membersHandler)
	v1.GET("/leader", h.leaderHandler)
	v1.POST("/leave", h.leaveHandler)

	v1.POST("/jobs", h.jobCreateOrUpdateHandler)
	v1.PATCH("/jobs", h.jobCreateOrUpdateHandler)
	// Place fallback routes last
	v1.GET("/jobs", h.jobsHandler)

	jobs := v1.Group("/jobs")
	jobs.DELETE("/:job", h.jobDeleteHandler)
	jobs.POST("/:job", h.jobRunHandler)
	jobs.POST("/:job/toggle", h.jobToggleHandler)

	// Place fallback routes last
	jobs.GET("/:job", h.jobGetHandler)
	jobs.GET("/:job/executions", h.executionsHandler)
}

// MetaMiddleware adds middleware to the gin Context.
func (h *HTTPTransport) MetaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Whom", h.agent.config.NodeName)
		c.Next()
	}
}

func renderJSON(c *gin.Context, status int, v interface{}) {
	if _, ok := c.GetQuery(pretty); ok {
		c.IndentedJSON(status, v)
	} else {
		c.JSON(status, v)
	}
}

func (h *HTTPTransport) indexHandler(c *gin.Context) {
	local := h.agent.serf.LocalMember()

	var status int
	if err := h.agent.Store.Healthy(); err != nil {
		status = http.StatusServiceUnavailable
	} else {
		status = http.StatusOK
	}

	stats := map[string]map[string]string{
		"agent": {
			"name":           local.Name,
			"version":        Version,
			"backend":        string(h.agent.config.Backend),
			"backend_status": strconv.FormatInt(int64(status), 10),
		},
		"serf": h.agent.serf.Stats(),
		"tags": local.Tags,
	}

	renderJSON(c, status, stats)
}

func (h *HTTPTransport) jobsHandler(c *gin.Context) {
	jobTags := c.QueryMap("tags")

	jobs, err := h.agent.Store.GetJobs(&JobOptions{ComputeStatus: true, Tags: jobTags})
	if err != nil {
		log.WithError(err).Error("api: Unable to get jobs, store not reachable.")
		return
	}
	renderJSON(c, http.StatusOK, jobs)
}

func (h *HTTPTransport) jobGetHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := h.agent.Store.GetJob(jobName, &JobOptions{ComputeStatus: true})
	if err != nil {
		log.Error(err)
	}
	if job == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	renderJSON(c, http.StatusOK, job)
}

func (h *HTTPTransport) jobCreateOrUpdateHandler(c *gin.Context) {
	// Init the Job object with defaults
	job := Job{
		Concurrency: ConcurrencyAllow,
	}

	// Parse values from JSON
	if err := c.BindJSON(&job); err != nil {
		c.Writer.WriteString("Incorrect or unexpected parameters")
		log.Error(err)
		return
	}

	// Save the job to the store
	if err := h.agent.Store.SetJob(&job, true); err != nil {
		c.AbortWithError(422, err)
		return
	}

	h.agent.SchedulerRestart()

	c.Header("Location", fmt.Sprintf("%s/%s", c.Request.RequestURI, job.Name))
	renderJSON(c, http.StatusCreated, &job)
}

func (h *HTTPTransport) jobDeleteHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := h.agent.Store.DeleteJob(jobName)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	h.agent.SchedulerRestart()
	renderJSON(c, http.StatusOK, job)
}

func (h *HTTPTransport) jobRunHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := h.agent.Store.GetJob(jobName, nil)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	ex := plugintypes.NewExecution(job.Name)
	h.agent.RunQuery(ex)

	c.Header("Location", c.Request.RequestURI)
	c.Status(http.StatusAccepted)
	renderJSON(c, http.StatusOK, job)
}

func (h *HTTPTransport) executionsHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := h.agent.Store.GetJob(jobName, nil)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	executions, err := h.agent.Store.GetExecutions(job.Name)
	if err != nil {
		if err == store.ErrKeyNotFound {
			renderJSON(c, http.StatusOK, &[]plugintypes.Execution{})
			return
		}
		log.Error(err)
		return

	}
	renderJSON(c, http.StatusOK, executions)
}

func (h *HTTPTransport) membersHandler(c *gin.Context) {
	renderJSON(c, http.StatusOK, h.agent.serf.Members())
}

func (h *HTTPTransport) leaderHandler(c *gin.Context) {
	member, err := h.agent.leaderMember()
	if err == nil {
		renderJSON(c, http.StatusOK, member)
	}
}

func (h *HTTPTransport) leaveHandler(c *gin.Context) {
	if err := h.agent.Stop(); err != nil {
		renderJSON(c, http.StatusOK, h.agent.ListServers())
	}
}

func (h *HTTPTransport) jobToggleHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := h.agent.Store.GetJob(jobName, nil)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	job.Disabled = !job.Disabled
	if err := h.agent.Store.SetJob(job, false); err != nil {
		c.AbortWithError(http.StatusPreconditionFailed, err)
		return
	}

	h.agent.SchedulerRestart()
	c.Header("Location", c.Request.RequestURI)
	renderJSON(c, http.StatusOK, job)
}
