package dcron

import (
	"bitbucket.org/victorcoder/dcron/cron"
	"encoding/json"
	"fmt"
)

type Scheduler struct {
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Load() {
	job := &Job{Name: "cron_job", Schedule: "@every 2s", Command: "date", Owner: "foo@bar.com"}
	job2 := &Job{Name: "cron_job_2", Schedule: "@every 3s", Command: "echo", Owner: "foo@bar.com"}

	if err := etcd.SetJob(job); err != nil {
		log.Fatal(err)
	}

	if err := etcd.SetJob(job2); err != nil {
		log.Fatal(err)
	}

	jobs, err := etcd.GetJobs()
	if err != nil {
		log.Fatal(err)
	}

	// if _, err := etcd.Client.Set("/dcron/jobs/"+job.Name+"/started_at", "1224553", 0); err != nil {
	// 	log.Fatal(err)
	// }

	res, err := etcd.Client.Get("/dcron/jobs", false, true)
	if err != nil {
		log.Fatal(err)
	}

	c := cron.New()

	for _, node := range res.Node.Nodes {
		for _, jobNode := range node.Nodes {
			var newJob Job
			fmt.Println(jobNode.Value)
			err := json.Unmarshal([]byte(jobNode.Value), &newJob)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(newJob)
			c.AddJob(newJob.Schedule, newJob)
		}
	}
	c.Start()
}

type Job struct {
	Name           string `json:"name"`
	Schedule       string `json:"schedule"`
	Command        string `json:"command"`
	Owner          string `json:"owner"`
	RunAsUser      string `json:"run_as_user"`
	SuccessfulRuns int    `json:"successful_runs"`
	FailedRuns     int    `json:"failed_runs"`
}

func (j Job) Run() {
	fmt.Println("Running: " + j.Command)
}
