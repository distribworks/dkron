package dkron

import (
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

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

	// Number of times to retry a job that failed an execution.
	Retries uint `json:"retries"`

	running sync.Mutex

	// Jobs that are dependent upon this one will be run after this job runs.
	DependentJobs []string `json:"dependent_jobs"`

	// List of ids of jobs that this job is dependent upon.
	ParentJobs []string `json:"parent_jobs"`
}

// Run the job
func (j *Job) Run() {
	j.running.Lock()
	defer j.running.Unlock()

	// Maybe we are testing or it's disabled
	if j.Agent != nil && j.Disabled == false {
		log.WithFields(logrus.Fields{
			"job":      j.Name,
			"schedule": j.Schedule,
		}).Debug("scheduler: Run job")

		ex := &Execution{
			JobName: j.Name,
			Group:   time.Now().UnixNano(),
			Job:     j,
			Attempt: 1,
		}

		cronInspect.Set(j.Name, j)
		j.Agent.RunQuery(ex)
	}
}

// Friendly format a job
func (j *Job) String() string {
	return fmt.Sprintf("\"Job: %s, scheduled at: %s, tags:%v\"", j.Name, j.Schedule, j.Tags)
}
