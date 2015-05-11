package dcron

import (
	"bitbucket.org/victorcoder/dcron/cron"
	"fmt"
	"log"
)

type Scheduler struct {
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Load() {
	if _, err := etcdClient.Set("/", "bar", 0); err != nil {
		log.Fatal(err)
	}

	c := cron.New()

	c.AddFunc("1 * * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.Start()
}
