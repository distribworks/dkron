package dcron

import (
	"testing"
	"time"

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

	args := []string{
		"-bind", "127.0.0.1:8950",
		"-join", "127.0.0.1:8947",
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

	args = []string{
		"-bind", "127.0.0.1:8950",
		"-join", "127.0.0.1:8946",
	}

	resultCh2 := make(chan int)
	go func() {
		resultCh2 <- a2.Run(args)
	}()

	// Send a shutdown request
	shutdownCh <- struct{}{}
	shutdownCh2 <- struct{}{}
}
