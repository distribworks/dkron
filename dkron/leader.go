package dkron

import (
	"fmt"
	"net"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
)

const (
	// barrierWriteTimeout is used to give Raft a chance to process a
	// possible loss of leadership event if we are unable to get a barrier
	// while leader.
	barrierWriteTimeout = 2 * time.Minute
)

// monitorLeadership is used to monitor if we acquire or lose our role
// as the leader in the Raft cluster. There is some work the leader is
// expected to do, so we must react to changes
func (a *Agent) monitorLeadership() {
	var weAreLeaderCh chan struct{}
	var leaderLoop sync.WaitGroup
	for {
		log.Info("dkron: monitoring leadership")
		select {
		case isLeader := <-a.leaderCh:
			switch {
			case isLeader:
				if weAreLeaderCh != nil {
					log.Error("dkron: attempted to start the leader loop while running")
					continue
				}

				weAreLeaderCh = make(chan struct{})
				leaderLoop.Add(1)
				go func(ch chan struct{}) {
					defer leaderLoop.Done()
					a.leaderLoop(ch)
				}(weAreLeaderCh)
				log.Info("dkron: cluster leadership acquired")

			default:
				if weAreLeaderCh == nil {
					log.Error("dkron: attempted to stop the leader loop while not running")
					continue
				}

				log.Debug("dkron: shutting down leader loop")
				close(weAreLeaderCh)
				leaderLoop.Wait()
				weAreLeaderCh = nil
				log.Info("dkron: cluster leadership lost")
			}

		case <-a.shutdownCh:
			return
		}
	}
}

// leaderLoop runs as long as we are the leader to run various
// maintenance activities
func (a *Agent) leaderLoop(stopCh chan struct{}) {
	var reconcileCh chan serf.Member
	establishedLeader := false

RECONCILE:
	// Setup a reconciliation timer
	reconcileCh = nil
	interval := time.After(a.config.ReconcileInterval)

	// Apply a raft barrier to ensure our FSM is caught up
	start := time.Now()
	barrier := a.raft.Barrier(barrierWriteTimeout)
	if err := barrier.Error(); err != nil {
		log.WithError(err).Error("dkron: failed to wait for barrier")
		goto WAIT
	}
	metrics.MeasureSince([]string{"dkron", "leader", "barrier"}, start)

	// Check if we need to handle initial leadership actions
	if !establishedLeader {
		if err := a.establishLeadership(stopCh); err != nil {
			log.WithError(err).Error("dkron: failed to establish leadership")

			// Immediately revoke leadership since we didn't successfully
			// establish leadership.
			if err := a.revokeLeadership(); err != nil {
				log.WithError(err).Error("dkron: failed to revoke leadership")
			}

			goto WAIT
		}

		establishedLeader = true
		defer func() {
			if err := a.revokeLeadership(); err != nil {
				log.WithError(err).Error("dkron: failed to revoke leadership")
			}
		}()
	}

	// Reconcile any missing data
	if err := a.reconcile(); err != nil {
		log.WithError(err).Error("dkron: failed to reconcile")
		goto WAIT
	}

	// Initial reconcile worked, now we can process the channel
	// updates
	reconcileCh = a.reconcileCh

	// Poll the stop channel to give it priority so we don't waste time
	// trying to perform the other operations if we have been asked to shut
	// down.
	select {
	case <-stopCh:
		return
	default:
	}

WAIT:
	// Wait until leadership is lost
	for {
		select {
		case <-stopCh:
			return
		case <-a.shutdownCh:
			return
		case <-interval:
			goto RECONCILE
		case member := <-reconcileCh:
			a.reconcileMember(member)
		}
	}
}

// reconcile is used to reconcile the differences between Serf
// membership and what is reflected in our strongly consistent store.
func (a *Agent) reconcile() error {
	defer metrics.MeasureSince([]string{"dkron", "leader", "reconcile"}, time.Now())
	members := a.serf.Members()
	for _, member := range members {
		if err := a.reconcileMember(member); err != nil {
			return err
		}
	}
	return nil
}

// reconcileMember is used to do an async reconcile of a single serf member
func (a *Agent) reconcileMember(member serf.Member) error {
	// Check if this is a member we should handle
	valid, parts := isServer(member)
	if !valid || parts.Region != a.config.Region {
		return nil
	}
	defer metrics.MeasureSince([]string{"dkron", "leader", "reconcileMember"}, time.Now())

	var err error
	switch member.Status {
	case serf.StatusAlive:
		err = a.addRaftPeer(member, parts)
	case serf.StatusLeft:
		err = a.removeRaftPeer(member, parts)
	}
	if err != nil {
		log.WithError(err).WithField("member", member).Error("failed to reconcile member")
		return err
	}
	return nil
}

