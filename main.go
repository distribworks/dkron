// Command that implements the main executable.
package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/victorcoder/dkron/dkron"
)

const (
	VERSION = "0.0.4"
)

func main() {
	c := cli.NewCLI("dkron", VERSION)
	c.Args = os.Args[1:]
	c.HelpFunc = cli.BasicHelpFunc("dkron")

	ui := &cli.BasicUi{Writer: os.Stdout}
	c.Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &dkron.AgentCommand{
				Ui:      ui,
				Version: VERSION,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
