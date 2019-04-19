---
title: Dkron vs. Other Software
wight: 5
toc: false
---

## Dkron vs. Chronos

Airbnb's Chronos is a job scheduler that is similar to dkron, it's distributed and fault tolerant thanks to the use of Zookeeper and Apache Mesos.

If you don't have/want to run a Mesos cluster and deal with the not easy configuration and maintenance of Zookeeper and you want something lighter, Dkron could help you.

## Dkron vs. Rundeck

Rundeck is a popular and mature platform to automate operations and schedule jobs.

It has cool features:

- Agentless
- Permissions and auditing

It's written in Java and it's not trivial to setup right.

It uses a central database to store job execution results and configuration data, that makes it vulnerable to failures, and you need to take care of providing an HA environment for the database yourself, and that's not an easy task to do with Rundeck's supported databases.

Dkron lacks some of its features but it's lightweight and fault-tolerant out-of-the-box.
