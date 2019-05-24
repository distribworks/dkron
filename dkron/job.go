package dkron

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"github.com/victorcoder/dkron/cron"
	"github.com/victorcoder/dkron/proto"
)

const (
	StatusNotSet = ""
	// Success is status of a job whose last run was a success.
	StatusSuccess = "success"
	// Running is status of a job whose last run has not finished.
	StatusRunning = "running"
	// Failed is status of a job whose last run was not successful on any nodes.
	StatusFailed = "failed"
	// PartialyFailed is status of a job whose last run was successful on only some nodes.
	StatusPartialyFailed = "partially_failed"

	// ConcurrencyAllow allows a job to execute concurrency.
	ConcurrencyAllow = "allow"
	// ConcurrencyForbid forbids a job from executing concurrency.
	ConcurrencyForbid = "forbid"
)

var (
	// ErrParentJobNotFound is returned when the parent job is not found.
	ErrParentJobNotFound = errors.New("Specified parent job not found")
	// ErrNoAgent is returned when the job's agent is nil.
	ErrNoAgent = errors.New("No agent defined")
	// ErrSameParent is returned when the job's parent is itself.
	ErrSameParent = errors.New("The job can not have itself as parent")
	// ErrNoParent is returned when the job has no parent.
	ErrNoParent = errors.New("The job doens't have a parent job set")
	// ErrNoCommand is returned when attempting to store a job that has no command.
	ErrNoCommand = errors.New("Unespecified command for job")
	// ErrWrongConcurrency is returned when Concurrency is set to a non existing setting.
	ErrWrongConcurrency = errors.New("Wrong concurrency policy value, use: allow/forbid")
)

// Job descibes a scheduled Job.
type Job struct {
	// Job name. Must be unique, acts as the id.
	Name string `json:"name"`

	// The timezone where the cron expression will be evaluated in.
	// Empty means local time.
	Timezone string `json:"timezone"`

	// Cron expression for the job. When to run the job.
	Schedule string `json:"schedule"`

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

	// Job metadata describes the job and allows filtering from the API.
	Metadata map[string]string `json:"metadata"`

	// Pointer to the calling agent.
	Agent *Agent `json:"-"`

	// Number of times to retry a job that failed an execution.
	Retries uint `json:"retries"`

	running sync.Mutex

	// Jobs that are dependent upon this one will be run after this job runs.
	DependentJobs []string `json:"dependent_jobs"`

	// Job id of job that this job is dependent upon.
	ParentJob string `json:"parent_job"`

	// Processors to use for this job
	Processors map[string]PluginConfig `json:"processors"`

	// Concurrency policy for this job (allow, forbid)
	Concurrency string `json:"concurrency"`

	// Executor plugin to be used in this job
	Executor string `json:"executor"`

	// Executor args
	ExecutorConfig ExecutorPluginConfig `json:"executor_config"`

	// Computed job status
	Status string `json:"status"`

	// Computed next execution
	Next time.Time `json:"next"`
}

func NewJobFromProto(in *proto.Job) *Job {
	lastSuccess, _ := ptypes.Timestamp(in.GetLastSuccess())
	lastError, _ := ptypes.Timestamp(in.GetLastError())
	return &Job{
		Name:           in.Name,
		Timezone:       in.Timezone,
		Schedule:       in.Schedule,
		Owner:          in.Owner,
		OwnerEmail:     in.OwnerEmail,
		SuccessCount:   int(in.SuccessCount),
		ErrorCount:     int(in.ErrorCount),
		Disabled:       in.Disabled,
		Tags:           in.Tags,
		Retries:        uint(in.Retries),
		DependentJobs:  in.DependentJobs,
		ParentJob:      in.ParentJob,
		Concurrency:    in.Concurrency,
		Executor:       in.Executor,
		ExecutorConfig: in.ExecutorConfig,
		Status:         in.Status,
		Metadata:       in.Metadata,
		LastSuccess:    lastSuccess,
		LastError:      lastError,
	}
}

// ToProto return the corresponding proto type
func (j *Job) ToProto() *proto.Job {
	lastSuccess, _ := ptypes.TimestampProto(j.LastSuccess)
	lastError, _ := ptypes.TimestampProto(j.LastError)
	return &proto.Job{
		Name:           j.Name,
		Timezone:       j.Timezone,
		Schedule:       j.Schedule,
		Owner:          j.Owner,
		OwnerEmail:     j.OwnerEmail,
		SuccessCount:   int32(j.SuccessCount),
		ErrorCount:     int32(j.ErrorCount),
		Disabled:       j.Disabled,
		Tags:           j.Tags,
		Retries:        uint32(j.Retries),
		DependentJobs:  j.DependentJobs,
		ParentJob:      j.ParentJob,
		Concurrency:    j.Concurrency,
		Executor:       j.Executor,
		ExecutorConfig: j.ExecutorConfig,
		Status:         j.Status,
		Metadata:       j.Metadata,
		LastSuccess:    lastSuccess,
		LastError:      lastError,
	}
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

			// This should be only called on the actual scheduler execution
			n, err := j.GetNext()
			if err != nil {
				log.WithError(err).WithField("job", j.Name).
					Fatal("agent: Error computing next execution")
			}
			j.Next = n

			j.Agent.RunQuery(j, ex)
		}
	}
}

// Friendly format a job
func (j *Job) String() string {
	return fmt.Sprintf("\"Job: %s, scheduled at: %s, tags:%v\"", j.Name, j.Schedule, j.Tags)
}

// Status returns the status of a job whether it's running, succeded or failed
func (j *Job) GetStatus() string {
	// Maybe we are testing
	if j.Agent == nil {
		return StatusNotSet
	}

	execs, _ := j.Agent.Store.GetLastExecutionGroup(j.Name)
	success := 0
	failed := 0
	for _, ex := range execs {
		if ex.FinishedAt.IsZero() {
			return StatusRunning
		}
	}

	var status string
	for _, ex := range execs {
		if ex.Success {
			success = success + 1
		} else {
			failed = failed + 1
		}
	}

	if failed == 0 {
		status = StatusSuccess
	} else if failed > 0 && success == 0 {
		status = StatusFailed
	} else if failed > 0 && success > 0 {
		status = StatusPartialyFailed
	}

	return status
}

// GetParent returns the parent job of a job
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

	parentJob, err := j.Agent.Store.GetJob(j.ParentJob, nil)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrParentJobNotFound
		}
		return nil, err

	}

	return parentJob, nil
}

// GetNext returns the job's next schedule
func (j *Job) GetNext() (time.Time, error) {
	if j.Schedule != "" {
		s, err := cron.Parse(j.Schedule)
		if err != nil {
			return time.Time{}, err
		}
		return s.Next(j.getLast()), nil
	}

	return time.Time{}, nil
}

func (j *Job) getLast() time.Time {
	if j.LastSuccess.After(j.LastError) {
		return j.LastSuccess
	}
	return j.LastError
}

func (j *Job) isRunnable() bool {
	if j.Concurrency == ConcurrencyForbid {
		j.Agent.RefreshJobStatus(j.Name)
	}
	j.Status = j.GetStatus()

	if j.Status == StatusRunning && j.Concurrency == ConcurrencyForbid {
		log.WithFields(logrus.Fields{
			"job":         j.Name,
			"concurrency": j.Concurrency,
			"job_status":  j.Status,
		}).Debug("scheduler: Skipping execution")
		return false
	}

	return true
}
