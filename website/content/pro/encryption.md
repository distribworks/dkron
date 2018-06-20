---
title: Encryption
---

SSL encryption is used for communicating dkron pro and the embedded store, and between storage nodes itself. Also client auth is enabled, so only dkron pro clients can talk to the embedded store. This means that no other software running on your local network will be able to talk to dkron's etcd server.

This ensures that no unexpected usage of the Dkron's store will happen, unless it is another Dkron pro instance.

SSL encryption is enabled by default in Dkron Pro and can not be disabled, you don't need to do nothing to use it.


