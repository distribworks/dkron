package dkron

import (
	"strings"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
)

const (
	// StatusReap is used to update the status of a node if we
	// are handling a EventMemberReap
	StatusReap = serf.MemberStatus(-1)
)

// nodeJoin is used to handle join events on the serf cluster
func (a *Agent) nodeJoin(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := isServer(m)
		if !ok {
			log.WithField("member", m.Name).Warn("non-server in gossip pool")
			continue
		}
		log.WithField("server", parts.Name).Info("adding server")

		// Check if this server is known
		found := false
		a.peerLock.Lock()
		existing := a.peers[parts.Region]
		for idx, e := range existing {
			if e.Name == parts.Name {
				existing[idx] = parts
				found = true
				break
			}
		}

		// Add ot the list if not known
		if !found {
			a.peers[parts.Region] = append(existing, parts)
		}

		// Check if a local peer
		if parts.Region == a.config.Region {
			a.localPeers[raft.ServerAddress(parts.Addr.String())] = parts
		}
		a.peerLock.Unlock()

		// If we still expecting to bootstrap, may need to handle this
		if a.config.BootstrapExpect != 0 {
			a.maybeBootstrap()
		}
	}
}

// maybeBootstrap is used to handle bootstrapping when a new server joins
func (a *Agent) maybeBootstrap() {
	// Bootstrap can only be done if there are no committed logs, remove our
	// expectations of bootstrapping. This is slightly cheaper than the full
	// check that BootstrapCluster will do, so this is a good pre-filter.
	var index uint64
	var err error
	if a.raftStore != nil {
		index, err = a.raftStore.LastIndex()
	} else if a.raftInmem != nil {
		index, err = a.raftInmem.LastIndex()
	} else {
		panic("neither raftInmem or raftStore is initialized")
	}
	if err != nil {
		log.WithError(err).Error("failed to read last raft index")
		return
	}

	// Bootstrap can only be done if there are no committed logs,
	// remove our expectations of bootstrapping
	if index != 0 {
		a.config.BootstrapExpect = 0
		return
	}

	// Scan for all the known servers
	members := a.serf.Members()
	var servers []ServerParts
	voters := 0
	for _, member := range members {
		valid, p := isServer(member)
		if !valid {
			continue
		}
		if p.Region != a.config.Region {
			continue
		}
		if p.Expect != 0 && p.Expect != a.config.BootstrapExpect {
			log.WithField("member", member).Error("peer has a conflicting expect value. All nodes should expect the same number")
			return
		}
		if p.Bootstrap {
			log.WithField("member", member).Error("peer has bootstrap mode. Expect disabled")
			return
		}
		if valid {
			voters++
		}
		servers = append(servers, *p)
	}

	// Skip if we haven't met the minimum expect count
	if voters < a.config.BootstrapExpect {
		return
	}

	// TODO: Query each of the servers and make sure they report no Raft peers.

	// Update the peer set
	// Attempt a live bootstrap!
	var configuration raft.Configuration
	var addrs []string

	for _, server := range servers {
		addr := server.Addr.String()
		addrs = append(addrs, addr)
		id := raft.ServerID(server.ID)
		suffrage := raft.Voter
		peer := raft.Server{
			ID:       id,
			Address:  raft.ServerAddress(addr),
			Suffrage: suffrage,
		}
		configuration.Servers = append(configuration.Servers, peer)
	}
	log.Info("agent: found expected number of peers, attempting to bootstrap cluster...",
		"peers", strings.Join(addrs, ","))
	future := a.raft.BootstrapCluster(configuration)
	if err := future.Error(); err != nil {
		log.WithError(err).Error("agent: failed to bootstrap cluster")
	}

	// Bootstrapping complete, or failed for some reason, don't enter this again
	a.config.BootstrapExpect = 0
}

// nodeFailed is used to handle fail events on the serf cluster
func (a *Agent) nodeFailed(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := isServer(m)
		if !ok {
			continue
		}
		log.Info("removing server", "server", parts)

		// Remove the server if known
		a.peerLock.Lock()
		existing := a.peers[parts.Region]
		n := len(existing)
		for i := 0; i < n; i++ {
			if existing[i].Name == parts.Name {
				existing[i], existing[n-1] = existing[n-1], nil
				existing = existing[:n-1]
				n--
				break
			}
		}

		// Trim the list there are no known servers in a region
		if n == 0 {
			delete(a.peers, parts.Region)
		} else {
			a.peers[parts.Region] = existing
		}

		// Check if local peer
		if parts.Region == a.config.Region {
			delete(a.localPeers, raft.ServerAddress(parts.Addr.String()))
		}
		a.peerLock.Unlock()
	}
}

// localMemberEvent is used to reconcile Serf events with the
// consistent store if we are the current leader.
func (a *Agent) localMemberEvent(me serf.MemberEvent) {
	// Do nothing if we are not the leader
	if !a.config.Server || !a.IsLeader() {
		return
	}

	// Check if this is a reap event
	isReap := me.EventType() == serf.EventMemberReap

	// Queue the members for reconciliation
	for _, m := range me.Members {
		// Change the status if this is a reap event
		if isReap {
			m.Status = StatusReap
		}
		select {
		case a.reconcileCh <- m:
		default:
		}
	}
}
