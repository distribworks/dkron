package dcron

import (
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	etcdc "github.com/coreos/go-etcd/etcd"
	"github.com/mitchellh/cli"
)

func TestAgentCommand_implements(t *testing.T) {
	var _ cli.Command = new(AgentCommand)
}

func TestAgentCommandRun(t *testing.T) {
	log.Level = logrus.FatalLevel

	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	args := []string{
		"-bind", "127.0.0.1:8946",
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a.Run(args)
	}()

	time.Sleep(10 * time.Millisecond)

	// Verify it runs "forever"
	select {
	case <-resultCh:
		t.Fatalf("ended too soon, err: %s", ui.ErrorWriter.String())
	case <-time.After(50 * time.Millisecond):
	}

	// Send a shutdown request
	shutdownCh <- struct{}{}

	select {
	case code := <-resultCh:
		if code != 0 {
			t.Fatalf("bad code: %d", code)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatalf("timeout")
	}
}

func TestAgentCommandElectLeader(t *testing.T) {
	log.Level = logrus.FatalLevel

	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	etcd := etcdc.NewClient([]string{})
	_, err := etcd.DeleteDir("dcron")
	if err != nil {
		if eerr, ok := err.(*etcdc.EtcdError); ok {
			if eerr.ErrorCode == etcdc.ErrCodeEtcdNotReachable {
				t.Fatal("etcd server needed to run tests")
			}
		}
	}

	args := []string{
		"-bind", "127.0.0.1:8947",
		"-join", "127.0.0.1:8948",
		"-node", "test1",
		"-server",
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a.Run(args)
	}()

	// Wait for the first agent to start and set itself as leader
	time.Sleep(2 * time.Second)
	test1Key := a.config.Tags["key"]
	t.Logf("test1 key %s", test1Key)

	// Start another agent
	shutdownCh2 := make(chan struct{})
	defer close(shutdownCh2)

	ui2 := new(cli.MockUi)
	a2 := &AgentCommand{
		Ui:         ui2,
		ShutdownCh: shutdownCh2,
	}

	args2 := []string{
		"-bind", "127.0.0.1:8948",
		"-join", "127.0.0.1:8947",
		"-node", "test2",
		"-server",
	}

	resultCh2 := make(chan int)
	go func() {
		resultCh2 <- a2.Run(args2)
	}()

	time.Sleep(2 * time.Second)
	test2Key := a2.config.Tags["key"]
	t.Logf("test2 key %s", test2Key)

	// Send a shutdown request
	shutdownCh <- struct{}{}

	receiver := make(chan *etcdc.Response)
	stop := make(chan bool)
	time.Sleep(2 * time.Second)

	go etcd.Watch("/dcron/leader", 0, false, receiver, stop)

	// Verify it runs "forever"
	select {
	case res := <-receiver:
		if res.Node.Value == test2Key {
			t.Logf("Leader changed: %s", res.Node.Value)
		}
		stop <- true
	case <-time.After(10 * time.Second):
		t.Fatal("No leader swap occurred")
		stop <- true
	}

	shutdownCh2 <- struct{}{}
}

func Test_processFilteredNodes(t *testing.T) {
	log.Level = logrus.DebugLevel

	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	etcd := etcdc.NewClient([]string{})
	_, err := etcd.DeleteDir("dcron")
	if err != nil {
		if eerr, ok := err.(*etcdc.EtcdError); ok {
			if eerr.ErrorCode == etcdc.ErrCodeEtcdNotReachable {
				t.Fatal("etcd server needed to run tests")
			}
		}
	}

	args := []string{
		"-bind", "127.0.0.1:8949",
		"-join", "127.0.0.1:8950",
		"-node", "test1",
		"-server",
		"-tag", "role=test",
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a.Run(args)
	}()

	// Start another agent
	shutdownCh2 := make(chan struct{})
	defer close(shutdownCh2)

	ui2 := new(cli.MockUi)
	a2 := &AgentCommand{
		Ui:         ui2,
		ShutdownCh: shutdownCh2,
	}

	args2 := []string{
		"-bind", "127.0.0.1:8950",
		"-join", "127.0.0.1:8949",
		"-node", "test2",
		"-server",
		"-tag", "role=test",
	}

	resultCh2 := make(chan int)
	go func() {
		resultCh2 <- a2.Run(args2)
	}()

	job := &Job{
		Name: "test_job_1",
		Tags: map[string]string{
			"role": "test:1",
		},
	}

	time.Sleep(100 * time.Millisecond)
	nodes, err := a.processFilteredNodes(job)

	t.Log(a.serf.Members())
	t.Log(nodes)

	// if nodes[0] != "test1" || nodes[1] != "test2" {
	// 	t.Fatal("Not expected returned nodes")
	// }

	// Send a shutdown request
	shutdownCh <- struct{}{}
	shutdownCh2 <- struct{}{}
}
