# Dkron - Distributed, fault tolerant job scheduling system [![GoDoc](https://godoc.org/github.com/victorcoder/dkron?status.svg)](https://godoc.org/github.com/victorcoder/dkron) [![Build Status](https://travis-ci.org/victorcoder/dkron.svg?branch=master)](https://travis-ci.org/victorcoder/dkron)

Website: http://dkron.io/

Dkron is a distributed cron service, easy to setup and fault tolerant with focus in:

- Easy: Easy to use with a great UI
- Reliable: Completly fault tolerant
- High scalable: Able to handle high volumes of scheduled jobs and thousands of nodes

Dkron is written in Go and leverage the power of distributed key-value stores and serf for providing fault tolerance, reliability and scalability while keeping simple and easily instalable.

Dkron is inspired by the google whitepaper [Reliable Cron across the Planet](https://queue.acm.org/detail.cfm?id=2745840)

Dkron runs on Linux, OSX and Windows. It can be used to run scheduled commands on a server cluster using any combination of servers for each job. It has no single points of failure due to the use of the Gossip protocol and fault tolerant distributed databases.

You can use Dkron to run the most important part of your company, scheduled jobs.

## Status

Currently Dkron is under heavy development and should be considered alpha stage.

Said that, I encourage you to try it, it's very easy to use, see how it works for you and report any bugs [creating an issue](https://github.com/victorcoder/dkron/issues) in the github project.

## Quick start

The best way to test and develop dkron is using docker, you will need [Docker Toolbox](https://www.docker.com/docker-toolbox) installed before proceding.

Clone the repository.

Install glide:

`go get github.com/Masterminds/glide`

Download dependencies:

`glide install`

Next, run the included Docker Compose config:

`docker-compose up`

This will start, etc, consul and Dkron instances. To add more Dkron instances to the clusters:

`docker-compose scale dkron=4`

Check the port mapping using `docker-compose ps` and use the browser to navigate to the Dkron dashboard using one of the ports mapped by compose.

To add jobs to the system read the API docs.

## Documentation

Full, comprehensive documentation is viewable on the [Dkron website](http://dkron.io)

## Get in touch

- Twitter: [@dkronio](https://twitter.com/dkronio) or [@victorcoder](https://twitter.com/victorcoder)
- Email: victor at victorcastell.com

# Sponsor

This project is possible thanks to the Support of Jobandtalent

![](https://upload.wikimedia.org/wikipedia/en/d/db/Jobandtalent_logo.jpg)
