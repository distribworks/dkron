#Getting started

Welcome to the intro guide to dkron! This will explain how to setup dkron, how easy is to use it, what problems could it help you to solve, etc.

## Introduction

Dkron nodes can work in two modes, agents or servers.

Servers are agents too. You can use servers to run jobs.

The main distinction is that servers order job executions and can be used to schedule jobs.

Dkron servers have a leader, the leader is responsible of executing jobs in the cluster.

Any Dkron agent or server acts as a cluster member and it's available to run scheduled tasks.

You can choose whether a job is run on a node or nodes by specifying tags and a count of target nodes having this tag do you want a job to run. For example you can specify to run a job at 5:00am in all servers with `role=web` tag or you can specify to run a job in just one server having the `role=web` tag:

```
role=web:1
```

dkron will try to run the job in the amount of nodes indicated by that count having that tag.

This gives an unprecedented level of flexibility in runnig jobs across a cluster of any size and with any combination of machines you need.

All the execution responses will be gathered by the scheduler and stored in the database.

## Requirements

Dkron relies on the key-value store for data storage, you can run an instance of the distributed store in the same machines as Dkron or connect it to your existing cluster.

It can use etcd, Consul or Zookeeper as data stores. To install any of this systems got to their web site:

- etcd: https://coreos.com/etcd/docs/latest/
- Consul: https://consul.io/intro/getting-started/install.html
- ZooKeeper: https://zookeeper.apache.org/doc/r3.3.3/zookeeperStarted.html

## Installation

Simply download the packaged archive for your platform from the [downloads page](https://github.com/victorcoder/dkron/releases), extract the package to a shared location in your drive, like `/opt/local` and run it from there.

There's a `.deb` package available too.

### Ubuntu

The recommended way to install Dkron is using the `.deb` package.

Sample upstart scripts for dkron are included in the `debian` folder

### Debian

Sample init scripts are included in the `debian` folder

## Configuration

See the [configuration section](configuration).

## Usage

By default dkron will try to use a local etcd server running in the same machine and in the default port. You can specify the store by setting the `backend` and `backend-machines` flag in the config file, env variables or as a command line flag.

To start a dkron server instance just run:

```
dkron agent -server
```
