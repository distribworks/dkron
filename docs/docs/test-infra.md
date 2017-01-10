# Testing/QA Environment

Dkron is tested continuously to ensure it doesn't break with new changes. Unit and integration tests are run in TravisCI while there is another environment where Dkron is tested for QA.

Several types and combinations of jobs are run continuously to ensure everything works as expected with latests releases.

The testing environment is composed by several single tenant bare metal machines kindly provided by [Packet](https://www.packet.net/).

This environment is public and rebuilt from time to time. Following some useful links:

etcd backend: [](http://test.dkron.io:8080/dashboard)
consul backend: [](http://test.dkron.io:8081/dashboard)
zookeeper backend: [](http://test.dkron.io:8082/dashboard)

Metrics: https://p.datadoghq.com/sb/cccdeb5eb-be392a1952
