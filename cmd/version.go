package cmd

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/serf/serf"
	"github.com/mitchellh/cli"
	"github.com/victorcoder/dkron/dkron"
)

// VersionCommand is a Command implementation prints the version.
type VersionCommand struct {
	Ui cli.Ui
}

func (c *VersionCommand) Help() string {
	return ""
}

func (c *VersionCommand) Run(_ []string) int {
	var versionString bytes.Buffer
	fmt.Fprintf(&versionString, "Dkron v%s", dkron.Version)

	c.Ui.Output(versionString.String())
	c.Ui.Output(fmt.Sprintf("Agent Protocol: %d (Understands back to: %d)",
		serf.ProtocolVersionMax, serf.ProtocolVersionMin))
	return 0
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the Dkron version"
}
