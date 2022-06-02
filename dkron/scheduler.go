package dkron

import (
	"context"
	"errors"
	"expvar"
	"strings"
	"sync"

	"github.com/armon/go-metrics"
	"github.com/distribworks/dkron/v3/extcron"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var (
	cronInspect      = expvar.NewMap("cron_entries")
	schedulerStarted = expvar.NewInt("scheduler_started")

	// ErrScheduleParse is the error returned when the schedule parsing fails.
	ErrScheduleParse = errors.New("can't parse job schedule")
)

type EntryJob struct {
	entry *cron.Entry
	job   *Job
}

// Scheduler represents a dkron scheduler instance, it stores the cron engine
// and the related parameters.
type Scheduler struct {
	// mu is to prevent concurrent edits to Cron and Started
	mu      sync.RWMutex
	Cron    *cron.Cron
	started bool
	logger  *logrus.Entry
}

// NewScheduler creates a new Scheduler instance
func NewScheduler(logger *logrus.Entry) *Scheduler {
	schedulerStarted.Set(0)
	return &Scheduler{
		Cron:    cron.New(cron.WithParser(extcron.NewParser())),
		started: false,
		logger:  logger,
	}
}

// Start the cron scheduler, adding its corresponding jobs and
// executing them on time.
func (s *Scheduler) Start(jobs []*Job, agent *Agent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return errors.New("scheduler: cron already started, should be stopped first")
	}
	s.ClearCron()

	metrics.IncrCounter([]string{"scheduler", "start"}, 1)
	for _, job := range jobs {
		job.Agent = agent
		if err := s.AddJob(job); err != nil {
			return err
		}
	}
	s.Cron.Start()
	s.started = true
	schedulerStarted.Set(1)

	return nil
}

// Stop stops the cron scheduler if it is running; otherwise it does nothing.
// A context is returned so the caller can wait for running jobs to complete.
func (s *Scheduler) Stop() context.Context {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := s.Cron.Stop()
	if s.started {
		s.logger.Debug("scheduler: Stopping scheduler")
		s.started = false

		// expvars
		cronInspect.Do(func(kv expvar.KeyValue) {
			kv.Value = nil
		})
	}
	schedulerStarted.Set(0)
	return ctx
}

// Restart the scheduler
func (s *Scheduler) Restart(jobs []*Job, agent *Agent) {
	// Stop the scheduler, running jobs will continue to finish but we
	// can not actively wait for them blocking the execution here.
	s.Stop()

	if err := s.Start(jobs, agent); err != nil {
		s.logger.Fatal(err)
	}
}

// ClearCron clears the cron scheduler
func (s *Scheduler) ClearCron() {
	for _, e := range s.Cron.Entries() {
		if j, ok := e.Job.(*Job); !ok {
			s.logger.Errorf("scheduler: Failed to cast job to *Job found type %T and removing it", e.Job)
			s.Cron.Remove(e.ID)
		} else {
			s.RemoveJob(j.Name)
		}
	}
}

// Started will safely return if the scheduler is started or not
func (s *Scheduler) Started() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.started
}

// GetEntryJob returns a EntryJob object from a snapshot in
// the current time, and whether or not the entry was found.
func (s *Scheduler) GetEntryJob(jobName string) (EntryJob, bool) {
	for _, e := range s.Cron.Entries() {
		if j, ok := e.Job.(*Job); !ok {
			s.logger.Errorf("scheduler: Failed to cast job to *Job found type %T", e.Job)
		} else {
			j.logger = s.logger
			if j.Name == jobName {
				return EntryJob{
					entry: &e,
					job:   j,
				}, true
			}
		}
	}
	return EntryJob{}, false
}

// AddJob Adds a job to the cron scheduler
func (s *Scheduler) AddJob(job *Job) error {
	// Check if the job is already set and remove it if exists
	if _, ok := s.GetEntryJob(job.Name); ok {
		s.RemoveJob(job.Name)
	}

	if job.Disabled || job.ParentJob != "" {
		return nil
	}

	s.logger.WithFields(logrus.Fields{
		"job": job.Name,
	}).Debug("scheduler: Adding job to cron")

	// If Timezone is set on the job, and not explicitly in its schedule,
	// AND its not a descriptor (that don't support timezones), add the
	// timezone to the schedule so robfig/cron knows about it.
	schedule := job.Schedule
	if job.Timezone != "" &&
		!strings.HasPrefix(schedule, "@") &&
		!strings.HasPrefix(schedule, "TZ=") &&
		!strings.HasPrefix(schedule, "CRON_TZ=") {
		schedule = "CRON_TZ=" + job.Timezone + " " + schedule
	}

	_, err := s.Cron.AddJob(schedule, job)
	if err != nil {
		return err
	}

	cronInspect.Set(job.Name, job)
	metrics.IncrCounterWithLabels([]string{"scheduler", "job_add"}, 1, []metrics.Label{{Name: "job", Value: job.Name}})

	return nil
}

// RemoveJob removes a job from the cron scheduler if it exists.
func (s *Scheduler) RemoveJob(jobName string) {
	s.logger.WithFields(logrus.Fields{
		"job": jobName,
	}).Debug("scheduler: Removing job from cron")

	if ej, ok := s.GetEntryJob(jobName); ok {
		s.Cron.Remove(ej.entry.ID)
		cronInspect.Delete(jobName)
		metrics.IncrCounterWithLabels([]string{"scheduler", "job_delete"}, 1, []metrics.Label{{Name: "job", Value: jobName}})
	}
}
