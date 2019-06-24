---
title: Upgrade from v1 to v2
---

Dkron v2 brings lots of changes to the previous version. To successfully upgrade from v1 to v2 you need to take care of certain changes:

## Migrating Jobs

To migrate jobs from v1 to v2 they should be exported from the v1 cluster and imported into v2.

A basic script to do that can be found [here](https://gist.github.com/pjz/94f4bd81a0897fd64db44593078e2156)

You can take the chance to change your job tags as explained in the following section.

## Changing tags and metadata

Tags behaves different in dkron 2. There is several dkron tags that changes its name and becomes reserved, the reserved tags are:

- dc
- region
- role
- version
- server
- bootstrap
- expect
- rpc_addr
- port

You can not set this tags with the `tags` param but tags `dc` and `region` could be set by `region` and `datacenter` params.

Dkron 2 change the job filtering API param from `tags` to `metadata`. Jobs have a new param `metadata` you could use to set any data to classificate jobs and then use the API call to filter results.

## Selecting storage dir and backup strategies

Dkron 2 implements a new param `data-dir` where the working data directory can be specified. This directory stores all working data of dkron and it should be backed up and handled with special attention.
