// Command that implements the main executable.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/cli"
	"github.com/victorcoder/dkron/cmd"
)

const (
	VERSION = "0.9.7"
)

func main() {
	c := cli.NewCLI("dkron", VERSION)
	c.Args = os.Args[1:]
	c.HelpFunc = cli.BasicHelpFunc("dkron")

	ui := &cli.BasicUi{Writer: os.Stdout}

	plugins := &Plugins{}
	if err := plugins.DiscoverPlugins(); err != nil {
		log.Fatal(err)
	}

	log.Println(plugins.Executors)

	// Make sure we clean up any managed plugins at the end of this

	//protoc -I proto/ proto/executor.proto --go_out=plugins=grpc:proto/

	c.Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &cmd.AgentCommand{
				Ui:               ui,
				Version:          VERSION,
				ProcessorPlugins: plugins.Processors,
				ExecutorPlugins:  plugins.Executors,
			}, nil
		},
		"keygen": func() (cli.Command, error) {
			return &cmd.KeygenCommand{
				Ui: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &cmd.VersionCommand{
				Version: VERSION,
				Ui:      ui,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
	}

	plugin.CleanupClients()
	os.Exit(exitStatus)
}
