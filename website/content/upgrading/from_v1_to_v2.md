---
title: Upgrade from v1 to v2
---

Dkron v2 brings lots of changes to the previous version. To successfully upgrade from v1 to v2 you need to take care of certain changes:

## Migrating Jobs

To migrate jobs from v1 to v2, export jobs from the v1 cluster and import them into v2.

A basic script to do that can be found [here](https://gist.github.com/pjz/94f4bd81a0897fd64db44593078e2156)

You can take the opportunity to change your job tags as explained in the following section.

## Changing tags and metadata

Tags behave different in dkron 2. There are several dkron tags that changed names and have become reserved. The reserved tags are:

- dc
- region
- role
- version
- server
- bootstrap
- expect
- rpc_addr
- port

You can not set these tags with the `tags` param. Tags `region` and `dc` can be set using the `region` and `datacenter` params, respectively.

In Dkron 2, the job filtering API now filters on the `metadata` instead of the `tags` field. Jobs have a new param `metadata` that can be used to set any data to classify jobs. These can then be used to filter results returned by the API.

## Selecting storage dir and backup strategies

Dkron 2 implements a new param `data-dir`, which specifies the working data directory. This directory stores all working data, and it should be backed up and handled with special care.
