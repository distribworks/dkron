package dcron

import (
	"sync"
	"time"

	"github.com/victorcoder/dcron/cron"
)

type Scheduler struct {
	Cron *cron.Cron
}

func NewScheduler() *Scheduler {
	c := cron.New()
	c.Start()
	return &Scheduler{Cron: c}
}

func (s *Scheduler) Start(jobs []*Job) {
	for _, job := range jobs {
		log.Debugf("Adding job to cron: %v", job)
		s.Cron.AddJob(job.Schedule, job)
	}
	s.Cron.Start()
}

func (s *Scheduler) Restart(jobs []*Job) {
	s.Cron.Stop()
	s.Cron.Stop()
	// entries := s.Cron.Entries()
	// entries = entries[:0]
	s.Cron = cron.New()
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

	running sync.Mutex
}

func (j Job) Run() {
	j.running.Lock()
	defer j.running.Unlock()

	log.Debugf("Running: %s %s", j.Name, j.Schedule)

	// Maybe we are testing
	if j.Agent != nil {
		j.Agent.RunQuery(&j)
	}
}

type Execution struct {
	JobName    string    `json:"job_name"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Success    bool      `json:"success"`
	Output     []byte    `json:"output"`
	Job        *Job      `json:"-"`
}
