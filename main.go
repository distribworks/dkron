// Command that implements the main executable.
package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/cli"
	"github.com/victorcoder/dkron/dkron"
)

const (
	VERSION = "0.9.2-b3"
)

func main() {
	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	c := cli.NewCLI("dkron", VERSION)
	c.Args = args
	c.HelpFunc = cli.BasicHelpFunc("dkron")

	ui := &cli.BasicUi{Writer: os.Stdout}

	plugins := &Plugins{}
	plugins.DiscoverPlugins()

	// Make sure we clean up any managed plugins at the end of this
	defer plugin.CleanupClients()

	c.Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &dkron.AgentCommand{
				Ui:               ui,
				Version:          VERSION,
				ProcessorPlugins: plugins.Processors,
			}, nil
		},
		"keygen": func() (cli.Command, error) {
			return &dkron.KeygenCommand{
				Ui: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &dkron.VersionCommand{
				Version: VERSION,
				Ui:      ui,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(exitStatus)
}
