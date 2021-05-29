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

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()
	aAddr := ip1.String()

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

	err = a.Store.SetJob(testJob, true)
	require.NoError(t, err)

	testChildJob := &Job{
		Name:           "child-test",
		ParentJob:      testJob.Name,
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/true"},
		Disabled:       false,
	}

	err = a.Store.SetJob(testChildJob, true)
	require.NoError(t, err)

	testExecution := &Execution{
		JobName:    "test",
		Group:      time.Now().UnixNano(),
		StartedAt:  time.Now(),
		NodeName:   "testNode",
		FinishedAt: time.Now(),
		Success:    true,
		Output:     "test",
	}

	log := getTestLogger()
	rc := NewGRPCClient(nil, a, log)
	rc.ExecutionDone(a.advertiseRPCAddr(), testExecution)
	execs, err := a.Store.GetExecutions("test", &ExecutionOptions{})
	require.NoError(t, err)

	assert.Len(t, execs, 1)
	assert.Equal(t, string(testExecution.Output), string(execs[0].Output))

	// Test run a dependent job
	execs, err = a.Store.GetExecutions("child-test", &ExecutionOptions{})
	require.NoError(t, err)

	assert.Len(t, execs, 1)

	// Test job with dependents no delete
	_, err = a.Store.DeleteJob(testJob.Name)
	require.Error(t, err)

	// Remove dependents and parent
	_, err = a.Store.DeleteJob(testChildJob.Name)
	require.NoError(t, err)
	_, err = a.Store.DeleteJob(testJob.Name)
	require.NoError(t, err)

	// Test store execution on a deleted job
	testExecution.FinishedAt = time.Now()
	err = rc.ExecutionDone(a.advertiseRPCAddr(), testExecution)

	assert.Error(t, err, ErrExecutionDoneForDeletedJob)

	// Test ephemeral jobs
	testJob.Ephemeral = true

	err = a.Store.SetJob(testJob, true)
	require.NoError(t, err)

	err = rc.ExecutionDone(a.advertiseRPCAddr(), testExecution)
	assert.NoError(t, err)

	j, err := a.Store.GetJob("test", nil)
	assert.Error(t, err)
	assert.Nil(t, j)
}
