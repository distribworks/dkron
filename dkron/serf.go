package dkron

import (
	"strings"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
)

const (
	// StatusReap is used to update the status of a node if we
	// are handling a EventMemberReap
	StatusReap = serf.MemberStatus(-1)

	// maxPeerRetries limits how many invalidate attempts are made
	maxPeerRetries = 6
)

// nodeJoin is used to handle join events on the serf cluster
func (a *Agent) nodeJoin(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := isServer(m)
		if !ok {
			a.logger.WithField("member", m.Name).Warn("non-server in gossip pool")
			continue
		}
		a.logger.WithField("server", parts.Name).Info("Adding LAN adding server")
		a.serverLookup.AddServer(parts)
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
		a.logger.WithError(err).Error("failed to read last raft index")
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
			a.logger.WithField("member", member).Error("peer has a conflicting expect value. All nodes should expect the same number")
			return
		}
		if p.Bootstrap {
			a.logger.WithField("member", member).Error("peer has bootstrap mode. Expect disabled")
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

	// Query each of the servers and make sure they report no Raft peers.
	for _, server := range servers {
		var peers []string

		// Retry with exponential backoff to get peer status from this server
		for attempt := uint(0); attempt < maxPeerRetries; attempt++ {
			configuration, err := a.GRPCClient.RaftGetConfiguration(server.RPCAddr.String())
			if err != nil {
				nextRetry := (1 << attempt) * time.Second
				a.logger.Error("Failed to confirm peer status for server (will retry).",
					"server", server.Name,
					"retry_interval", nextRetry.String(),
					"error", err,
				)
				time.Sleep(nextRetry)
			} else {
				for _, peer := range configuration.Servers {
					peers = append(peers, peer.Id)
				}
				break
			}
		}

		// Found a node with some Raft peers, stop bootstrap since there's
		// evidence of an existing cluster. We should get folded in by the
		// existing servers if that's the case, so it's cleaner to sit as a
		// candidate with no peers so we don't cause spurious elections.
		// It's OK this is racy, because even with an initial bootstrap
		// as long as one peer runs bootstrap things will work, and if we
		// have multiple peers bootstrap in the same way, that's OK. We
		// just don't want a server added much later to do a live bootstrap
		// and interfere with the cluster. This isn't required for Raft's
		// correctness because no server in the existing cluster will vote
		// for this server, but it makes things much more stable.
		if len(peers) > 0 {
			a.logger.Info("Existing Raft peers reported by server, disabling bootstrap mode", "server", server.Name)
			a.config.BootstrapExpect = 0
			return
		}
	}

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
	a.logger.Info("agent: found expected number of peers, attempting to bootstrap cluster...",
		"peers", strings.Join(addrs, ","))
	future := a.raft.BootstrapCluster(configuration)
	if err := future.Error(); err != nil {
		a.logger.WithError(err).Error("agent: failed to bootstrap cluster")
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
		a.logger.Info("removing server ", parts)

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
		a.serverLookup.RemoveServer(parts)
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

func (a *Agent) lanNodeUpdate(me serf.MemberEvent) {
	for _, m := range me.Members {
		ok, parts := isServer(m)
		if !ok {
			continue
		}
		a.logger.WithField("server", parts.String()).Info("Updating LAN server")

		// Update server lookup
		a.serverLookup.AddServer(parts)
	}
}
