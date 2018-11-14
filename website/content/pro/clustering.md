---
title: Clustering
---

## Configure a cluster

Dkron can run in HA mode, avoiding SPOFs, this mode provides better scalability and better reliability for users that wants a high level of confidence in the cron jobs they need to run.

To form a cluster, server nodes need to know the address of its peers as in the following example:

```yaml
# dkron.yml
join:
- 10.19.3.9
- 10.19.4.64
- 10.19.7.215
```

On the other side, the embedded store also needs to know its peers, it needs its own configuration as in the following example:

```yaml
# etcd.conf.yaml
# Initial cluster configuration for bootstrapping.
initial-cluster: dkron1=https://10.19.3.9:2380,dkron2=https://10.19.4.64:2380,dkron3=https://10.19.7.215:2380
```

With this configuration Dkron Pro should start in cluster mode with embedded storage.
