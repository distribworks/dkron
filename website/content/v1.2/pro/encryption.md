---
title: Encryption
---

SSL encryption is used for communicating dkron pro and the embedded store, and between storage nodes itself. Also if client auth is enabled, only dkron pro clients can talk to the embedded store. This means that no other software running on your local network will be able to talk to dkron's etcd server.

This ensures that no unexpected usage of the Dkron's store will happen, unless it is another Dkron pro instance.

SSL encryption is enabled by default in Dkron Pro and can not be disabled, you don't need to do nothing to use it.

By default Dkron Pro runs with automatically generated SSL certificates, this is enough for using it in a trusted environment but to have a better grade of confidence, it is recommended to run Dkron Pro with custom SSL certificates.

Follow [this tutorial](https://coreos.com/os/docs/latest/generate-self-signed-certificates.html) to generate autosigned SSL certificates for your instances.

{{% notice note %}}
You don't need a client certificate for Dkron server, just add "client auth" usage to your server cert.
{{% /notice %}}

```yaml
# dkron.yaml
auto-tls: false # Set to false to use custom certs
key-file: server-key.pem
cert-file: server.pem
trusted-ca-file: ca.pem
client-cert-auth: true # Enable it to only allow certs signed by the same CA
```

