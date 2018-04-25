package dkron

import (
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRPCExecutionDone(t *testing.T) {
	store := NewStore("etcd", []string{etcdAddr}, nil, "dkron", nil)
	viper.Reset()

	// Cleanup everything
	err := store.Client.DeleteTree("dkron")
	if err != nil {
		t.Logf("error cleaning up: %s", err)
	}

	aAddr := testutil.GetBindAddr().String()

	args := []string{
		"-bind-addr", aAddr,
		"-backend-machine", etcdAddr,
		"-node-name", "test1",
		"-server",
		"-keyspace", "dkron",
		"-log-level", logLevel,
	}

	c := NewConfig(args)
	a := NewAgent(c, nil)
	a.Start()

	time.Sleep(2 * time.Second)

	testJob := &Job{
		Name:     "test",
		Schedule: "@every 1m",
		Command:  "/bin/false",
		Disabled: true,
	}

	if err := store.SetJob(testJob, nil); err != nil {
		t.Fatalf("error creating job: %s", err)
	}

	testExecution := &Execution{
		JobName:    "test",
		Group:      time.Now().UnixNano(),
		StartedAt:  time.Now(),
		NodeName:   "testNode",
		FinishedAt: time.Now(),
		Success:    true,
		Output:     []byte("type"),
	}

	rc := &RPCClient{
		ServerAddr: a.getRPCAddr(),
	}
	rc.callExecutionDone(testExecution)
	execs, _ := store.GetExecutions("test")

	assert.Len(t, execs, 1)
	assert.Equal(t, string(testExecution.Output), string(execs[0].Output))

	// Test store execution on a deleted job
	store.DeleteJob(testJob.Name)

	testExecution.FinishedAt = time.Now()
	err = rc.callExecutionDone(testExecution)

	assert.Error(t, err, ErrExecutionDoneForDeletedJob)
}
