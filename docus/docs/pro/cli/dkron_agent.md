---
date: 2020-05-15
title: "dkron agent"
slug: dkron_agent
url: /pro/cli/dkron_agent/
---
## dkron agent

Start a dkron agent

### Synopsis

Start a dkron agent that schedule jobs, listen for executions and run executors.
It also runs a web UI.

```
dkron agent [flags]
```

### Options

```
      --advertise-addr string           Address used to advertise to other nodes in the cluster. By default, the bind address is advertised. The value supports go-sockaddr/template format.
      --advertise-rpc-port int          Use the value of rpc-port by default
      --auto-tls                        Client TLS using generated certificates (default true)
      --bind-addr string                Specifies which address the agent should bind to for network services, including the internal gossip protocol and RPC mechanism. This should be specified in IP format, and can be used to easily bind all network services to the same address. The value supports go-sockaddr/template format. (default "{{ GetPrivateIP }}:8946")
      --bootstrap-expect int            Provides the number of expected servers in the datacenter. Either this value should not be provided or the value must agree with other servers in the cluster. When provided, Dkron waits until the specified number of servers are available and then bootstraps the cluster. This allows an initial leader to be elected automatically. This flag requires server mode.
      --cert-file string                Path to the client server TLS cert file
      --client-cert-auth                Enable client cert authentication
      --client-crl-file string          Path to the client certificate revocation list file
      --data-dir string                 Specifies the directory to use for server-specific data, including the replicated log. By default, this is the top-level data-dir, like [/var/lib/dkron] (default "dkron.data")
      --datacenter string               Specifies the data center of the local agent. All members of a datacenter should share a local LAN connection. (default "dc1")
      --disable-http-tls                Disable TLS for HTTP WebUI/API regardless of TLS configuration
      --dog-statsd-addr string          DataDog Agent address
      --dog-statsd-tags strings         Datadog tags, specified as key:value
      --enable-prometheus               Enable serving prometheus metrics
      --encrypt string                  Key for encrypting network traffic. Must be a base64-encoded 16-byte key
      --federation-mode string          Federation mode between clusters in different regions (default "active")
  -h, --help                            help for agent
      --http-addr string                Address to bind the UI web server to. Only used when server. The value supports go-sockaddr/template format. (default ":8080")
      --join strings                    An initial agent to join with. This flag can be specified multiple times
      --key-file string                 Path to the client server TLS key file
      --log-level string                Log level (debug|info|warn|error|fatal|panic) (default "info")
      --mail-from string                From email address to use
      --mail-host string                Mail server host address to use for notifications
      --mail-password string            Mail server password to use
      --mail-payload string             Notification mail payload
      --mail-port uint16                Mail server port
      --mail-subject-prefix string      Notification mail subject prefix (default "[Dkron]")
      --mail-username string            Mail server username used for authentication
      --node-name string                Name of this node. Must be unique in the cluster (default "pris.local")
      --password string                 authentication password
      --profile string                  Profile is used to control the timing profiles used (default "lan")
      --raft-multiplier int             An integer multiplier used by servers to scale key Raft timing parameters. Omitting this value or setting it to 0 uses default timing described below. Lower values are used to tighten timing and increase sensitivity while higher values relax timings and reduce sensitivity. Tuning this affects the time it takes to detect leader failures and to perform leader elections, at the expense of requiring more network and CPU resources for better performance. By default, Dkron will use a lower-performance timing that's suitable for minimal Dkron servers, currently equivalent to setting this to a value of 5 (this default may be changed in future versions of Dkron, depending if the target minimum server profile changes). Setting this to a value of 1 will configure Raft to its highest-performance mode is recommended for production Dkron servers. The maximum allowed value is 10. (default 1)
      --region string                   Specifies the region the Dkron agent is a member of. A region typically maps to a geographic region, for example us, with potentially multiple zones, which map to datacenters such as us-west and us-east (default "global")
      --retry-interval string           Time to wait between join attempts. (default "30s")
      --retry-join strings              Address of an agent to join at start time with retries enabled. Can be specified multiple times.
      --retry-max int                   Maximum number of join attempts. Defaults to 0, which will retry indefinitely.
      --rpc-port int                    RPC Port used to communicate with clients. Only used when server. The RPC IP Address will be the same as the bind address (default 6868)
      --serf-reconnect-timeout string   This is the amount of time to attempt to reconnect to a failed node before giving up and considering it completely gone. In Kubernetes, you might need this to about 5s, because there is no reason to try reconnects for default 24h value. Also Raft behaves oddly if node is not reaped and returned with same ID, but different IP. Format there: https://golang.org/pkg/time/#ParseDuration (default "24h")
      --server                          This node is running in server mode
      --statsd-addr string              Statsd address
      --tag strings                     Tag can be specified multiple times to attach multiple key/value tag pairs to the given node, specified as key=value
      --trusted-ca-file string          Path to the client server TLS trusted CA cert file
      --username string                 authentication username
      --webhook-headers strings         Headers to use when calling the webhook URL. Can be specified multiple times
      --webhook-payload string          Body of the POST request to send on webhook call
      --webhook-url string              Webhook url to call for notifications
```

### Options inherited from parent commands

```
      --config string   config file (default is /etc/dkron/dkron.yml)
```

### SEE ALSO

* [dkron](/pro/cli/dkron/)	 - Professional distributed job scheduling system

###### Auto generated by spf13/cobra on 15-May-2020
