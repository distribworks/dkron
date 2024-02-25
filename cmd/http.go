package cmd

import (
	dkplugin "github.com/distribworks/dkron/v4/plugin"
	"github.com/distribworks/dkron/v4/plugin/http"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Hidden: true,
	Use:    "http",
	Short:  "Run the http plugin",
	Long:   ``,
	Run: func(cmd *cobra.Command, args []string) {
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: dkplugin.Handshake,
			Plugins: map[string]plugin.Plugin{
				"executor": &dkplugin.ExecutorPlugin{Executor: http.New()},
			},

			// A non-nil value here enables gRPC serving for this plugin...
			GRPCServer: plugin.DefaultGRPCServer,
		})
	},
}

func init() {
	dkronCmd.AddCommand(httpCmd)
}
