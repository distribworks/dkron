package cmd

import (
	dkplugin "github.com/distribworks/dkron/v3/plugin"
	"github.com/distribworks/dkron/v3/plugin/shell"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: dkplugin.Handshake,
			Plugins: map[string]plugin.Plugin{
				"executor": &dkplugin.ExecutorPlugin{Executor: &shell.Shell{}},
			},

			// A non-nil value here enables gRPC serving for this plugin...
			GRPCServer: plugin.DefaultGRPCServer,
		})
	},
}

func init() {
	dkronCmd.AddCommand(shellCmd)
}
