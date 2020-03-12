## This is a Dkron HA solution for running in Kubernetes

### Overview
#### Bootstrap process
1. StatefulSet launches requested amount of dkron-server pods at the same time, so that cluster could be automatically bootstrapped
2. Dkron itself is launched with wrapper init script
    * Script does the following job
      1. Checks when all dkron-server pods fqdns are resolving to POD ip
      1. Checks if the cluster is already bootstrapped, by checking if there's already a leader.  
      If there's a leader, then init script is not passing --bootstrap-expect parameter to Dkron arguments, to not cause a failover on a working cluster on dkron server pod restart.     
      Otherwise --bootstrap-expect parameter is passed with a value of INITIAL_CLUSTER_SUZE environment variable which should match statefulset replicas values during initial launch.
      1. launches dkron binary with a FQDN list of cluster peers
    4. Now all pods are able to discover each other and select a leader
    5. Waits for SIGTERM

#### Pod restart process
1.  Init script receives SIGTERM on pod termination
1. Init script sends SIGTERM to dkron-server, causing it to shutdown
1. Sleeps for 75 seconds, as it's the calculated by trial and error time to dkron server to forget about the ex-node. It's required, because IP address of a pod changes after a restart, but Raft expects node to come back with same IP
1. Init script sends "raft remove-peer" request to a dkron to remove a node from Raft. That's where Raft forgets about a node that will never comeback with same IP.
1. Exits container
1. Pod terminated
1. New pod started
1. Init script checks when all dkron-server pods fqdns are resolving to POD ip
1. Init script checks if the cluster is already bootstrapped, by checking if there's already a leader. If there's a leader, then init script is not passing --bootstrap-expect parameter to Dkron arguments, as it causes failover on a pod restart. Otherwise --bootstrap-expect parameter is passed with a value of INITIAL_CLUSTER_SUZE environment variable which should match statefulset replicas values during initial launch.
1. Init script launches dkron server with a FQDN list of cluster peers
1. Dkron server joins to cluster with a new IP address
1. Init script waits for SIGTERM

#### dkron-server-leader service
This is a service pointing to a cluster leader.
That is required for checking if there's a cluster is already bootstrapped and removing pod from Raft on shutdown.

#### Label updater
There is a labelupdater sidecar container that is checking self dkron if it's leader or follower, and updating self pod label according to current role.

### Configuration
#### Environment variables
* INITIAL_CLUSTER_SIZE  
  Set the number same as replicas values in StatefulSet.  
  This is required for proper bootstrapping

* STATEFULSET_NAME  
  Set this same as your statefulset name.  
  This is required for proper FQDN build


#### Service account and worker agent discovery
Service account with provided binding to role is used to allow dkron-server pod to update leader label.  
Second thing that Dkron worker agents are using Kubernetes service discovery by labels in order to find the cluster to join.  
All Dkron worker agents should be run with appropriate service account that has permission to list pods by label.  

