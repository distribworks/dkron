package dkron

import (
	"testing"
	"time"

	s "github.com/abronan/valkeyrie/store"
)

func TestStore(t *testing.T) {
	store := NewStore("etcd", []string{etcdAddr}, nil, "dkron-test", nil)

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

	// Check that we still get an empty job list
	jobs, err := store.GetJobs()
	if err != nil {
		t.Fatalf("error getting jobs: %s", err)
	} else if jobs == nil {
		t.Fatal("jobs empty, expecting empty slice")
	}

	if err := store.SetJob(testJob, nil); err != nil {
		t.Fatalf("error creating job: %s", err)
	}

	jobs, err = store.GetJobs()
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
