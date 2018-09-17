package cmd

import (
	"fmt"

	"github.com/hashicorp/serf/serf"
	"github.com/spf13/cobra"
	"github.com/victorcoder/dkron/dkron"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  `Show the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Dkron %s\n", dkron.Version)
		fmt.Printf("Agent Protocol: %d (Understands back to: %d)\n",
			serf.ProtocolVersionMax, serf.ProtocolVersionMin)
	},
}

func init() {
	dkronCmd.AddCommand(versionCmd)
}
