package dkron

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGRPCExecutionDone(t *testing.T) {
	store := NewStore("etcdv3", []string{etcdAddr}, nil, "dkron", nil)
	viper.Reset()

	// Cleanup everything
	err := store.Client.DeleteTree("dkron")
	if err != nil {
		t.Logf("error cleaning up: %s", err)
	}

	aAddr := testutil.GetBindAddr().String()

	c := DefaultConfig()
	c.BindAddr = aAddr
	c.BackendMachines = []string{etcdAddr}
	c.NodeName = "test1"
	c.Server = true
	c.LogLevel = logLevel
	c.Keyspace = "dkron"
	c.Backend = "etcdv3"
	c.BackendMachines = []string{os.Getenv("DKRON_BACKEND_MACHINE")}

	a := NewAgent(c, nil)
	a.Start()

	time.Sleep(2 * time.Second)

	testJob := &Job{
		Name:           "test",
		Schedule:       "@every 1m",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/false"},
		Disabled:       true,
	}

	if err := store.SetJob(testJob, true); err != nil {
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

	rc := NewGRPCClient(nil)
	rc.CallExecutionDone(a.getRPCAddr(), testExecution)
	execs, _ := store.GetExecutions("test")

	assert.Len(t, execs, 1)
	assert.Equal(t, string(testExecution.Output), string(execs[0].Output))

	// Test store execution on a deleted job
	store.DeleteJob(testJob.Name)

	testExecution.FinishedAt = time.Now()
	err = rc.CallExecutionDone(a.getRPCAddr(), testExecution)

	assert.Error(t, err, ErrExecutionDoneForDeletedJob)
}
