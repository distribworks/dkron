---
title: Configuration
wight: 20
---

## Configuration sources

Settings can be specified in three ways (in order of precedence): 

1. Command line arguments.
1. Environment variables starting with **`DKRON_`**
1. **`dkron.json`** config file

### Config file example

```yaml
# Dkron example configuration file
# backend: etcd
# backend-machine: 127.0.0.1:2379
# server: false
# log-level: debug
# tags:
#   role: web
#   datacenter: east
# keyspace: dkron
# encrypt: a-valid-key-generated-with-dkron-keygen
# join:
#   - 10.0.0.1
#   - 10.0.0.2
#   - 10.0.0.3
# webhook-url: https://hooks.slack.com/services/XXXXXX/XXXXXXX/XXXXXXXXXXXXXXXXXXXX
# webhook-payload: "payload={\"text\": \"{{.Report}}\", \"channel\": \"#foo\"}"
# webhook-headers: Content-Type:application/x-www-form-urlencoded
# mail-host: email-smtp.eu-west-1.amazonaws.com
# mail-port: 25
# mail-username": mailuser
# mail-password": mailpassword
# mail-from": cron@example.com
# mail-subject_prefix: [Dkron]
```

### SEE ALSO

* [dkron agent](/cli/dkron_agent/)	 - Start a dkron agent
* [dkron doc](/cli/dkron_doc/)	 - Generate Markdown documentation for the Dkron CLI.
* [dkron keygen](/cli/dkron_keygen/)	 - Generates a new encryption key
* [dkron version](/cli/dkron_version/)	 - Show version
