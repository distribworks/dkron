---
title: "Intro"
weight: 80
icon: "<b>1. </b>"
---

## Dkron - Distributed, fault tolerant job scheduling system

Welcome to the Dkron documentation! This is the reference guide on how to use Dkron. If you want a getting started guide refer to the [getting started guide](getting-started).

## What is Dkron

Dkron is a distributed system to run scheduled jobs against a server or a group of servers of any size. One of the machines is the leader and the others will be followers. If the leader fails or becomes unreachable, any other one will take over and reschedule all jobs to keep the system healthy.

In case the old leader becomes alive again, it'll become a follower.

Dkron is a distributed cron drop-in replacement, easy to setup and fault tolerant with focus in:

- Easy: Easy to use with a great UI
- Reliable: Completely fault tolerant
- Highly scalable: Able to handle high volumes of scheduled jobs and thousands of nodes

Dkron is written in Go and leverages the power of distributed key value stores and [Serf](https://www.serfdom.io/) for providing fault tolerance, reliability and scalability while remaining simple and easily installable.

Dkron is inspired by the google whitepaper [Reliable Cron across the Planet](https://queue.acm.org/detail.cfm?id=2745840)

Dkron runs on Linux, OSX and Windows. It can be used to run scheduled commands on a server cluster using any combination of servers for each job. It has no single points of failure due to the use of the fault tolerant distributed databases and can work at large scale thanks to the efficient and lightweight gossip protocol.

Dkron uses the efficient and lightweight [gossip protocol](https://www.serfdom.io/docs/internals/gossip.html) underneath to communicate with nodes. Failure notification and task handling are run efficiently across an entire cluster of any size.

## Web UI
<center class="hidden-xs">
![](/img/screenshot1.png)
</center>

## Dkron design

Dkron is designed to solve one problem well, executing commands in given intervals. Following the unix philosophy of doing one thing and doing it well (like the battle-tested cron) but with the given addition of being designed for the cloud era, removing single points of failure in environments where scheduled jobs are needed to be run in multiple servers.
