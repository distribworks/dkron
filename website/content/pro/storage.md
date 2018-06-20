---
title: Embedded storage
---

Dkron Pro has an embedded distributed KV store engine based on etcd. This works out of the box on each node dkron server is started.

This ensures a dead easy install and setup, basically run dkron and you will have a full working node and at the same time provides you with a fully tested well supported store for its use with dkron.

## Configuration

The embedded etcd instance configuration can be tunned using the standard etcd config yaml file located in `config/etcd.conf.yml` but several reserved parameters are autoconfigured by dkron. Refer to the [official etcd documentation](https://coreos.com/etcd/docs/latest/v2/configuration.html)
