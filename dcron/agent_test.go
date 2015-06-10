package dcron

import (
	"testing"
	"time"

	etcdc "github.com/coreos/go-etcd/etcd"
	"github.com/mitchellh/cli"
)

func TestProcessFilteredNodes(t *testing.T) {

}

func TestAgentCommand_implements(t *testing.T) {
	var _ cli.Command = new(AgentCommand)
}

func TestAgentCommandRun(t *testing.T) {
	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	c := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	args := []string{
		"-bind", "127.0.0.1:8946",
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- c.Run(args)
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

	time.Sleep(5 * time.Second)

	leader := a.etcd.GetLeader()
	t.Log(leader)

	// Send a shutdown request
	shutdownCh <- struct{}{}
	shutdownCh2 <- struct{}{}
}
