package cmd

import (
	"github.com/distribworks/dkron/v4/dkron"
	"github.com/hashicorp/serf/serf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  `Show the version`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Name: %s\n", dkron.Name)
		log.Infof("Version: %s\n", dkron.Version)
		log.Infof("Codename: %s\n", dkron.Codename)
		log.Infof("Agent Protocol: %d (Understands back to: %d)\n",
			serf.ProtocolVersionMax, serf.ProtocolVersionMin)
	},
}

func init() {
	dkronCmd.AddCommand(versionCmd)
}
