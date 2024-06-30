package dkron

import (
	"testing"

	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/require"
)

type testAddr struct {
	addr string
}

func (ta *testAddr) Network() string {
	return "tcp"
}

func (ta *testAddr) String() string {
	return ta.addr
}

func TestAddServer(t *testing.T) {
	// arrange
	lookup := NewServerLookup()
	id1, addr1 := "server-1", "127.0.0.1:8300"
	id2, addr2 := "server-2", "127.0.0.2:8300"
	server1, server2 := buildServerParts(id1, addr1), buildServerParts(id2, addr2)

	// act
	lookup.AddServer(server1)
	lookup.AddServer(server2)

	// assert
	servers := lookup.Servers()
	require.Containsf(t, servers, server1, "Expected %v to contain %+v", servers, server1)
	require.Containsf(t, servers, server2, "Expected %v to contain %+v", servers, server2)

	got, err := lookup.ServerAddr(raft.ServerID(id1))
	require.NoErrorf(t, err, "Unexpected error: %v", err)
	require.EqualValuesf(t, addr1, string(got), "Expected %v but got %v", addr1, got)

	server := lookup.Server(raft.ServerAddress(addr1))
	strAddr := server.RPCAddr.String()
	require.EqualValuesf(t, addr1, strAddr, "Expected lookup to return address %v but got %v", addr1, strAddr)

	got, err = lookup.ServerAddr(raft.ServerID(id2))
	require.NoErrorf(t, err, "Unexpected error: %v", err)
	require.EqualValuesf(t, addr2, string(got), "Expected %v but got %v", addr2, got)

	server = lookup.Server(raft.ServerAddress(addr2))
	strAddr = server.RPCAddr.String()
	require.EqualValuesf(t, addr2, strAddr, "Expected lookup to return address %v but got %v", addr2, strAddr)
}

func TestRemoveServer(t *testing.T) {
	// arrange
	lookup := NewServerLookup()
	id1, addr1 := "server-1", "127.0.0.1:8300"
	id2, addr2 := "server-2", "127.0.0.2:8300"
	server1, server2 := buildServerParts(id1, addr1), buildServerParts(id2, addr2)
	lookup.AddServer(server1)
	lookup.AddServer(server2)

	// act
	lookup.RemoveServer(server1)

	// assert
	servers := lookup.Servers()
	expectedServers := []*ServerParts{server2}
	require.EqualValuesf(t, expectedServers, servers, "Expected %v but got %v", expectedServers, servers)

	require.Nilf(t, lookup.Server(raft.ServerAddress(addr1)), "Expected lookup to return nil")
	addr, err := lookup.ServerAddr(raft.ServerID(id1))
	require.Errorf(t, err, "Expected lookup to return error")
	require.EqualValuesf(t, "", string(addr), "Expected empty address but got %v", addr)

	got, err := lookup.ServerAddr(raft.ServerID(id2))
	require.NoErrorf(t, err, "Unexpected error: %v", err)
	require.EqualValuesf(t, addr2, string(got), "Expected %v but got %v", addr2, got)

	server := lookup.Server(raft.ServerAddress(addr2))
	strAddr := server.RPCAddr.String()
	require.EqualValuesf(t, addr2, strAddr, "Expected lookup to return address %v but got %v", addr2, strAddr)
}

func buildServerParts(id, addr string) *ServerParts {
	return &ServerParts{
		ID:      id,
		RPCAddr: &testAddr{addr},
	}
}
