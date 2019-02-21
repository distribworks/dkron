package cmd

import (
	"github.com/spf13/cobra"
	"github.com/victorcoder/dkron/dkron"
)

var rpcAddr string

// versionCmd represents the version command
var leaveCmd = &cobra.Command{
	Use:   "leave",
	Short: "Force an agent to leave the cluster",
	Long: `Stop stops an agent, if the agent is a server and is running for election
	stop running for election, if this server was the leader
	this will force the cluster to elect a new leader and start a new scheduler.
	If this is a server and has the scheduler started stop it, ignoring if this server
	was participating in leader election or not (local storage).
	Then actually leave the cluster.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var gc dkron.DkronGRPCClient
		gc = dkron.NewGRPCClient(nil)

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
