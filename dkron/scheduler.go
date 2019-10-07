package dkron

import (
	"errors"
	"expvar"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/distribworks/dkron/v2/cron"
	"github.com/distribworks/dkron/v2/extcron"
	"github.com/sirupsen/logrus"
)

var (
	cronInspect      = expvar.NewMap("cron_entries")
	schedulerStarted = expvar.NewInt("scheduler_started")

	// ErrScheduleParse is the error returned when the schdule parsing fails.
	ErrScheduleParse = errors.New("Can't parse job schedule")
)

// Scheduler represents a dkron scheduler instance, it stores the cron engine
// and the related parameters.
type Scheduler struct {
	Cron    *cron.Cron
	Started bool
}

// NewScheduler creates a new Scheduler instance
func NewScheduler() *Scheduler {
	schedulerStarted.Set(0)
	return &Scheduler{Cron: nil, Started: false}
}

// Start the cron scheduler, adding its corresponding jobs and
// executing them on time.
func (s *Scheduler) Start(jobs []*Job) {
	s.Cron = cron.New(cron.WithParser(extcron.NewParser()))
	metrics.IncrCounter([]string{"scheduler", "start"}, 1)
	for _, job := range jobs {
		if job.Disabled || job.ParentJob != "" {
			continue
		}

		log.WithFields(logrus.Fields{
			"job": job.Name,
		}).Debug("scheduler: Adding job to cron")

		cronInspect.Set(job.Name, job)
		metrics.EmitKey([]string{"scheduler", "job", "add", job.Name}, 1)

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
		s.Cron.AddJob(schedule, job)
	}
	s.Cron.Start()
	s.Started = true

	schedulerStarted.Set(1)
}

// Stop stops the scheduler effectively not running any job.
func (s *Scheduler) Stop() {
	if s.Started {
		log.Debug("scheduler: Stopping scheduler")
		s.Cron.Stop()
		s.Started = false
		s.Cron = nil

		// expvars
		cronInspect.Do(func(kv expvar.KeyValue) {
			kv.Value = nil
		})
	}
	schedulerStarted.Set(0)
}

// Restart the scheduler
func (s *Scheduler) Restart(jobs []*Job) {
	s.Stop()
	s.Start(jobs)
}

// GetEntry returns a scheduler entry from a snapshot in
// the current time, and whether or not the entry was found.
func (s *Scheduler) GetEntry(job *Job) (cron.Entry, bool) {
	for _, e := range s.Cron.Entries() {
		j, _ := e.Job.(*Job)
		if j.Name == job.Name {
			return e, true
		}
	}
	return cron.Entry{}, false
}
