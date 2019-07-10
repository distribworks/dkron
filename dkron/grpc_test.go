package dkron

import (
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGRPCExecutionDone(t *testing.T) {
	viper.Reset()

	aAddr := testutil.GetBindAddr().String()

	c := DefaultConfig()
	c.BindAddr = aAddr
	c.NodeName = "test-grpc"
	c.Server = true
	c.LogLevel = logLevel
	c.BootstrapExpect = 1
	c.DevMode = true

	a := NewAgent(c)
	a.Start()

	for {
		if a.IsLeader() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	testJob := &Job{
		Name:           "test",
		Schedule:       "@every 1m",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/false"},
		Disabled:       true,
	}

	if err := a.Store.SetJob(testJob, true); err != nil {
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

	rc := NewGRPCClient(nil, a)
	rc.ExecutionDone(a.getRPCAddr(), testExecution)
	execs, _ := a.Store.GetExecutions("test")

	assert.Len(t, execs, 1)
	assert.Equal(t, string(testExecution.Output), string(execs[0].Output))

	// Test store execution on a deleted job
	a.Store.DeleteJob(testJob.Name)

	testExecution.FinishedAt = time.Now()
	err := rc.ExecutionDone(a.getRPCAddr(), testExecution)

	assert.Error(t, err, ErrExecutionDoneForDeletedJob)
}
