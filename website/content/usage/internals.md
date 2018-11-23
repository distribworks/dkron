---
title: Internals
---

This document is a WIP, it's intended to describe the reasons that lead to design decisions in Dkron.

## Execution results

Dkron store the result of each job execution in each node.

Every time dkron executes a job it assigns it an execution group, generating a new uuid and send a serf query to target machines and waits for a response.

Each target machine that will run the job, then responds with an execution object saying it started to run the job.

This allows dkron to know how many machines will be running the job.

The design takes into account the differences of how the different storage backends work.

Due to this issue https://github.com/docker/libkv/issues/20 executions are grouped using the group id in the execution object.

## Executions commands output

When a node has finished executing a job it gathers the output of the executed command and sends it back to the a server using an RPC call. This is designed after two main reasons:

1. Scallability, in case of thousands of nodes responding to job the responses are sent to dkron servers in an evenly way selecting a random Dkron server of the ones that are available at the moment and send the response. In the future, Dkron should retry sending the command result with an exponential backoff.

2. Due to the limitations of Serf the queries payload can't be bigger that 1KB, this renders impossible to send a minimal command output togheter with the execution metadata.
