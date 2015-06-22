package dcron

import (
	"testing"
	"time"
)

func testJobs() []*Job {
	var jobs []*Job

	job := &Job{
		Name: "cron_job", Schedule: "@every 2s", Command: "date", Owner: "John Dough", OwnerEmail: "foo@bar.com",
	}
	jobs = append(jobs, job)

	return jobs
}

func TestScheduleReshedule(t *testing.T) {
	sched := NewScheduler()

	testJob1 := &Job{
		Name: "cron_job", Schedule: "@every 2s", Command: "echo 'test1'", Owner: "John Dough", OwnerEmail: "foo@bar.com",
	}
	sched.Start([]*Job{testJob1})

	for _, entry := range sched.Cron.Entries() {
		log.Debug(*entry)
	}

	log.Debug(len(sched.Cron.Entries()))

	sched.Cron.Stop()
	time.Sleep(10 * time.Second)

	testJob2 := &Job{
		Name: "cron_job", Schedule: "@every 5s", Command: "echo 'test2'", Owner: "John Dough", OwnerEmail: "foo@bar.com",
	}
	sched.Restart([]*Job{testJob2})

	for _, entry := range sched.Cron.Entries() {
		log.Debug(*entry)
	}

	log.Debug(len(sched.Cron.Entries()))
	time.Sleep(10 * time.Second)
}
