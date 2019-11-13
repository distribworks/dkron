package dkron

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGRPCExecutionDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	viper.Reset()

	aAddr := testutil.GetBindAddr().String()

	c := DefaultConfig()
	c.BindAddr = aAddr
	c.NodeName = "test-grpc"
	c.Server = true
	c.LogLevel = logLevel
	c.BootstrapExpect = 1
	c.DevMode = true
	c.DataDir = dir

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
		Schedule:       "@manually",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/true"},
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
	err = rc.ExecutionDone(a.getRPCAddr(), testExecution)

	assert.Error(t, err, ErrExecutionDoneForDeletedJob)
}
