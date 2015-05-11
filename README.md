# dcron - Distributed cron

dcron is a distributed cron service, ultraeasy to setup, fault tolerant with focus in:

Easy: Easy to use with a great UI
Reliable: Completly fault tolerant
High scalable: Able to handle high volumes of scheduled jobs

dcron is written in Go and leverage the power of etcd for maintaining configuration and serf for providing fault tolerance while keeping simple and easily instalable.

## Getting started

Cron jobs are configuration, thus it's configuration is made in /etc/cron.d directory, this makes etcd a great choice to store cron jobs configuration.
