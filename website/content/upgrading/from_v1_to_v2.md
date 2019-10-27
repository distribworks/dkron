---
title: Upgrade from v1 to v2
---

Dkron v2 brings lots of changes to the previous version. To successfully upgrade from v1 to v2 you need to take care of certain steps.

## Storage

Dkron v2 brings native storage, no external storage engine is needed. This change can be shocking for existing users and should be carefully planned to consider Dkron as a stateful service in contrast of a stateless service as in v1.

You must plan a backup strategy for the data directories, you can configure it using the `--data-dir` parameter, by default `./dkron.data`.

This will ensure that you can recover cluster data in case of an unexpected failure.

## Migrating Jobs

To migrate jobs from v1 to v2, export jobs from the v1 cluster and import them into v2.

A basic script to do that can be found [here](https://gist.github.com/pjz/94f4bd81a0897fd64db44593078e2156)

You can take the opportunity to update your job definitions as explained in the following section:

## Change job name

In v2 dkron doesn't accept jobs with invalid naming, adapt your scripts or api calls if necessary to set the job name with valid characters.

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
