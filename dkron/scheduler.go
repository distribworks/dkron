package dkron

import (
	"errors"
	"expvar"

	"github.com/Sirupsen/logrus"
	"github.com/victorcoder/dkron/cron"
)

var (
	cronInspect      = expvar.NewMap("cron_entries")
	schedulerStarted = expvar.NewString("scheduler_started")

	ErrScheduleParse = errors.New("Can't parse job schedule")
)

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
		cronInspect.Set(job.Name, job)

		if job.Disabled || job.ParentJob != "" {
			continue
		}

		log.WithFields(logrus.Fields{
			"job": job.Name,
		}).Debug("scheduler: Adding job to cron")

		s.Cron.AddJob(job.Schedule, job)
	}
	s.Cron.Start()
	s.Started = true

	schedulerStarted.Set("true")
}

func (s *Scheduler) Stop() {
	if s.Started {
		log.Debug("scheduler: Stopping scheduler")
		s.Cron.Stop()
		s.Started = false
		s.Cron = cron.New()

		// expvars
		cronInspect.Do(func(kv expvar.KeyValue) {
			kv.Value = nil
		})
		schedulerStarted.Set("false")
	}
}

func (s *Scheduler) Restart(jobs []*Job) {
	s.Stop()
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
