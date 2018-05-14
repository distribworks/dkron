// Command that implements the main executable.
package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/victorcoder/dkron/cmd"
	"github.com/victorcoder/dkron/dkron"
)

func main() {
	c := cli.NewCLI("dkron", dkron.Version)
	c.Args = os.Args[1:]
	c.HelpFunc = cli.BasicHelpFunc("dkron")

	ui := &cli.BasicUi{Writer: os.Stdout}

	c.Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &cmd.AgentCommand{
				Ui: ui,
			}, nil
		},
		"keygen": func() (cli.Command, error) {
			return &cmd.KeygenCommand{
				Ui: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &cmd.VersionCommand{
				Ui: ui,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
	}
	os.Exit(exitStatus)
}
