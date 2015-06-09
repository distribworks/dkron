package dcron

import (
	"net"
	"testing"
	"time"

	"github.com/mitchellh/cli"
)

func TestProcessFilteredNodes(t *testing.T) {

}

func TestCommand_implements(t *testing.T) {
	var _ cli.Command = new(AgentCommand)
}

func TestCommandRun(t *testing.T) {
	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	c := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	args := []string{
		"-bind", net.IPv4(127, 0, 0, 1).String(),
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

// func testAgent(logOutput io.Writer) *AgentCommand {
// 	return testAgentWithConfig(DefaultConfig(), serf.DefaultConfig(), logOutput)
// }
//
// func testAgentWithConfig(agentConfig *Config, serfConfig *serf.Config,
// 	logOutput io.Writer) *Agent {
//
// 	if logOutput == nil {
// 		logOutput = os.Stderr
// 	}
// 	serfConfig.MemberlistConfig.ProbeInterval = 100 * time.Millisecond
// 	serfConfig.MemberlistConfig.BindAddr = testutil.GetBindAddr().String()
// 	serfConfig.NodeName = serfConfig.MemberlistConfig.BindAddr
//
// 	agent, err := Create(agentConfig, serfConfig, logOutput)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return agent
// }
