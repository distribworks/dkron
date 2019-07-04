
---
title: Installation
weight: 10
---

## Running the binary

Download the packaged archive for your platform from the [downloads page](https://github.com/distribworks/dkron/releases) and extract the package to a shared location in your drive, like /opt/local/bin.

Run Dkron with default setting: `dkron agent --server`

Navigate to http://localhost:8080


{{% notice note %}}
By default dkron will start with a file based, embedded KV store called BoltDB, it is functional for a single node demo but does not offers clustering or HA.
{{% /notice %}}

## Installing the package

### Debian repo

APT repository: 
```
deb [trusted=yes] https://apt.fury.io/distribworks/ /
```

Then install: `sudo apt-get install dkron`

### YUM repo

YUM repository:

```
[dkron]
name=Dkron Pro Private Repo
baseurl=https://yum.fury.io/distribworks/
enabled=1
gpgcheck=0
```

Then install: `sudo apt-get install dkron`

This will start Dkron as a system service and the place example configuration file under `/etc/dkron/dkron.yml`

## Running in Docker

Dkron provides an official Docker image vi Dockerhub that can be used for deployment on any system running Docker.

### Launching Dkron on a new container

Here’s a quick one-liner to get you off the ground (please note, we recommend further configuration for production deployments below):

```
docker run -d -p 8080:8080 --name dkron dkron/dkron agent --server
```

This will launch a Dkron server on port 8080 by default. You can use docker logs -f dkron to follow the rest of the initialization progress. Once the Dkron startup completes you can access the app at localhost:8080

Since Docker containers have their own ports and we just map them to the system ports as needed it’s easy to move Dkron onto a different system port if you wish. For example running Dkron on port 12345:

```
docker run -d -p 12345:8080 --name dkron dkron/dkron
```

### Mounting a mapped file storage volume

In its default configuration Dkron uses the local filesystem to run a BoltDB embedded database to store its own application data. The end result is that your Dkron application data will be on disk inside your container and lost if you ever remove the container.

To persist your data outside of the container and make it available for use between container launches we can mount a local file path inside our container.

```
docker run -d -p 8080:8080 -v ~/dkron.db:/dkron.db --name dkron dkron/dkron agent --server
```

Now when you launch your container we are mounting that file from our local filesystem into the container.
