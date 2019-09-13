package dkron

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/distribworks/dkron/v2/cron"
	"github.com/distribworks/dkron/v2/ntime"
	"github.com/distribworks/dkron/v2/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
)

const (
	// StatusNotSet is the initial job status.
	StatusNotSet = ""
	// StatusSuccess is status of a job whose last run was a success.
	StatusSuccess = "success"
	// StatusRunning is status of a job whose last run has not finished.
	StatusRunning = "running"
	// StatusFailed is status of a job whose last run was not successful on any nodes.
	StatusFailed = "failed"
	// StatusPartialyFailed is status of a job whose last run was successful on only some nodes.
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
	ErrNoParent = errors.New("The job doesn't have a parent job set")
	// ErrNoCommand is returned when attempting to store a job that has no command.
	ErrNoCommand = errors.New("Unspecified command for job")
	// ErrWrongConcurrency is returned when Concurrency is set to a non existing setting.
	ErrWrongConcurrency = errors.New("Wrong concurrency policy value, use: allow/forbid")
)

// Job descibes a scheduled Job.
type Job struct {
	// Job name. Must be unique, acts as the id.
	Name string `json:"name"`

	// Display name of the job. If present, displayed instead of the name
	DisplayName string `json:"displayname"`

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
	LastSuccess ntime.NullableTime `json:"last_success"`

	// Last time this job failed.
	LastError ntime.NullableTime `json:"last_error"`

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

	// running indicates that the Run method is still broadcasting
	running bool

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

// NewJobFromProto create a new Job from a PB Job struct
func NewJobFromProto(in *proto.Job) *Job {
	next, _ := ptypes.Timestamp(in.GetNext())

	job := &Job{
		Name:           in.Name,
		DisplayName:    in.Displayname,
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
		Next:           next,
	}
	if in.GetLastSuccess().GetHasValue() {
		t, _ := ptypes.Timestamp(in.GetLastSuccess().GetTime())
		job.LastSuccess.Set(t)
	}
	if in.GetLastError().GetHasValue() {
		t, _ := ptypes.Timestamp(in.GetLastError().GetTime())
		job.LastError.Set(t)
	}

	procs := make(map[string]PluginConfig)
	for k, v := range in.Processors {
		procs[k] = v.Config
	}
	job.Processors = procs

	return job
}

// ToProto return the corresponding representation of this Job in proto struct
func (j *Job) ToProto() *proto.Job {
	lastSuccess := &proto.Job_NullableTime{
		HasValue: j.LastSuccess.HasValue(),
	}
	if j.LastSuccess.HasValue() {
		lastSuccess.Time, _ = ptypes.TimestampProto(j.LastSuccess.Get())
	}
	lastError := &proto.Job_NullableTime{
		HasValue: j.LastError.HasValue(),
	}
	if j.LastError.HasValue() {
		lastError.Time, _ = ptypes.TimestampProto(j.LastError.Get())
	}
	next, _ := ptypes.TimestampProto(j.Next)

	processors := make(map[string]*proto.PluginConfig)
	for k, v := range j.Processors {
		processors[k] = &proto.PluginConfig{Config: v}
	}
	return &proto.Job{
		Name:           j.Name,
		Displayname:    j.DisplayName,
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
		Processors:     processors,
		Executor:       j.Executor,
		ExecutorConfig: j.ExecutorConfig,
		Status:         j.Status,
		Metadata:       j.Metadata,
		LastSuccess:    lastSuccess,
		LastError:      lastError,
		Next:           next,
	}
}

// Run the job
func (j *Job) Run() {
	// Maybe we are testing or it's disabled
	if j.Agent != nil && j.Disabled == false {
		// Check if it's runnable
		if j.isRunnable() {
			j.running = true
			defer func() { j.running = false }()

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

// GetStatus returns the status of a job whether it's running, succeeded or failed
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

// GetNext returns the job's next schedule from now
func (j *Job) GetNext() (time.Time, error) {
	if j.Schedule != "" {
		s, err := cron.Parse(j.Schedule)
		if err != nil {
			return time.Time{}, err
		}
		return s.Next(time.Now()), nil
	}

	return time.Time{}, nil
}

func (j *Job) isRunnable() bool {
	if j.Agent.GlobalLock {
		log.WithField("job", j.Name).
			Warning("job: Skipping execution because active global lock")
		return false
	}

	if j.Concurrency == ConcurrencyForbid {
		j.Agent.RefreshJobStatus(j.Name)
	}
	j.Status = j.GetStatus()

	if j.Status == StatusRunning && j.Concurrency == ConcurrencyForbid {
		log.WithFields(logrus.Fields{
			"job":         j.Name,
			"concurrency": j.Concurrency,
			"job_status":  j.Status,
		}).Info("job: Skipping concurrent execution")
		return false
	}

	if j.running {
		log.WithField("job", j.Name).
			Warning("job: Skipping execution because last execution still broadcasting, consider increasing schedule interval")
		return false
	}

	return true
}
