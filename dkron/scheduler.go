package dkron

import (
	"errors"
	"expvar"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/distribworks/dkron/v3/extcron"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var (
	cronInspect      = expvar.NewMap("cron_entries")
	schedulerStarted = expvar.NewInt("scheduler_started")

	// ErrScheduleParse is the error returned when the schdule parsing fails.
	ErrScheduleParse = errors.New("can't parse job schedule")
)

// Scheduler represents a dkron scheduler instance, it stores the cron engine
// and the related parameters.
type Scheduler struct {
	Cron        *cron.Cron
	Started     bool
	EntryJobMap map[string]cron.EntryID
}

// NewScheduler creates a new Scheduler instance
func NewScheduler() *Scheduler {
	schedulerStarted.Set(0)
	return &Scheduler{
		Cron:        nil,
		Started:     false,
		EntryJobMap: make(map[string]cron.EntryID),
	}
}

// Start the cron scheduler, adding its corresponding jobs and
// executing them on time.
func (s *Scheduler) Start(jobs []*Job, agent *Agent) error {
	s.Cron = cron.New(cron.WithParser(extcron.NewParser()))

	metrics.IncrCounter([]string{"scheduler", "start"}, 1)
	for _, job := range jobs {
		job.Agent = agent
		s.AddJob(job)
	}
	s.Cron.Start()
	s.Started = true
	schedulerStarted.Set(1)

	return nil
}

// Stop stops the scheduler effectively not running any job.
func (s *Scheduler) Stop() {
	if s.Started {
		log.Debug("scheduler: Stopping scheduler")
		s.Cron.Stop()
		s.Started = false
		// Keep Cron exists and let the jobs which have been scheduled can continue to finish,
		// even the node's leadership will be revoked.
		// Ignore the running jobs and make s.Cron to nil may cause whole process crashed.
		//s.Cron = nil

		// expvars
		cronInspect.Do(func(kv expvar.KeyValue) {
			kv.Value = nil
		})
	}
	schedulerStarted.Set(0)
}

// Restart the scheduler
func (s *Scheduler) Restart(jobs []*Job, agent *Agent) {
	s.Stop()
	s.ClearCron()
	s.Start(jobs, agent)
}

// Clear cron separately, this can only be called when agent will be stop.
func (s *Scheduler) ClearCron() {
	s.Cron = nil
}

// GetEntry returns a scheduler entry from a snapshot in
// the current time, and whether or not the entry was found.
func (s *Scheduler) GetEntry(jobName string) (cron.Entry, bool) {
	for _, e := range s.Cron.Entries() {
		j, _ := e.Job.(*Job)
		if j.Name == jobName {
			return e, true
		}
	}
	return cron.Entry{}, false
}

// AddJob Adds a job to the cron scheduler
func (s *Scheduler) AddJob(job *Job) error {
	// Check if the job is already set and remove it if exists
	if _, ok := s.EntryJobMap[job.Name]; ok {
		s.RemoveJob(job)
	}

	if job.Disabled || job.ParentJob != "" {
		return nil
	}

	log.WithFields(logrus.Fields{
		"job": job.Name,
	}).Debug("scheduler: Adding job to cron")

	cronInspect.Set(job.Name, job)
	metrics.EmitKey([]string{"scheduler", "job/update", "add", job.Name}, 1)

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

	id, err := s.Cron.AddJob(schedule, job)
	if err != nil {
		return err
	}
	s.EntryJobMap[job.Name] = id

	return nil
}

// RemoveJob removes a job from the cron scheduler
func (s *Scheduler) RemoveJob(job *Job) {
	log.WithFields(logrus.Fields{
		"job": job.Name,
	}).Debug("scheduler: Removing job from cron")
	s.Cron.Remove(s.EntryJobMap[job.Name])
	delete(s.EntryJobMap, job.Name)
}
