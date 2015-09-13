package dkron

import (
	"testing"
	"time"
)

func TestEtcdClient(t *testing.T) {
	etcd := NewEtcdClient([]string{}, nil, "dkron-test")

	// Cleanup everything
	_, err := etcd.Client.Delete("dkron-test", true)
	if err != nil {
		t.Fatalf("error cleaning up: %s", err)
	}

	testJob := &Job{
		Name:     "test",
		Schedule: "@every 2s",
		Disabled: true,
	}

	if err := etcd.SetJob(testJob); err != nil {
		t.Fatalf("error creating job: %s", err)
	}

	jobs, err := etcd.GetJobs()
	if err != nil {
		t.Fatalf("error getting jobs: %s", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("error in number of expected jobs: %v", jobs)
	}

	if _, err := etcd.DeleteJob("test"); err != nil {
		t.Fatalf("error deleting job: %s", err)
	}

	if _, err := etcd.DeleteJob("test"); err == nil {
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

	_, err = etcd.SetExecution(testExecution)
	if err != nil {
		t.Fatalf("error setting the execution: %s", err)
	}

	execs, err := etcd.GetExecutions("test")
	if err != nil {
		t.Fatalf("error getting executions: %s", err)
	}

	if execs[0].StartedAt != testExecution.StartedAt {
		t.Fatalf("error on retrieved excution expected: %s got: %s", testExecution.StartedAt, execs[0].StartedAt)
	}

	if len(execs) != 1 {
		t.Fatalf("error in number of expected executions: %v", execs)
	}
}
