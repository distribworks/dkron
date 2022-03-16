---
title: Clustering
---

## Configure a cluster

First follow the Dkron [clustering guide](/usage/clustering) then you can continue with this guide.

The embedded store also needs to know its peers, it needs its own configuration as in the following example:

```yaml
# etcd.conf.yaml
# Initial cluster configuration for bootstrapping.
initial-cluster: dkron1=https://10.19.3.9:2380,dkron2=https://10.19.4.64:2380,dkron3=https://10.19.7.215:2380
```

With this configuration Dkron Pro should start in cluster mode with embedded storage.

For a more in detail guide of clustering with etcd follow this guide: https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/clustering.md
