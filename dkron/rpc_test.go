package dkron

import (
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/mitchellh/cli"
)

func TestRPCExecutionDone(t *testing.T) {
	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	aAddr := testutil.GetBindAddr().String()
	rpcAddr := testutil.GetBindAddr().String() + ":6868"

	args := []string{
		"-bind", aAddr,
		"-rpc-addr", rpcAddr,
		"-node", "test1",
		"-server",
		"-debug",
		"-keyspace", "dkron-test",
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a.Run(args)
	}()
	time.Sleep(2 * time.Second)

	store := NewStore("etcd", []string{"127.0.0.1:4001"}, nil, "dkron-test")

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
		ServerAddr: rpcAddr,
	}

	rc.callExecutionDone(testExecution)
	execs, _ := store.GetExecutions("test")

	if len(execs) == 0 {
		t.Fatal("executions result is empty")
	}

	if string(execs[0].Output) != string(testExecution.Output) {
		t.Fatalf("error on retrieved excution expected: %s got: %s", testExecution.Output, execs[0].Output)
	}
}
