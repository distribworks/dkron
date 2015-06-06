package main

import (
	"bitbucket.org/victorcoder/dcron/dcron"
	"github.com/mitchellh/cli"
	"log"
	"os"
)

const (
	VERSION = "0.0.2"
)

func main() {
	c := cli.NewCLI("dcron", VERSION)
	c.Args = os.Args[1:]
	c.HelpFunc = cli.BasicHelpFunc("dcron")

	ui := &cli.BasicUi{Writer: os.Stdout}
	c.Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &dcron.AgentCommand{
				Ui: ui,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
