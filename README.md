# Dkron - Distributed, fault tolerant job scheduling system

Website: http://dkron.io/

Dkron is a distributed cron service, easy to setup and fault tolerant with focus in:

- Easy: Easy to use with a great UI
- Reliable: Completly fault tolerant
- High scalable: Able to handle high volumes of scheduled jobs and thousands of nodes

Dkron is written in Go and leverage the power of etcd and serf for providing fault tolerance and, reliability and scalability while keeping simple and easily instalable.

Dkron is inspired by the google whitepaper [Reliable Cron across the Planet](https://queue.acm.org/detail.cfm?id=2745840)

Dkron runs on Linux, OSX and Windows. It can be used to run scheduled commands on a server cluster using any combination of servers for each job. It has no single points of failure due to the use of the Gossip protocol and the fault tolerant distributed database etcd.

You can use Dkron to run the most important part of your company, scheduled jobs.

## Status

Currently Dkron is under heavy development and should be considered alpha stage.

Said that, I encourage you to try it, it's very easy to use, see how it works for you and report any bugs [creating an issue](https://github.com/victorcoder/dkron/issues) in the github project.

## Quick start

Clone the repository.

*NOTE*: The included etcd binary is compiled for OSX, if you're in another platform, download the apporpriate etcd binary for your platform.

Setup goreman:

`go get github.com/mattn/goreman`

Next, run the included Procfile

`goreman start`

This will start etcd and some Dkron instances that will form a cluster.

Now you can view the web panel at: http://localhost:8081

To add jobs to the system read the API docs or take a look to the `job.json` file.

## Documentation

Full, comprehensive documentation is viewable on the Dkron website: http://dkron.io

## Get in touch

Twitter: [@dkronio](https://twitter.com/dkronio) or [@victorcoder](https://twitter.com/victorcoder)
Email: victor at victorcastell.com
