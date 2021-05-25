package dkron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/expvar"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/serf/serf"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	status "google.golang.org/grpc/status"
)

const (
	pretty        = "pretty"
	apiPathPrefix = "v1"
)

// Transport is the interface that wraps the ServeHTTP method.
type Transport interface {
	ServeHTTP()
}

// HTTPTransport stores pointers to an agent and a gin Engine.
type HTTPTransport struct {
	Engine *gin.Engine

	agent  *Agent
	logger *logrus.Entry
}

// NewTransport creates an HTTPTransport with a bound agent.
func NewTransport(a *Agent, log *logrus.Entry) *HTTPTransport {
	return &HTTPTransport{
		agent:  a,
		logger: log,
	}
}

func (h *HTTPTransport) ServeHTTP() {
	h.Engine = gin.Default()
	h.Engine.HTMLRender = CreateMyRender(h.logger)
	h.Engine.Use(h.Options)

	rootPath := h.Engine.Group("/")

	rootPath.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	rootPath.Use(h.MetaMiddleware())

	h.APIRoutes(rootPath)
	if h.agent.config.UI {
		h.UI(rootPath)
	} else {
		h.agent.DashboardRoutes(rootPath)
	}

	h.logger.WithFields(logrus.Fields{
		"address": h.agent.config.HTTPAddr,
	}).Info("api: Running HTTP server")

	go h.Engine.Run(h.agent.config.HTTPAddr)
}

// APIRoutes registers the api routes on the gin RouterGroup.
func (h *HTTPTransport) APIRoutes(r *gin.RouterGroup, middleware ...gin.HandlerFunc) {
	r.GET("/debug/vars", expvar.Handler())

	h.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	if h.agent.config.EnablePrometheus {
		// Prometheus metrics scrape endpoint
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	r.GET("/v1", h.indexHandler)
	v1 := r.Group("/v1")
	v1.Use(middleware...)
	v1.GET("/", h.indexHandler)
	v1.GET("/members", h.membersHandler)
	v1.GET("/leader", h.leaderHandler)
	v1.GET("/isleader", h.isLeaderHandler)
	v1.POST("/leave", h.leaveHandler)
	v1.POST("/restore", h.restoreHandler)

	v1.GET("/busy", h.busyHandler)

	v1.POST("/jobs", h.jobCreateOrUpdateHandler)
	v1.PATCH("/jobs", h.jobCreateOrUpdateHandler)
	// Place fallback routes last
	v1.GET("/jobs", h.jobsHandler)

	jobs := v1.Group("/jobs")
	jobs.DELETE("/:job", h.jobDeleteHandler)
	jobs.POST("/:job", h.jobRunHandler)
	jobs.POST("/:job/run", h.jobRunHandler)
	jobs.POST("/:job/toggle", h.jobToggleHandler)
	jobs.PUT("/:job", h.jobCreateOrUpdateHandler)

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

func (h *HTTPTransport) Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		gh := cors.Default()
		gh(c)

		c.AbortWithStatus(http.StatusOK)
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

	stats := map[string]map[string]string{
		"agent": {
			"name":    local.Name,
			"version": Version,
		},
		"serf": h.agent.serf.Stats(),
		"tags": local.Tags,
	}

	renderJSON(c, http.StatusOK, stats)
}

func (h *HTTPTransport) jobsHandler(c *gin.Context) {
	metadata := c.QueryMap("metadata")
	sort := c.DefaultQuery("_sort", "id")
	if sort == "id" {
		sort = "name"
	}
	order := c.DefaultQuery("_order", "ASC")
	q := c.Query("q")

	jobs, err := h.agent.Store.GetJobs(
		&JobOptions{
			Metadata: metadata,
			Sort:     sort,
			Order:    order,
			Query:    q,
			Status:   c.Query("status"),
			Disabled: c.Query("disabled"),
		},
	)
	if err != nil {
		h.logger.WithError(err).Error("api: Unable to get jobs, store not reachable.")
		return
	}

	start, ok := c.GetQuery("_start")
	if !ok {
		start = "0"
	}
	s, _ := strconv.Atoi(start)

	end, ok := c.GetQuery("_end")
	e := 0
	if !ok {
		e = len(jobs)
	} else {
		e, _ = strconv.Atoi(end)
		if e > len(jobs) {
			e = len(jobs)
		}
	}

	c.Header("X-Total-Count", strconv.Itoa(len(jobs)))
	renderJSON(c, http.StatusOK, jobs[s:e])
}

func (h *HTTPTransport) jobGetHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := h.agent.Store.GetJob(jobName, nil)
	if err != nil {
		h.logger.Error(err)
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
		c.Writer.WriteString(fmt.Sprintf("Unable to parse payload: %s.", err))
		h.logger.Error(err)
		return
	}

	// Validate job
	if err := job.Validate(); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		c.Writer.WriteString(fmt.Sprintf("Job contains invalid value: %s.", err))
		return
	}

	// Call gRPC SetJob
	if err := h.agent.GRPCClient.SetJob(&job); err != nil {
		s := status.Convert(err)
		if s.Message() == ErrParentJobNotFound.Error() {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.Writer.WriteString(s.Message())
		return
	}

	// Immediately run the job if so requested
	if _, exists := c.GetQuery("runoncreate"); exists {
		h.agent.GRPCClient.RunJob(job.Name)
	}

	c.Header("Location", fmt.Sprintf("%s/%s", c.Request.RequestURI, job.Name))
	renderJSON(c, http.StatusCreated, &job)
}

