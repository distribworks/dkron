package dkron

import (
	"fmt"
	"sync"

	"github.com/hashicorp/raft"
)

// ServerLookup encapsulates looking up servers by id and address
type ServerLookup struct {
	lock            sync.RWMutex
	addressToServer map[raft.ServerAddress]*ServerParts
	idToServer      map[raft.ServerID]*ServerParts
}

func NewServerLookup() *ServerLookup {
	return &ServerLookup{
		lock:            sync.RWMutex{},
		addressToServer: make(map[raft.ServerAddress]*ServerParts),
		idToServer:      make(map[raft.ServerID]*ServerParts),
	}
}

func (sl *ServerLookup) AddServer(server *ServerParts) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	sl.addressToServer[raft.ServerAddress(server.RPCAddr.String())] = server
	sl.idToServer[raft.ServerID(server.ID)] = server
}

func (sl *ServerLookup) RemoveServer(server *ServerParts) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	delete(sl.addressToServer, raft.ServerAddress(server.RPCAddr.String()))
	delete(sl.idToServer, raft.ServerID(server.ID))
}

// Implements the ServerAddressProvider interface
func (sl *ServerLookup) ServerAddr(id raft.ServerID) (raft.ServerAddress, error) {
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	svr, ok := sl.idToServer[id]
	if !ok {
		return "", fmt.Errorf("Could not find address for server id %v", id)
	}
	return raft.ServerAddress(svr.RPCAddr.String()), nil
}

// Server looks up the server by address, returns a boolean if not found
func (sl *ServerLookup) Server(addr raft.ServerAddress) *ServerParts {
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	return sl.addressToServer[addr]
}

func (sl *ServerLookup) Servers() []*ServerParts {
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	var ret []*ServerParts
	for _, svr := range sl.addressToServer {
		ret = append(ret, svr)
	}
	return ret
}

func (sl *ServerLookup) CheckServers(fn func(srv *ServerParts) bool) {
	sl.lock.RLock()
	defer sl.lock.RUnlock()

	for _, srv := range sl.addressToServer {
		if !fn(srv) {
			return
		}
	}
}
