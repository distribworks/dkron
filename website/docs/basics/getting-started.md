---
sidebar_position: 1
---
# Getting started

## Introduction

Dkron is a distributed job scheduling system that runs across multiple nodes. Understanding its architecture will help you deploy it effectively.

### Node Types

Dkron nodes can work in two modes:

- **Agents**: Cluster members that handle job executions, run scripts, and return the output to servers.
- **Servers**: Cluster members that send job execution queries and also handle scheduling, data storage, and leader election.

```mermaid
flowchart TD
    Leader[Server Node\n(Leader)] <--> Follower1[Server Node\n(Follower)]
    Leader <--> Follower2[Server Node\n(Follower)]
    Leader --> Agent1[Agent Node]
    Leader --> Agent2[Agent Node]
    Follower1 --> Agent1
    Follower2 --> Agent2
```

### Leadership and Job Execution

- A single server node is elected as the **leader** in each Dkron cluster
- The leader is responsible for initiating job execution queries across the cluster
- Any Dkron agent or server can run scheduled jobs
- By default, all nodes execute each job, but this is configurable

### Job Targeting

You can control which nodes run a job by:
- Specifying tags on nodes
- Setting a count of target nodes with specific tags in your job definition

This provides flexibility in running jobs across clusters of any size with precise node targeting.

### Execution Flow

1. The leader schedules a job based on its timing configuration
2. Target nodes are selected based on tags and count
3. Selected nodes execute the job
4. Execution responses are collected by the scheduler
5. Results are stored in the database

## State Storage

Dkron deployment is just a single binary. It stores state in an internal BoltDB instance and replicates all changes between server nodes using the Raft protocol. No external storage system is required.

## Installation

See the [installation](installation.md).

## Configuration

See the [configuration](configuration.md).

## Usage

By default Dkron uses the following ports:

- `8946` for serf layer between agents
- `8080` for HTTP for the API and Dashboard
- `6868` for gRPC and raft layer comunication between agents.

:::info
Be sure you have opened this ports (or the ones that you configured) in your firewall or AWS security groups.
:::

### Starting a single node

Works out of the box, good for non HA installations.

- System service: If no changes are done to the default config files, dkron will start as a service in single mode.
- Command line: Running a single node with default config can be done by running:

```
dkron agent --server --bootstrap-expect=1
```

Check your server is working: `curl localhost:8080/v1`

:::info
For a full list of configuration parameters and its description, see the [CLI docs](/docs/cli/dkron_agent)
:::

### Create a Job

:::info
This job will only run in just one `server` node due to the node count in the tag. Refer to the [target node spec](/docs/usage/target-nodes-spec) for details.
:::

```bash
curl localhost:8080/v1/jobs -XPOST -d '{
  "name": "job1",
  "schedule": "@every 10s",
  "timezone": "Europe/Berlin",
  "owner": "Platform Team",
  "owner_email": "platform@example.com",
  "disabled": false,
  "tags": {
    "server": "true:1"
  },
  "metadata": {
    "user": "12345"
  },
  "concurrency": "allow",
  "executor": "shell",
  "executor_config": {
    "command": "date"
  }
}'
```

For full Job params description refer to the Job model in the [API guide](/api)

That's it!

#### To start configuring an HA installation of Dkron follow the [clustering guide](/docs/usage/clustering)
