package dkron

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libkv/store"
	gin "gopkg.in/gin-gonic/gin.v1"
)

const pretty = "pretty"

func (a *AgentCommand) ServeHTTP() {
	r := gin.Default()
	if log.Level >= logrus.InfoLevel {
		gin.SetMode(gin.ReleaseMode)
	}

	a.apiRoutes(r)
	a.dashboardRoutes(r)

	r.Use(a.metaMiddleware())
	//r.GET("/debug/vars", expvar.Handler())

	log.WithFields(logrus.Fields{
		"address": a.config.HTTPAddr,
	}).Info("api: Running HTTP server")

	go r.Run(a.config.HTTPAddr)
}

func (a *AgentCommand) apiRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.GET("/", a.indexHandler)
	v1.GET("/members", a.membersHandler)
	v1.GET("/leader", a.leaderHandler)
	v1.POST("/leave", a.leaveHandler)

	v1.POST("/jobs", a.jobCreateOrUpdateHandler)
	v1.PATCH("/jobs", a.jobCreateOrUpdateHandler)
	// Place fallback routes last
	v1.GET("/jobs", a.jobsHandler)

	jobs := v1.Group("/jobs")
	jobs.DELETE("/:job", a.jobDeleteHandler)
	jobs.POST("/:job", a.jobRunHandler)
	// Place fallback routes last
	jobs.GET("/:job", a.jobGetHandler)
	jobs.GET("/:job/executions", a.executionsHandler)
}

func (a *AgentCommand) metaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Whom", a.config.NodeName)
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

func (a *AgentCommand) indexHandler(c *gin.Context) {
	local := a.serf.LocalMember()
	stats := map[string]map[string]string{
		"agent": {
			"name":    local.Name,
			"version": a.Version,
			"backend": a.config.Backend,
		},
		"serf": a.serf.Stats(),
		"tags": local.Tags,
	}
	renderJSON(c, http.StatusOK, stats)
}

func (a *AgentCommand) jobsHandler(c *gin.Context) {
	jobs, err := a.store.GetJobs()
	if err != nil {
		log.WithError(err).Error("api: Unable to get jobs, store not reachable.")
		return
	}
	renderJSON(c, http.StatusOK, jobs)
}

func (a *AgentCommand) jobGetHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := a.store.GetJob(jobName)
	if err != nil {
		log.Error(err)
	}
	if job == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	renderJSON(c, http.StatusOK, job)
}

func (a *AgentCommand) jobCreateOrUpdateHandler(c *gin.Context) {
	// Init the Job object with defaults
	job := Job{
		Concurrency: ConcurrencyAllow,
	}
	c.BindJSON(&job)

	// Get if the requested job already exist
	ej, err := a.store.GetJob(job.Name)
	if err != nil && err != store.ErrKeyNotFound {
		c.AbortWithError(422, err)
		return
	}

	// If it's an existing job, lock it
	if ej != nil {
		ej.Lock()
		defer ej.Unlock()
	}

	// Save the job to the store
	if err = a.store.SetJob(&job); err != nil {
		c.AbortWithError(422, err)
		return
	}

	// Save the job parent
	if err = a.store.SetJobDependencyTree(&job, ej); err != nil {
		c.AbortWithError(422, err)
		return
	}

	a.schedulerRestartQuery(string(a.store.GetLeader()))

	c.Header("Location", fmt.Sprintf("%s/%s", c.Request.RequestURI, job.Name))
	renderJSON(c, http.StatusCreated, job)
}

func (a *AgentCommand) jobDeleteHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := a.store.DeleteJob(jobName)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	a.schedulerRestartQuery(string(a.store.GetLeader()))
	renderJSON(c, http.StatusOK, job)
}

func (a *AgentCommand) jobRunHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := a.store.GetJob(jobName)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	ex := NewExecution(job.Name)
	a.RunQuery(ex)

	c.Header("Location", c.Request.RequestURI)
	c.Status(http.StatusAccepted)
	renderJSON(c, http.StatusOK, job)
}

func (a *AgentCommand) executionsHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := a.store.GetJob(jobName)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	executions, err := a.store.GetExecutions(job.Name)
	if err != nil {
		if err == store.ErrKeyNotFound {
			renderJSON(c, http.StatusOK, &[]Execution{})
			return
		} else {
			log.Error(err)
			return
		}
	}
	renderJSON(c, http.StatusOK, executions)
}

func (a *AgentCommand) membersHandler(c *gin.Context) {
	renderJSON(c, a.serf.Members())
}

func (a *AgentCommand) leaderHandler(c *gin.Context) {
	member, err := a.leaderMember()
	if err == nil {
		renderJSON(c, http.StatusOK, member)
	}
}

func (a *AgentCommand) leaveHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		log.Warn("/leave GET is deprecated and will be removed, use POST")
	}
	if err := a.serf.Leave(); err != nil {
		renderJSON(c, http.StatusOK, a.listServers())
	}
}
