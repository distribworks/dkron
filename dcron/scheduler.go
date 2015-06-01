package dcron

import (
	"time"

	"bitbucket.org/victorcoder/dcron/cron"
)

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() *Scheduler {
	c := cron.New()
	c.Start()
	return &Scheduler{cron: c}
}

func (s *Scheduler) Start(jobs []*Job) {
	for _, job := range jobs {
		log.Debugf("Adding job to cron: %v", job)
		s.cron.AddJob(job.Schedule, job)
	}
	s.cron.Start()
}

func (s *Scheduler) Restart(jobs []*Job) {
	s.cron.Stop()
	s.cron = cron.New()
	s.Start(jobs)
}

type Job struct {
	Name         string            `json:"name"`
	Schedule     string            `json:"schedule"`
	Command      string            `json:"command"`
	Owner        string            `json:"owner"`
	OwnerEmail   string            `json:"owner_email"`
	RunAsUser    string            `json:"run_as_user"`
	SuccessCount int               `json:"success_count"`
	ErrorCount   int               `json:"error_count"`
	LastSuccess  time.Time         `json:"last_success"`
	LastError    time.Time         `json:"last_error"`
	Disabled     bool              `json:"disabled"`
	Tags         map[string]string `json:"tags"`

	Executions []*Execution  `json:"-"`
	Agent      *AgentCommand `json:"-"`
}

func (j Job) Run() {
	log.Debug("Running: " + j.Name)
	j.Agent.RunQuery(&j)
}

type Execution struct {
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	ExitStatus int       `json:"exit_status"`
	Job        *Job      `json:"-"`
}
