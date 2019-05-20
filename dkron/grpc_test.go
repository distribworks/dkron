package dkron

import (
	"os"
	"testing"
	"time"

	"github.com/abronan/valkeyrie/store"
	"github.com/hashicorp/serf/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/victorcoder/dkron/plugintypes"
)

func TestGRPCExecutionDone(t *testing.T) {
	s := NewStore(store.Backend(backend), []string{backendMachine}, nil, "dkron", nil)
	viper.Reset()

	// Cleanup everything
	err := s.Client().DeleteTree("dkron")
	if err != nil {
		t.Logf("error cleaning up: %s", err)
	}

	aAddr := testutil.GetBindAddr().String()

	c := DefaultConfig()
	c.BindAddr = aAddr
	c.BackendMachines = []string{backendMachine}
	c.NodeName = "test1"
	c.Server = true
	c.LogLevel = logLevel
	c.Keyspace = "dkron"
	c.Backend = store.Backend(backend)
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

	if err := s.SetJob(testJob, true); err != nil {
		t.Fatalf("error creating job: %s", err)
	}

	testExecution := &plugintypes.Execution{
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
	execs, _ := s.GetExecutions("test")

	assert.Len(t, execs, 1)
	assert.Equal(t, string(testExecution.Output), string(execs[0].Output))

	// Test store execution on a deleted job
	s.DeleteJob(testJob.Name)

	testExecution.FinishedAt = time.Now()
	err = rc.CallExecutionDone(a.getRPCAddr(), testExecution)

	assert.Error(t, err, ErrExecutionDoneForDeletedJob)
}
