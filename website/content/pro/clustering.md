---
title: Clustering
---

## Configure a cluster

Run in HA mode, good for companies.

Configure the peers to join:

```yaml
# dkron.yml
join:
- 10.19.3.9
- 10.19.4.64
- 10.19.7.215
```

```yaml
# etcd.conf.yaml
# Initial cluster configuration for bootstrapping.
initial-cluster: dkron1=https://10.19.3.9:2380,dkron2=https://10.19.4.64:2380,dkron3=https://10.19.7.215:2380
```

With this configuration Dkron Pro should start in cluster mode with embedded storage.
