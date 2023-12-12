package dkron

import (
	"fmt"
	"testing"
	"time"

	"github.com/distribworks/dkron/v3/extcron"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	err := sched.Start([]*Job{testJob1}, &Agent{})
	require.NoError(t, err)

	assert.True(t, sched.started)
	assert.True(t, sched.Started())
	now := time.Now().Truncate(time.Second)

	ej, _ := sched.GetEntryJob(testJob1.Name)
	assert.Equal(t, now.Add(time.Second*2), ej.entry.Next)

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

	sched.Stop()
}

func TestClearCron(t *testing.T) {
	log := getTestLogger()
	sched := NewScheduler(log)

	testJob := &Job{
		Name:           "cron_job",
		Schedule:       "@every 2s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "echo 'test1'", "shell": "true"},
		Owner:          "John Dough",
		OwnerEmail:     "foo@bar.com",
	}
	err := sched.AddJob(testJob)
	require.NoError(t, err)
	assert.Len(t, sched.Cron.Entries(), 1)

	sched.ClearCron()
	assert.Len(t, sched.Cron.Entries(), 0)
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
	_ = sched.Start([]*Job{tzJob}, &Agent{})

	assert.True(t, sched.started)
	assert.True(t, sched.Started())
	assert.Len(t, sched.Cron.Entries(), 1)
	sched.Stop()
}

func TestScheduleStop(t *testing.T) {
	log := getTestLogger()
	sched := NewScheduler(log)

	sched.Cron = cron.New(cron.WithParser(extcron.NewParser()))
	_, err := sched.Cron.AddFunc("@every 2s", func() {
		time.Sleep(time.Second * 5)
		fmt.Println("function done")
	})
	require.NoError(t, err)
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
	err = sched.Start([]*Job{testJob1}, &Agent{})
	assert.Error(t, err)

	// Wait for the job to start
	time.Sleep(time.Second * 2)
	<-sched.Stop().Done()
	err = sched.Start([]*Job{testJob1}, &Agent{})
	assert.NoError(t, err)

	sched.Stop()
	assert.False(t, sched.Started())
}
