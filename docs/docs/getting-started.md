#Getting started

Welcome to the intro guide to dkron! This will explain how to setup dkron, how easy is to use it, what problems could it help you to solve, etc.

## Introduction

dkron nodes can work in two modes, agents or servers.

Servers are agents too. You can use servers to run jobs.

The main distinction is that servers order job executions and can be used to schedule jobs.

dkron servers have a leader, the leader is responsible of executing jobs in the cluster.

Any dkron agent or server acts as a cluster member and it's available to run scheduled tasks.

You can choose whether a job is run on a node or nodes by specifying tags and a count of target nodes having this tag do you want a job to run. For example you can specify to run a job at 5:00am in all servers with role=web tag or you can specify to run a job in just one server of having the role=web tag:

```
role=web:1
```

dkron will try to run the job in the amount of nodes indicated by that count having that tag.

This gives an unprecedented level of flexibility in runnig jobs across a cluster of any size and with any combination of machines you need.

All the execution responses will be gathered by the scheduler and stored in the database.

## Installation

Simply download the packaged archive for your platform from the downloads page, extract the package to a shared location in your drive, like `/opt/local` and run it from there.

### Ubuntu

Sample upstart scripts for dkron are included in the `extras` folder

### Debian

Sample init scripts are included in the `extras` folder

## Configuration

See the [configuration section](configuration).

## Usage

dkron relies on etcd for data storage, the etcd executable is included in the package and can be used to run an etcd node along with dkron servers.

By default dkron will try to use a local etcd server running in the same machine and in the default port.

If you want to use an existing etcd cluster of your own, you can specify it by setting the `etcd-machies` flag in the config file or as a command line flag.

To start a dkron server instance just run:

```
dkron agent -server
```
