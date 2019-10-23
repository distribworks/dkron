package cmd

import (
	"github.com/distribworks/dkron/v2/dkron"
	"github.com/spf13/cobra"
)

var rpcAddr string

// versionCmd represents the version command
var leaveCmd = &cobra.Command{
	Use:   "leave",
	Short: "Force an agent to leave the cluster",
	Long: `Stop stops an agent, if the agent is a server and is running for election
	stop running for election, if this server was the leader
	this will force the cluster to elect a new leader and start a new scheduler.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var gc dkron.DkronGRPCClient
		gc = dkron.NewGRPCClient(nil, nil)

		if err := gc.Leave(rpcAddr); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	dkronCmd.AddCommand(leaveCmd)
	leaveCmd.PersistentFlags().StringVar(&rpcAddr, "rpc-addr", "127.0.0.1:6868", "gRPC address of the agent")
}
