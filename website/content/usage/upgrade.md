---
title: Rolling upgrade
---

Use the following procedure to rotate all cluster nodes, one server at a time:

1. Add a new servers to the cluster with a configuration that joins them to the existing cluter.
1. Stop dkron service in one of the old servers, if it was the leader allow a new leader to be ellected, note that it is better to remove the current leader at the end, to ensure a leader is elected between the new nodes.
1. Use `dkron raft list-peers` to list current cluster nodes
1. Use `dkron raft remove-peer` to forcefully remove the old server
1. Repeat steps util all old cluster nodes has been rotated
