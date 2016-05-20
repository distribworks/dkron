package dkron

import (
	"expvar"
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/victorcoder/dkron/cron"
)

var cronInspect = expvar.NewMap("cron_entries")
var schedulerStarted = expvar.NewString("scheduler_started")

type Scheduler struct {
	Cron    *cron.Cron
	Started bool
}

func NewScheduler() *Scheduler {
	c := cron.New()
	schedulerStarted.Set("false")
	return &Scheduler{Cron: c, Started: false}
}

func (s *Scheduler) Start(jobs []*Job) {
	for _, job := range jobs {
		log.WithFields(logrus.Fields{
			"job": job.Name,
		}).Debug("scheduler: Adding job to cron")

		s.Cron.AddJob(job.Schedule, job)
		cronInspect.Set(job.Name, job)
	}
	s.Cron.Start()
	s.Started = true
	schedulerStarted.Set("true")
}

func (s *Scheduler) Restart(jobs []*Job) {
	s.Cron.Stop()
	s.Cron = cron.New()
	s.Start(jobs)
}

func (s *Scheduler) GetEntry(job *Job) *cron.Entry {
	for _, e := range s.Cron.Entries() {
		j, _ := e.Job.(*Job)
		if j.Name == job.Name {
			return e
		}
	}
	return nil
}

type Job struct {
	// Job name. Must be unique, acts as the id.
	Name string `json:"name"`

	// Cron expression for the job. When to run the job.
	Schedule string `json:"schedule"`

	// Use shell to run the command.
	Shell bool `json:"shell"`

	// Command to run. Must be a shell command to execute.
	Command string `json:"command"`

	// Owner of the job.
	Owner string `json:"owner"`

	// Owner email of the job.
	OwnerEmail string `json:"owner_email"`

	// Actual user to use when running the command.
	RunAsUser string `json:"run_as_user"`

	// Number of successful executions of this job.
	SuccessCount int `json:"success_count"`

	// Number of errors running this job.
	ErrorCount int `json:"error_count"`

	// Last time this job executed succesful.
	LastSuccess time.Time `json:"last_success"`

	// Last time this job failed.
	LastError time.Time `json:"last_error"`

	// Is this job disabled?
	Disabled bool `json:"disabled"`

	// Tags of the target servers to run this job against.
	Tags map[string]string `json:"tags"`

	// Pointer to the calling agent.
	Agent *AgentCommand `json:"-"`

	running sync.Mutex
}

// Run the job
func (j *Job) Run() {
	j.running.Lock()
	defer j.running.Unlock()

	// Maybe we are testing
	if j.Agent != nil && j.Disabled == false {

		log.WithFields(logrus.Fields{
			"job":      j.Name,
			"schedule": j.Schedule,
		}).Debug("scheduler: Run job")

		cronInspect.Set(j.Name, j)
		j.Agent.RunQuery(j)
	}
}

// Friendly format a job
func (j *Job) String() string {
	return fmt.Sprintf("\"Job: %s, scheduled at: %s, tags:%v\"", j.Name, j.Schedule, j.Tags)
}

type Execution struct {
	// Name of the job this executions refers to.
	JobName string `json:"job_name,omitempty"`

	// Start time of the execution.
	StartedAt time.Time `json:"started_at,omitempty"`

	// When the execution finished running.
	FinishedAt time.Time `json:"finished_at,omitempty"`

	// If this execution executed succesfully.
	Success bool `json:"success,omitempty"`

	// Partial output of the execution.
	Output []byte `json:"output,omitempty"`

	// Node name of the node that run this execution.
	NodeName string `json:"node_name,omitempty"`

	// Execution group to what this execution belongs to.
	Group int64 `json:"group,omitempty"`

	// The job used to generate this execution.
	Job *Job `json:"job,omitempty"`
}

// Used to enerate the execution Id
func (e *Execution) Key() string {
	return fmt.Sprintf("%d-%s", e.StartedAt.UnixNano(), e.NodeName)
}