// establishLeadership is invoked once we become leader and are able
// to invoke an initial barrier. The barrier is used to ensure any
// previously inflight transactions have been committed and that our
// state is up-to-date.
func (a *Agent) establishLeadership(stopCh chan struct{}) error {
	defer metrics.MeasureSince([]string{"dkron", "leader", "establish_leadership"}, time.Now())

	log.Info("agent: Starting scheduler")
	jobs, err := a.Store.GetJobs(nil)
	if err != nil {
		log.Fatal(err)
	}
	a.sched.Start(jobs, a)

	return nil
}

// revokeLeadership is invoked once we step down as leader.
// This is used to cleanup any state that may be specific to a leader.
func (a *Agent) revokeLeadership() error {
	defer metrics.MeasureSince([]string{"dkron", "leader", "revoke_leadership"}, time.Now())
	a.sched.Stop()

	return nil
}

// addRaftPeer is used to add a new Raft peer when a dkron server joins
func (a *Agent) addRaftPeer(m serf.Member, parts *ServerParts) error {
	// Check for possibility of multiple bootstrap nodes
	members := a.serf.Members()
	if parts.Bootstrap {
		for _, member := range members {
			valid, p := isServer(member)
			if valid && member.Name != m.Name && p.Bootstrap {
				log.Errorf("dkron: '%v' and '%v' are both in bootstrap mode. Only one node should be in bootstrap mode, not adding Raft peer.", m.Name, member.Name)
				return nil
			}
		}
	}

	// Processing ourselves could result in trying to remove ourselves to
	// fix up our address, which would make us step down. This is only
	// safe to attempt if there are multiple servers available.
	addr := (&net.TCPAddr{IP: m.Addr, Port: parts.Port}).String()
	configFuture := a.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		log.WithError(err).Error("dkron: failed to get raft configuration")
		return err
	}

	if m.Name == a.config.NodeName {
		if l := len(configFuture.Configuration().Servers); l < 3 {
			log.WithField("peer", m.Name).Debug("dkron: Skipping self join check since the cluster is too small")
			return nil
		}
	}

	// See if it's already in the configuration. It's harmless to re-add it
	// but we want to avoid doing that if possible to prevent useless Raft
	// log entries. If the address is the same but the ID changed, remove the
	// old server before adding the new one.
	for _, server := range configFuture.Configuration().Servers {

		// If the address or ID matches an existing server, see if we need to remove the old one first
		if server.Address == raft.ServerAddress(addr) || server.ID == raft.ServerID(parts.ID) {
			// Exit with no-op if this is being called on an existing server and both the ID and address match
			if server.Address == raft.ServerAddress(addr) && server.ID == raft.ServerID(parts.ID) {
				return nil
			}
			future := a.raft.RemoveServer(server.ID, 0, 0)
			if server.Address == raft.ServerAddress(addr) {
				if err := future.Error(); err != nil {
					return fmt.Errorf("error removing server with duplicate address %q: %s", server.Address, err)
				}
				log.WithField("server", server.Address).Info("dkron: removed server with duplicate address")
			} else {
				if err := future.Error(); err != nil {
					return fmt.Errorf("error removing server with duplicate ID %q: %s", server.ID, err)
				}
				log.WithField("server", server.ID).Info("dkron: removed server with duplicate ID")
			}
		}
	}

	// Attempt to add as a peer
	switch {
	case minRaftProtocol >= 3:
		addFuture := a.raft.AddVoter(raft.ServerID(parts.ID), raft.ServerAddress(addr), 0, 0)
		if err := addFuture.Error(); err != nil {
			log.WithError(err).Error("dkron: failed to add raft peer")
			return err
		}
	}

	return nil
}

// removeRaftPeer is used to remove a Raft peer when a dkron server leaves
// or is reaped
func (a *Agent) removeRaftPeer(m serf.Member, parts *ServerParts) error {
	// See if it's already in the configuration. It's harmless to re-remove it
	// but we want to avoid doing that if possible to prevent useless Raft
	// log entries.
	configFuture := a.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		log.WithError(err).Error("dkron: failed to get raft configuration")
		return err
	}

	// Pick which remove API to use based on how the server was added.
	for _, server := range configFuture.Configuration().Servers {
		// If we understand the new add/remove APIs and the server was added by ID, use the new remove API
		if minRaftProtocol >= 2 && server.ID == raft.ServerID(parts.ID) {
			log.WithField("server", server.ID).Info("dkron: removing server by ID")
			future := a.raft.RemoveServer(raft.ServerID(parts.ID), 0, 0)
			if err := future.Error(); err != nil {
				log.WithError(err).WithField("server", server.ID).Error("dkron: failed to remove raft peer")
				return err
			}
			break
		}
	}

	return nil
}
