package dcron

import (
	"bitbucket.org/victorcoder/dcron/cron"
	"crypto/md5"
	"fmt"
	"log"
)

type Scheduler struct {
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Load() {
	job := `{"name": "cron job", "schedule": "1 * * * * *", "command": "date"}`
	jobId := fmt.Sprintf("%x", md5.Sum([]byte(job)))

	if _, err := etcdClient.Set("/jobs/"+jobId+"/job", job, 0); err != nil {
		log.Fatal(err)
	}

	if _, err := etcdClient.Set("/jobs/"+jobId+"/started_at", "1224553", 0); err != nil {
		log.Fatal(err)
	}

	res, err := etcdClient.Get("/jobs/"+jobId, false, true)
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range res.Node.Nodes {
		fmt.Println(node.Value)
	}

	c := cron.New()

	c.AddFunc("1 * * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.Start()
}
