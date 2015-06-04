# Dcron - Distributed cron

Website: http://www.dcron.io

Dcron is a distributed cron service, easy to setup and fault tolerant with focus in:

- Easy: Easy to use with a great UI
- Reliable: Completly fault tolerant
- High scalable: Able to handle high volumes of scheduled jobs and thousands of nodes

Dcron is written in Go and leverage the power of etcd and serf for providing fault tolerance and, reliability and scalability while keeping simple and easily instalable.

Dcron is inspired by the google whitepaper [Reliable Cron across the Planet](https://queue.acm.org/detail.cfm?id=2745840)

Dcron runs on Linux, OSX and Windows. It can be used to run scheduled commands on a server cluster using any combination of servers for each job. It has no single points of failure due to the use of the Gossip protocol and the fault tolerant distributed database etcd.

You can use Dcron to run the most important part of your company, scheduled jobs.

## Quick start

First, download a pre-built Dcron binary for your operating system or [compile Dcron yourself](#developing-dcron).

Setup goreman

`go get github.com/mattn/goreman`

Next, run the included Procfile

`goreman start`

This will start etcd and some Dcron instances that will form a cluster.

Now you can view the web panel at: http://localhost:8080

To add jobs to the system read the API docs.

## Documentation

Full, comprehensive documentation is viewable on the Dcron website:

http://www.dcron.io

## Developing Dcron

If you wish to work on Dcron itself, you'll first need Go installed (version 1.2+ is required). Make sure you have Go properly installed, including setting up your GOPATH.

Next, clone this repository into $GOPATH/src/github.com/hashicorp/dcron and then just type `go build *.go`. In a few moments, you'll have a working serf executable:

$ go build *.go
...
$ bin/dcron
...
note: make will also place a copy of the binary in the first part of your $GOPATH

You can run tests by typing make test.

If you make any changes to the code, run make format in order to automatically format the code according to Go standards.
