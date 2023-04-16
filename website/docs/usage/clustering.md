---
title: Clustering
---

## Configure a cluster

Dkron can run in HA mode, avoiding SPOFs, this mode provides better scalability and better reliability for users that wants a high level of confidence in the cron jobs they need to run.

Manually bootstrapping a Dkron cluster does not rely on additional tooling, but does require operator participation in the cluster formation process. When bootstrapping, Dkron servers and clients must be started and informed with the address of at least one Dkron server.

As you can tell, this creates a chicken-and-egg problem where one server must first be fully bootstrapped and configured before the remaining servers and clients can join the cluster. This requirement can add additional provisioning time as well as ordered dependencies during provisioning.

First, we bootstrap a single Dkron server and capture its IP address. After we have that nodes IP address, we place this address in the configuration.

1. First bootstrap a node with a configuration like this:

```yaml
# dkron.yml
server: true
bootstrap-expect: 1
```

2. Then stop the bootstrapped server and capture the server IP address.

3. To form a cluster, configure server nodes (including the bootstrapped server) with the address of its peers as in the following example:

```yaml
# dkron.yml
server: true
bootstrap-expect: 3
retry-join:
- 10.19.3.9
- 10.19.4.64
- 10.19.7.215
```

## Deployment Table

Below is a table that shows quorum size and failure tolerance for various
cluster sizes. The recommended deployment is either 3 or 5 servers. A single
server deployment is _**highly**_ discouraged as data loss is inevitable in a
failure scenario.

<table>
  <thead>
    <tr>
      <th>Servers</th>
      <th>Quorum Size</th>
      <th>Failure Tolerance</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>1</td>
      <td>1</td>
      <td>0</td>
    </tr>
    <tr>
      <td>2</td>
      <td>2</td>
      <td>0</td>
    </tr>
    <tr class="warning">
      <td>3</td>
      <td>2</td>
      <td>1</td>
    </tr>
    <tr>
      <td>4</td>
      <td>3</td>
      <td>1</td>
    </tr>
    <tr class="warning">
      <td>5</td>
      <td>3</td>
      <td>2</td>
    </tr>
    <tr>
      <td>6</td>
      <td>4</td>
      <td>2</td>
    </tr>
    <tr>
      <td>7</td>
      <td>4</td>
      <td>3</td>
    </tr>
  </tbody>
</table>
