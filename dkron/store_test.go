package dkron

import (
	s "github.com/docker/libkv/store"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	store := NewStore("etcd", []string{etcdAddr}, nil, "dkron-test")

	// Cleanup everything
	err := store.Client.DeleteTree("dkron-test")
	if err != s.ErrKeyNotFound {
		t.Logf("error cleaning up: %s", err)
	}

	testJob := &Job{
		Name:     "test",
		Schedule: "@every 2s",
		Command:  "/bin/false",
		Disabled: true,
	}

	if err := store.SetJob(testJob); err != nil {
		t.Fatalf("error creating job: %s", err)
	}

	jobs, err := store.GetJobs()
	if err != nil {
		t.Fatalf("error getting jobs: %s", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("error in number of expected jobs: %v", jobs)
	}
	if jobs[0].Name != "test" {
		t.Fatalf("expected job name: %s got: %s", testJob.Name, jobs[0].Name)
	}

	if _, err := store.DeleteJob("test"); err != nil {
		t.Fatalf("error deleting job: %s", err)
	}

	if _, err := store.DeleteJob("test"); err == nil {
		t.Fatalf("error job deletion should fail: %s", err)
	}

	testExecution := &Execution{
		JobName:    "test",
		StartedAt:  time.Now(),
		FinishedAt: time.Now(),
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode",
	}

	_, err = store.SetExecution(testExecution)
	if err != nil {
		t.Fatalf("error setting the execution: %s", err)
	}

	execs, err := store.GetExecutions("test")
	if err != nil {
		t.Fatalf("error getting executions: %s", err)
	}

	if len(execs) == 0 {
		t.Fatal("executions result is empty")
	}

	if !execs[0].StartedAt.Equal(testExecution.StartedAt) {
		t.Fatalf("error on retrieved excution expected: %s got: %s", testExecution.StartedAt, execs[0].StartedAt)
	}

	if len(execs) != 1 {
		t.Fatalf("error in number of expected executions: %v", execs)
	}
}

func TestEmptyStoreShouldReturnEmptyJobsList(t *testing.T) {
	store := NewStore("etcd", []string{etcdAddr}, nil, "dkron-test")
	jobs, err := store.GetJobs()
	assert.Nil(t, err, "Getting empty jobs should not return any errors")
	assert.NotNil(t, jobs, "Getting empty jobs should not return a nil value")
	assert.Empty(t, jobs, "Jobs should be an empty list")
}
