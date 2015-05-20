package dcron

import (
	"bitbucket.org/victorcoder/dcron/cron"
	"fmt"
	"time"
)

var sched = NewScheduler()

type Scheduler struct {
	Cron *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Load() {
	jobs, err := etcd.GetJobs()
	if err != nil {
		log.Fatal(err)
	}

	s.Cron = cron.New()

	for _, job := range jobs {
		log.Debugf("Adding job to cron: %v", job)
		s.Cron.AddJob(job.Schedule, job)
	}
	s.Cron.Start()
}

func (s *Scheduler) Reload() {
	s.Cron.Stop()
	s.Load()
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

	Executions []*Execution `json:"-"`
}

func (j Job) Run() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error in job execution", err)
			serf.Terminate()
		}
	}()
	fmt.Println("Running: " + j.Command)
}

type Execution struct {
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	ExitStatus int       `json:"exit_status"`
	Job        *Job      `json:"-"`
}
