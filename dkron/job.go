package dkron

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/abronan/valkeyrie/store"
)

const (
	Success = iota
	Running
	Failed
	PartialyFailed

	ConcurrencyAllow  = "allow"
	ConcurrencyForbid = "forbid"
)

var (
	ErrParentJobNotFound = errors.New("Specified parent job not found")
	ErrNoAgent           = errors.New("No agent defined")
	ErrSameParent        = errors.New("The job can not have itself as parent")
	ErrNoParent          = errors.New("The job doens't have a parent job set")
	ErrNoCommand         = errors.New("Unespecified command for job")
	ErrWrongConcurrency  = errors.New("Wrong concurrency policy value, use: allow/forbid")
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

	// Extra environment variable to give to the command to execute.
	EnvironmentVariables []string `json:"environment_variables"`

	// Owner of the job.
	Owner string `json:"owner"`

	// Owner email of the job.
	OwnerEmail string `json:"owner_email"`

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

	// Job id of job that this job is dependent upon.
	ParentJob string `json:"parent_job"`

	lock store.Locker

	// Processors to use for this job
	Processors map[string]PluginConfig `json:"processors"`

	// Concurrency policy for this job (allow, forbid)
	Concurrency string `json:"concurrency"`

	// Executor plugin to be used in this job
	Executor string `json:"executor"`

	// Executor args
	ExecutorConfig ExecutorPluginConfig `json:"executor_config"`
}

// Run the job
func (j *Job) Run() {
	j.running.Lock()
	defer j.running.Unlock()

	// Maybe we are testing or it's disabled
	if j.Agent != nil && j.Disabled == false {
		// Check if it's runnable
		if j.isRunnable() {
			log.WithFields(logrus.Fields{
				"job":      j.Name,
				"schedule": j.Schedule,
			}).Debug("scheduler: Run job")

			cronInspect.Set(j.Name, j)

			// Simple execution wrapper
			ex := NewExecution(j.Name)
			j.Agent.RunQuery(ex)
		}
	}
}

// Friendly format a job
func (j *Job) String() string {
	return fmt.Sprintf("\"Job: %s, scheduled at: %s, tags:%v\"", j.Name, j.Schedule, j.Tags)
}

// Return the status of a job
// Wherever it's running, succeded or failed
func (j *Job) Status() int {
	// Maybe we are testing
	if j.Agent == nil {
		return -1
	}

	execs, _ := j.Agent.Store.GetLastExecutionGroup(j.Name)
	success := 0
	failed := 0
	for _, ex := range execs {
		if ex.FinishedAt.IsZero() {
			return Running
		}
	}

	var status int
	for _, ex := range execs {
		if ex.Success {
			success = success + 1
		} else {
			failed = failed + 1
		}
	}

	if failed == 0 {
		status = Success
	} else if failed > 0 && success == 0 {
		status = Failed
	} else if failed > 0 && success > 0 {
		status = PartialyFailed
	}

	return status
}

// Get the parent job of a job
func (j *Job) GetParent() (*Job, error) {
	// Maybe we are testing
	if j.Agent == nil {
		return nil, ErrNoAgent
	}

	if j.Name == j.ParentJob {
		return nil, ErrSameParent
	}

	if j.ParentJob == "" {
		return nil, ErrNoParent
	}

	parentJob, err := j.Agent.Store.GetJob(j.ParentJob)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return nil, ErrParentJobNotFound
		} else {
			return nil, err
		}
	}

	return parentJob, nil
}

// Lock the job in store
func (j *Job) Lock() error {
	// Maybe we are testing
	if j.Agent == nil {
		return ErrNoAgent
	}

	lockKey := fmt.Sprintf("%s/job_locks/%s", j.Agent.Store.keyspace, j.Name)
	// TODO: LockOptions empty is a temporary fix until https://github.com/docker/libkv/pull/99 is fixed
	l, err := j.Agent.Store.Client.NewLock(lockKey, &store.LockOptions{RenewLock: make(chan (struct{}))})
	if err != nil {
		return err
	}
	j.lock = l

	_, err = j.lock.Lock(nil)
	if err != nil {
		return err
	}

	return nil
}

// Unlock the job in store
func (j *Job) Unlock() error {
	// Maybe we are testing
	if j.Agent == nil {
		return ErrNoAgent
	}

	if err := j.lock.Unlock(); err != nil {
		return err
	}

	return nil
}

func (j *Job) isRunnable() bool {
	status := j.Status()

	if status == Running {
		if j.Concurrency == ConcurrencyAllow {
			return true
		} else if j.Concurrency == ConcurrencyForbid {
			log.WithFields(logrus.Fields{
				"job":         j.Name,
				"concurrency": j.Concurrency,
				"job_status":  status,
			}).Debug("scheduler: Skipping execution")
			return false
		}
	}

	return true
}
