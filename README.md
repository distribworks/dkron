<p align="center">
<img width="400" src="docs/images/DKRON_STICKER_OK_CMYK_RGB_CONV_300.png" alt="Dkron" title="Dkron" />
</p>

# Dkron - Distributed, fault tolerant job scheduling system for cloud native environments [![GoDoc](https://godoc.org/github.com/distribworks/dkron?status.svg)](https://godoc.org/github.com/distribworks/dkron) [![Actions Status](https://github.com/distribworks/dkron/workflows/Test/badge.svg)](https://github.com/distribworks/dkron/actions) [![Gitter](https://badges.gitter.im/distribworks/dkron.svg)](https://gitter.im/distribworks/dkron) [![Gurubase](https://img.shields.io/badge/Gurubase-Ask%20Dkron%20Guru-006BFF)](https://gurubase.io/g/dkron)

Website: http://dkron.io/

Dkron is a distributed cron service, easy to setup and fault tolerant with focus in:

- Easy: Easy to use with a great UI
- Reliable: Completely fault tolerant
- Highly scalable: Able to handle high volumes of scheduled jobs and thousands of nodes

Dkron is written in Go and leverage the power of the Raft protocol and Serf for providing fault tolerance, reliability
and scalability while keeping simple and easily installable.

Dkron is inspired by the google
whitepaper [Reliable Cron across the Planet](https://queue.acm.org/detail.cfm?id=2745840) and by Airbnb Chronos
borrowing the same features from it.

Dkron runs on Linux, OSX and Windows. It can be used to run scheduled commands on a server cluster using any combination
of servers for each job. It has no single points of failure due to the use of the Gossip protocol and fault tolerant
distributed databases.

You can use Dkron to run the most important part of your company, scheduled jobs.

## Installation

[Installation instructions](https://dkron.io/docs/basics/installation)

Full, comprehensive documentation is accessible on the [Dkron website](http://dkron.io)

## Quickstart

### Deploying Dkron using Docker

The best way to test and develop dkron is using docker, you will need [Docker](https://www.docker.com/) with Docker
compose installed before proceeding.

```bash
docker compose up -d
```

The UI should be available on http://localhost:8080/ui.

### Using Dkron

To add jobs to the system read the [API docs](https://dkron.io/api/).

### Scaling the cluster

To add more Dkron instances to the cluster:

```bash
docker compose up -d --scale dkron-server=4
docker compose up -d --scale dkron-agent=10
```

## Development

To develop Dkron, you can deploy the cluster with local changes applied with the following steps:

1. Clone the repository.

2. Run the `docker compose`:

    ```bash
    docker compose -f docker-compose.dev.yml up
    ```

### Email Testing

For testing email notifications during development, MailHog is included in the development docker-compose setup. MailHog provides a local SMTP server that captures outgoing emails without sending them to real recipients.

Start MailHog with the development environment:

```bash
docker compose -f docker-compose.dev.yml up mailhog
```

Or run it standalone:

```bash
docker run -p 8025:8025 -p 1025:1025 mailhog/mailhog
```

View captured emails at: http://localhost:8025

For more information, see [docs/EMAIL_TESTING.md](docs/EMAIL_TESTING.md).

### Testing CI Locally

To validate that your changes will pass in GitHub Actions before pushing:

```bash
./scripts/test-ci-locally.sh
```

This script:
- Starts MailHog (simulating the CI service container)
- Runs tests with the same configuration as GitHub Actions
- Provides clear pass/fail results
- Allows you to inspect emails in the MailHog UI

See [.github/TESTING.md](.github/TESTING.md) for more information about CI testing.

### Frontend development

Dkron dashboard is built using [React Admin](https://marmelab.com/react-admin/) as a single page application.

To start developing the dashboard enter the `ui` directory and run `npm install` to get the frontend dependencies and
then start the local server with `npm start` it should start a new local web server and open a new browser window
serving de web ui.

Make your changes to the code, then run `make ui` to generate assets files. This is a method of embedding resources in
Go applications.

### Resources

Chef cookbook
https://supermarket.chef.io/cookbooks/dkron

Python Client Library
https://github.com/oldmantaiter/pydkron

Ruby client
https://github.com/jobandtalent/dkron-rb

PHP client
https://github.com/gromo/dkron-php-adapter

Terraform provider
https://github.com/bozerkins/terraform-provider-dkron

Manage and run jobs in Dkron from your django project
https://github.com/surface-security/django-dkron

## Contributors

<a href="https://github.com/distribworks/dkron/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=distribworks/dkron" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

## Get in touch

- Twitter: [@distribworks](https://twitter.com/distribworks)
- Chat: https://gitter.im/distribworks/dkron
- Email: victor at distrib.works

