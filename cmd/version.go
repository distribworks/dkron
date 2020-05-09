package cmd

import (
	"fmt"

	"github.com/distribworks/dkron/v3/dkron"
	"github.com/hashicorp/serf/serf"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  `Show the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Name: %s\n", dkron.Name)
		fmt.Printf("Version: %s\n", dkron.Version)
		fmt.Printf("Codename: %s\n", dkron.Codename)
		fmt.Printf("Agent Protocol: %d (Understands back to: %d)\n",
			serf.ProtocolVersionMax, serf.ProtocolVersionMin)
	},
}

func init() {
	dkronCmd.AddCommand(versionCmd)
}
