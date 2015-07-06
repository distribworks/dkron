# Dcron - Distributed, fault tolerant job scheduling system

Welcome to the Dcron documentation! This is the reference guide on how to use Drcon. If you want a getting started guide refer to the [getting started guide](getting-started/) of the Dcron documentation.

<div class="alert alert-warning" role="alert">Warning: Under active development. Used at X in production</div>
<div class="alert alert-info" role="alert">Note: Version in this documentation, 0.1</div>

## What is Dcron

Dcron is a system service that runs scheduled tasks at given intervals or times, just like the cron unix service. It differs from it in the sense that it's distributed in several machines in a cluster and if one of that machines (the leader) fails, any other one can take this responsability and keep executing the sheduled tasks without human intervention.

Dcron is a distributed cron service, easy to setup and fault tolerant with focus in:

- Easy: Easy to use with a great UI
- Reliable: Completly fault tolerant
- High scalable: Able to handle high volumes of scheduled jobs and thousands of nodes

Dcron is written in Go and leverage the power of [etcd](https://coreos.com/etcd/) and [Serf](https://www.serfdom.io/) for providing fault tolerance and, reliability and scalability while keeping simple and easily instalable.

Dcron is inspired by the google whitepaper [Reliable Cron across the Planet](https://queue.acm.org/detail.cfm?id=2745840)

Dcron runs on Linux, OSX and Windows. It can be used to run scheduled commands on a server cluster using any combination of servers for each job. It has no single points of failure due to the use of the fault tolerant distributed database etcd and can work large scale thanks to the efficient and lightweight gossip protocol.

Dcron uses the efficient and lightweight [gossip protocol](https://www.serfdom.io/docs/internals/gossip.html) underneath to communicate with nodes. Failure notification and task handling are run efficiently across an entire cluster of any size.

## Dcron design

Dcron is designed to do one task well, executing commands in given intervals, following the unix philosophy of doing one thing and doing it well, like the classic and battle tested cron unix service, with the given addition of being designed for the cloud era, removing single points of failure and clusters of any size are needed to execute scheduled tasks in a decentralized fashion.
