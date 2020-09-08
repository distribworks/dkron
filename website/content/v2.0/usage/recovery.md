---
title: Outage recovery
---

## Outage Recovery

Don't panic! This is a critical first step.

Depending on your deployment configuration, it may take only a single server failure for cluster unavailability. Recovery requires an operator to intervene, but the process is straightforward.

{{% notice note %}}
This guide is for recovery from a Dkron outage due to a majority of server nodes in a datacenter being lost. If you are looking to add or remove servers, see the [clustering](/usage/clustering) guide.
{{% /notice %}}

## Failure of a Single Server Cluster

If you had only a single server and it has failed, simply restart it. A single server configuration requires the -bootstrap-expect=1 flag. If the server cannot be recovered, you need to bring up a new server. See the [clustering](/usage/clustering) guide for more detail.

In the case of an unrecoverable server failure in a single server cluster, data loss is inevitable since data was not replicated to any other servers. This is why a single server deploy is never recommended.

## Failure of a Server in a Multi-Server Cluster

If you think the failed server is recoverable, the easiest option is to bring it back online and have it rejoin the cluster with the same IP address, returning the cluster to a fully healthy state. Similarly, even if you need to rebuild a new Dkron server to replace the failed node, you may wish to do that immediately. Keep in mind that the rebuilt server needs to have the same IP address as the failed server. Again, once this server is online and has rejoined, the cluster will return to a fully healthy state.

Both of these strategies involve a potentially lengthy time to reboot or rebuild a failed server. If this is impractical or if building a new server with the same IP isn't an option, you need to remove the failed server.

Both of these strategies involve a potentially lengthy time to reboot or rebuild a failed server. If this is impractical or if building a new server with the same IP isn't an option, you need to remove the failed server. Usually, you can issue a `dkron leave` command to remove the failed server if it's still a member of the cluster.

If `dkron leave` isn't able to remove the server, you can use the `dkron raft remove-peer` command to remove the stale peer server on the fly with no downtime.

You can use the `dkron raft list-peers` command to inspect the Raft configuration:

```
$ dkron raft list-peers
Node                   ID               Address          State     Voter
dkron-server01.global  10.10.11.5:4647  10.10.11.5:4647  follower  true
dkron-server02.global  10.10.11.6:4647  10.10.11.6:4647  leader    true
dkron-server03.global  10.10.11.7:4647  10.10.11.7:4647  follower  true
```

## Failure of Multiple Servers in a Multi-Server Cluster

In the event that multiple servers are lost, causing a loss of quorum and a complete outage, partial recovery is possible using data on the remaining servers in the cluster. There may be data loss in this situation because multiple servers were lost, so information about what's committed could be incomplete. The recovery process implicitly commits all outstanding Raft log entries, so it's also possible to commit data that was uncommitted before the failure.

See the section below for details of the recovery procedure. You simply include just the remaining servers in the raft/peers.json recovery file. The cluster should be able to elect a leader once the remaining servers are all restarted with an identical raft/peers.json configuration.

Any new servers you introduce later can be fresh with totally clean data directories.

In extreme cases, it should be possible to recover with just a single remaining server by starting that single server with itself as the only peer in the raft/peers.json recovery file.

The raft/peers.json recovery file is final, and a snapshot is taken after it is ingested, so you are guaranteed to start with your recovered configuration. This does implicitly commit all Raft log entries, so should only be used to recover from an outage, but it should allow recovery from any situation where there's some cluster data available.

## Manual Recovery Using peers.json

To begin, stop all remaining servers. You can attempt a graceful leave, but it will not work in most cases. Do not worry if the leave exits with an error. The cluster is in an unhealthy state, so this is expected.

The peers.json file will be deleted after Dkron starts and ingests this file.

Using raft/peers.json for recovery can cause uncommitted Raft log entries to be implicitly committed, so this should only be used after an outage where no other option is available to recover a lost server. Make sure you don't have any automated processes that will put the peers file in place on a periodic basis.

The next step is to go to the `data-dir` of each Dkron server. Inside that directory, there will be a raft/ sub-directory. We need to create a raft/peers.json file. It should look something like:

```json
[
  {
    "id": "node1",
    "address": "10.1.0.1:4647"
  },
  {
    "id": "node2",
    "address": "10.1.0.2:4647"
  },
  {
    "id": "node3",
    "address": "10.1.0.3:4647"
  }
]
```

Simply create entries for all remaining servers. You must confirm that servers you do not include here have indeed failed and will not later rejoin the cluster. Ensure that this file is the same across all remaining server nodes.

At this point, you can restart all the remaining servers. In Dkron 0.5.5 and later you will see them ingest recovery file:

```
Recovery log placeholder
```

It should be noted that any existing member can be used to rejoin the cluster as the gossip protocol will take care of discovering the server nodes.

At this point, the cluster should be in an operable state again. One of the nodes should claim leadership and emit a log like:

`[INFO] Dkron: cluster leadership acquired`

You can use the `dkron raft list-peers` command to inspect the Raft configuration:

```
$ dkron raft list-peers
Node   ID     Address          State     Voter
node1  node1  10.10.11.5:4647  follower  true
node2  node2  10.10.11.6:4647  leader    true
node3  node3  10.10.11.7:4647  follower  true
```

* id (string: <required>) - Specifies the node ID of the server. This is the `name` of the node.

* address (string: <required>) - Specifies the IP and port of the server in ip:port format. The port is the server's gRPC port used for cluster communications, typically `6868`.
