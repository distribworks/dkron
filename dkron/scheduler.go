package dkron

import (
	"errors"
	"expvar"

	"github.com/armon/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/victorcoder/dkron/cron"
)

var (
	cronInspect      = expvar.NewMap("cron_entries")
	schedulerStarted = expvar.NewInt("scheduler_started")

	ErrScheduleParse = errors.New("Can't parse job schedule")
)

type Cron interface {
	Start()
	Stop()
	Schedule(schedule cron.Schedule, cmd cron.Job)
	Entries() []*cron.Entry
	AddFunc(spec string, cmd func()) error
	AddJob(spec string, cmd cron.Job) error
	AddTimezoneSensitiveJob(spec, timezone string, cmd cron.Job) error
}

type Scheduler struct {
	Cron    Cron
	Started bool
}

func NewScheduler() *Scheduler {
	c := cron.New()
	schedulerStarted.Set(0)
	return &Scheduler{Cron: c, Started: false}
}

func (s *Scheduler) Start(jobs []*Job) {
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

		if job.Timezone != "" {
			s.Cron.AddTimezoneSensitiveJob(job.Schedule, job.Timezone, job)
		} else {
			s.Cron.AddJob(job.Schedule, job)
		}
	}
	s.Cron.Start()
	s.Started = true

	schedulerStarted.Set(1)
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
	}
	schedulerStarted.Set(0)
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
