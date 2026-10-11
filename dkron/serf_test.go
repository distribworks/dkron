package dkron

import (
	"net"
	"testing"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeJoinDoesNotBootstrapOnLocalJoin(t *testing.T) {
	raftInmem := raft.NewInmemStore()
	require.NoError(t, raftInmem.StoreLog(&raft.Log{Index: 1, Term: 1}))

	a := &Agent{
		config: &Config{
			NodeName:        "node1",
			Region:          "global",
			BootstrapExpect: 3,
		},
		raftInmem:    raftInmem,
		peers:        make(map[string][]*ServerParts),
		localPeers:   make(map[raft.ServerAddress]*ServerParts),
		serverLookup: NewServerLookup(),
		logger:       getTestLogger(),
	}

	a.nodeJoin(serf.MemberEvent{
		Type: serf.EventMemberJoin,
		Members: []serf.Member{
			{
				Name:   "node1",
				Addr:   net.ParseIP("127.0.0.1"),
				Status: serf.StatusAlive,
				Tags: map[string]string{
					"role":    "dkron",
					"server":  "true",
					"region":  "global",
					"dc":      "dc1",
					"version": Version,
					"port":    "6868",
					"expect":  "3",
				},
			},
		},
	})

	assert.Equal(t, 3, a.config.BootstrapExpect)
}
