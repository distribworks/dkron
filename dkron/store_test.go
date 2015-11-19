package dkron

import (
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	store := NewStore("consul", []string{"127.0.0.1:8500"}, nil, "dkron-test")

	// Cleanup everything
	err := store.Client.DeleteTree("dkron-test")
	if err != nil {
		t.Logf("error cleaning up: %s", err)
	}

	testJob := &Job{
		Name:     "test",
		Schedule: "@every 2s",
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
