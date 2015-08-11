# Dcron - Distributed, fault tolerant job scheduling system

Welcome to the Dcron documentation! This is the reference guide on how to use Dcron. If you want a getting started guide refer to the [getting started guide](user-guide/getting-started).

<div class="alert alert-warning" role="alert">Warning: Dcron is Under heavy development, it's considered alpha state, expect dragons.</div>
<div class="alert alert-info" role="alert">Note: Version in this documentation, 0.0.2</div>

## What is Dcron

Dcron it's distributed system to run scheduled jobs against a server or a group of servers of any size. One of the machines is the leader and the others will be followers. If the leader fails or becomes unreachable, any other one will take over and reschedule all jobs to keep the system healthy.

In case the old leader becomes alive again, it'll become a follower.

Dcron is a distributed cron drop-in replacement, easy to setup and fault tolerant with focus in:

- Easy: Easy to use with a great UI
- Reliable: Completely fault tolerant
- High scalable: Able to handle high volumes of scheduled jobs and thousands of nodes

Dcron is written in Go and leverage the power of [etcd](https://coreos.com/etcd/) and [Serf](https://www.serfdom.io/) for providing fault tolerance and, reliability and scalability while keeping simple and easily installable.

Dcron is inspired by the google whitepaper [Reliable Cron across the Planet](https://queue.acm.org/detail.cfm?id=2745840)

Dcron runs on Linux, OSX and Windows. It can be used to run scheduled commands on a server cluster using any combination of servers for each job. It has no single points of failure due to the use of the fault tolerant distributed database etcd and can work large scale thanks to the efficient and lightweight gossip protocol.

Dcron uses the efficient and lightweight [gossip protocol](https://www.serfdom.io/docs/internals/gossip.html) underneath to communicate with nodes. Failure notification and task handling are run efficiently across an entire cluster of any size.

## Dcron design

Dcron is designed to solve one problem well, executing commands in given intervals. Following the unix philosophy of doing one thing and doing it well (like the battle-tested cron) but with the given addition of being designed for the cloud era, removing single points of failure in environments where scheduled jobs are needed to be run in multiple servers.
