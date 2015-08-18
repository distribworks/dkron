package dkron

import (
	"testing"

	"github.com/Sirupsen/logrus"
)

func testJobs() []*Job {
	var jobs []*Job

	job := &Job{
		Name: "cron_job", Schedule: "@every 2s", Command: "date", Owner: "John Dough", OwnerEmail: "foo@bar.com",
	}
	jobs = append(jobs, job)

	return jobs
}

func TestSchedule(t *testing.T) {
	log.Level = logrus.FatalLevel

	sched := NewScheduler()

	if sched.Started == true {
		t.Fatal("The scheduler should be stopped.")
	}

	testJob1 := &Job{
		Name: "cron_job", Schedule: "@every 2s", Command: "echo 'test1'", Owner: "John Dough", OwnerEmail: "foo@bar.com",
	}
	sched.Start([]*Job{testJob1})

	if sched.Started != true {
		t.Fatal("The scheduler should be started.")
	}

	testJob2 := &Job{
		Name: "cron_job", Schedule: "@every 5s", Command: "echo 'test2'", Owner: "John Dough", OwnerEmail: "foo@bar.com",
	}
	sched.Restart([]*Job{testJob2})

	if sched.Started != true {
		t.Fatal("The scheduler should be started.")
	}

	if len(sched.Cron.Entries()) > 1 {
		t.Fatal("The scheduler has more jobs than expected.")
	}
}
