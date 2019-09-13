package cmd

import (
	"fmt"

	"github.com/distribworks/dkron/v2/dkron"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var raftCmd = &cobra.Command{
	Use:   "raft [command]",
	Short: "Command to perform some raft operations",
	Long:  ``,
}

var raftListCmd = &cobra.Command{
	Use:   "list-peers",
	Short: "Command to list raft peers",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var gc dkron.DkronGRPCClient
		gc = dkron.NewGRPCClient(nil, nil)

		reply, err := gc.RaftGetConfiguration(rpcAddr)
		if err != nil {
			return err
		}

		// Format it as a nice table.
		result := []string{"Node|ID|Address|State|Voter"}
		for _, s := range reply.Servers {
			state := "follower"
			if s.Leader {
				state = "leader"
			}
			result = append(result, fmt.Sprintf("%s|%s|%s|%s|%v",
				s.Node, s.Id, s.Address, state, s.Voter))
		}

		fmt.Println(columnize.SimpleFormat(result))

		return nil
	},
}

var peerID string

var raftRemovePeerCmd = &cobra.Command{
	Use:   "remove-peer",
	Short: "Command to list raft peers",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var gc dkron.DkronGRPCClient
		gc = dkron.NewGRPCClient(nil, nil)

		if err := gc.RaftRemovePeerByID(rpcAddr, peerID); err != nil {
			return err
		}
		fmt.Println("Peer removed")

		return nil
	},
}

func init() {
	raftCmd.PersistentFlags().StringVar(&rpcAddr, "rpc-addr", "127.0.0.1:6868", "gRPC address of the agent")
	raftRemovePeerCmd.Flags().StringVar(&peerID, "peer-id", "", "Remove a Dkron server with the given ID from the Raft configuration.")

	raftCmd.AddCommand(raftListCmd)
	raftCmd.AddCommand(raftRemovePeerCmd)

	dkronCmd.AddCommand(raftCmd)
}
