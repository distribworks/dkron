---
title: Using SSL
---

### Using SSL

By default Dkron Pro runs with automatically generated SSL certificates, this is enough if using it in a trusted environment but to have a better grade of confidence, it is recommended to run Dkron Pro with custom SSL certificates.

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
