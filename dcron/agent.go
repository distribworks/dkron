package dcron

import (
	"flag"
	"github.com/mitchellh/cli"
	"strings"
)

// AgentCommand run dcron agent
type AgentCommand struct {
	Ui cli.Ui
}

func (s *AgentCommand) Help() string {
	helpText := `
Usage: dcron agent [options]
	Provides debugging information for operators
Options:
  -format                  If provided, output is returned in the specified
                           format. Valid formats are 'json', and 'text' (default)
`
	return strings.TrimSpace(helpText)
}

func (a *AgentCommand) Run(args []string) int {
	var format string
	cmdFlags := flag.NewFlagSet("agent", flag.ContinueOnError)
	cmdFlags.Usage = func() { a.Ui.Output(a.Help()) }
	cmdFlags.StringVar(&format, "format", "text", "output format")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	InitSerfAgent()
	return 0
}

func (s *AgentCommand) Synopsis() string {
	return "Run dcron agent"
}
