package dkron

import (
	"fmt"
	"testing"
	"time"

	"github.com/distribworks/dkron/v3/extcron"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

func TestSchedule(t *testing.T) {
	log := getTestLogger()
	sched := NewScheduler(log)

	assert.False(t, sched.started)
	assert.False(t, sched.Started())

	testJob1 := &Job{
		Name:           "cron_job",
		Schedule:       "@every 2s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "echo 'test1'", "shell": "true"},
		Owner:          "John Dough",
		OwnerEmail:     "foo@bar.com",
	}
	sched.Start([]*Job{testJob1}, &Agent{})

	assert.True(t, sched.started)
	assert.True(t, sched.Started())
	now := time.Now().Truncate(time.Second)

	entry, _ := sched.GetEntry(testJob1.Name)
	assert.Equal(t, now.Add(time.Second*2), entry.Next)

	testJob2 := &Job{
		Name:           "cron_job",
		Schedule:       "@every 5s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "echo 'test2'", "shell": "true"},
		Owner:          "John Dough",
		OwnerEmail:     "foo@bar.com",
	}
	sched.Restart([]*Job{testJob2}, &Agent{})

	assert.True(t, sched.started)
	assert.True(t, sched.Started())
	assert.Len(t, sched.Cron.Entries(), 1)

	sched.ClearCron()
	assert.Len(t, sched.Cron.Entries(), 0)

	sched.Stop()
}

func TestTimezoneAwareJob(t *testing.T) {
	log := getTestLogger()
	sched := NewScheduler(log)

	tzJob := &Job{
		Name:           "cron_job",
		Timezone:       "Europe/Amsterdam",
		Schedule:       "@every 2s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "echo 'test1'", "shell": "true"},
	}
	sched.Start([]*Job{tzJob}, &Agent{})

	assert.True(t, sched.started)
	assert.True(t, sched.Started())
	assert.Len(t, sched.Cron.Entries(), 1)
	sched.Stop()
}

func TestScheduleStop(t *testing.T) {
	log := getTestLogger()
	sched := NewScheduler(log)

	sched.Cron = cron.New(cron.WithParser(extcron.NewParser()))
	sched.Cron.AddFunc("@every 2s", func() {
		time.Sleep(time.Second * 5)
		fmt.Println("function done")
	})
	sched.Cron.Start()
	sched.started = true

	testJob1 := &Job{
		Name:           "cron_job",
		Schedule:       "@every 2s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "echo 'test1'", "shell": "true"},
		Owner:          "John Dough",
		OwnerEmail:     "foo@bar.com",
	}
	err := sched.Start([]*Job{testJob1}, &Agent{})
	assert.Error(t, err)

	// Wait for the job to start
	time.Sleep(time.Second * 2)
	<-sched.Stop().Done()
	err = sched.Start([]*Job{testJob1}, &Agent{})
	assert.NoError(t, err)

	sched.Stop()
	assert.False(t, sched.Started())
}
