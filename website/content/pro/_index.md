---
title: Dkron Pro
---

## Getting started

Dkron Pro provides a backend store out of the box, based on etcd.

To configure the storage a sample `etcd.conf.yaml` file is provided in `/etc/dkron` path. Editing the file, allows to configure several options for the embedded store.

The location of the store configuration can be set in the command line or in the dkron config file `/etc/dkron/dkron.yml` using `etcd-config-file-path` parameter.

### Starting a single node

Works out of the box, good for non HA installations.

- System service: If no changes are done to the default config files, dkron will start as a service in single mode.
- Command line: Running a single node with default config can be done by running: `dkron agent --server`

Check your server is working: `curl localhost:8080/v1`

Simple as that, now it is time to add some jobs:

```bash
curl localhost:8080/v1/jobs -XPOST -d '{
  "name": "job1",
  "schedule": "@every 10s",
  "timezone": "Europe/Berlin",
  "owner": "Platform Team",
  "owner_email": "platform@example.com",
  "disabled": false,
  "tags": {
    "dkron_server": "true"
  },
  "concurrency": "allow",
  "executor": "shell",
  "executor_config": {
    "command": "date"
  }
}`
```

## Configuration

Dkron Pro uses the [same parameters](/basics/configuration) as Dkron OSS and add some extra parameters.

### Command line options

* `--etcd-config-file-path` - Etcd node config
* `--username` - Authentication username
* `--password` - Authentication password
* `--cert-file` - Path to the client server TLS cert file
* `--key-file` - Path to the client server TLS key file
* `--client-crl-file` - Path to the client certificate revocation list file
* `--trusted-ca-file` - Path to the client server TLS trusted CA cert file
* `--client-cert-auth` - Enable client cert authentication
* `--auto-tls` - Client TLS using generated certificates

{{% children style="card" depth="2"  description="true" %}}