func (h *HTTPTransport) jobDeleteHandler(c *gin.Context) {
	jobName := c.Param("job")

	// Call gRPC DeleteJob
	job, err := h.agent.GRPCClient.DeleteJob(jobName)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	renderJSON(c, http.StatusOK, job)
}

func (h *HTTPTransport) jobRunHandler(c *gin.Context) {
	jobName := c.Param("job")

	// Call gRPC RunJob
	job, err := h.agent.GRPCClient.RunJob(jobName)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Header("Location", c.Request.RequestURI)
	c.Status(http.StatusAccepted)
	renderJSON(c, http.StatusOK, job)
}

// Restore jobs from file.
// Overwrite job if the job is exist.
func (h *HTTPTransport) restoreHandler(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var jobs []*Job
	err = json.Unmarshal(data, &jobs)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	jobTree, err := generateJobTree(jobs)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	result := h.agent.recursiveSetJob(jobTree)
	resp, err := json.Marshal(result)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	renderJSON(c, http.StatusOK, string(resp))
}

func (h *HTTPTransport) executionsHandler(c *gin.Context) {
	jobName := c.Param("job")

	sort := c.DefaultQuery("_sort", "")
	if sort == "id" {
		sort = "started_at"
	}
	order := c.DefaultQuery("_order", "DESC")

	job, err := h.agent.Store.GetJob(jobName, nil)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	executions, err := h.agent.Store.GetExecutions(job.Name,
		&ExecutionOptions{
			Sort:     sort,
			Order:    order,
			Timezone: job.GetTimeLocation(),
		},
	)
	if err == buntdb.ErrNotFound {
		executions = make([]*Execution, 0)
	} else if err != nil {
		h.logger.Error(err)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(len(executions)))
	renderJSON(c, http.StatusOK, executions)
}

type MId struct {
	serf.Member

	Id string `json:"id"`
}

func (h *HTTPTransport) membersHandler(c *gin.Context) {
	mems := []*MId{}
	for _, m := range h.agent.serf.Members() {
		id, _ := uuid.GenerateUUID()
		mid := &MId{m, id}
		mems = append(mems, mid)
	}
	c.Header("X-Total-Count", strconv.Itoa(len(mems)))
	renderJSON(c, http.StatusOK, mems)
}

func (h *HTTPTransport) leaderHandler(c *gin.Context) {
	member, err := h.agent.leaderMember()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	if member == nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	renderJSON(c, http.StatusOK, member)
}

func (h *HTTPTransport) isLeaderHandler(c *gin.Context) {
	isleader := h.agent.IsLeader()
	if isleader {
		renderJSON(c, http.StatusOK, "I am a leader")
	} else {
		renderJSON(c, http.StatusNotFound, "I am a follower")
	}
}

func (h *HTTPTransport) leaveHandler(c *gin.Context) {
	if err := h.agent.Stop(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	renderJSON(c, http.StatusOK, h.agent.peers)
}

func (h *HTTPTransport) jobToggleHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := h.agent.Store.GetJob(jobName, nil)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// Toggle job status
	job.Disabled = !job.Disabled

	// Call gRPC SetJob
	if err := h.agent.GRPCClient.SetJob(job); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	c.Header("Location", c.Request.RequestURI)
	renderJSON(c, http.StatusOK, job)
}

func (h *HTTPTransport) busyHandler(c *gin.Context) {
	executions := []*Execution{}

	exs, err := h.agent.GetActiveExecutions()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, e := range exs {
		executions = append(executions, NewExecutionFromProto(e))
	}

	sort.SliceStable(executions, func(i, j int) bool {
		return executions[i].StartedAt.Before(executions[j].StartedAt)
	})

	c.Header("X-Total-Count", strconv.Itoa(len(executions)))
	renderJSON(c, http.StatusOK, executions)
}
