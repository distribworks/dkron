# Installation

This guide covers different methods of installing Dkron for development, testing, and production environments.

## Quick Start Methods

### Running the binary

Download the packaged archive for your platform from the [downloads page](https://github.com/distribworks/dkron/releases) and extract the package to a shared location in your drive, like /opt/local/bin.

Run Dkron with default setting: `dkron agent --server --bootstrap-expect=1`

Navigate to http://localhost:8080/ui

### Installing from package repositories

#### Debian/Ubuntu

APT repository: 
```
deb [trusted=yes] https://repo.distrib.works/apt/ /
```

Setup and install:
```bash
# Add the repository
echo "deb [trusted=yes] https://repo.distrib.works/apt/ /" | sudo tee /etc/apt/sources.list.d/dkron.list

# Update package list
sudo apt-get update

# Install Dkron
sudo apt-get install dkron
```

#### RHEL/CentOS/Fedora

YUM repository:

```
[dkron]
name=Dkron Repository
baseurl=https://repo.distrib.works/yum/
enabled=1
gpgcheck=0
```

Setup and install:
```bash
# Create repo file
sudo tee /etc/yum.repos.d/dkron.repo << EOF
[dkron]
name=Dkron Repository
baseurl=https://repo.distrib.works/yum/
enabled=1
gpgcheck=0
EOF

# Install Dkron
sudo yum install dkron
```

After package installation, Dkron will be started as a system service and an example configuration file will be placed under `/etc/dkron/dkron.yml`.

## Running in Docker

Dkron provides official Docker images via Docker Hub that can be used for deployment on any system running Docker.

:::info
If you only plan to use the built-in executors, `http` and `shell` you can use the Dkron Light edition that only includes a single binary as the plugins are built-in.
:::

### Launching Dkron as a new container

Here's a quick one-liner to get you off the ground (please note, we recommend further configuration for production deployments below):

```
docker run -d -p 8080:8080 --name dkron dkron/dkron agent --server --bootstrap-expect=1 --node-name=node1
```

This will launch a Dkron server on port 8080 by default. You can use `docker logs -f dkron` to follow the rest of the initialization progress. Once the Dkron startup completes you can access the app at localhost:8080

Since Docker containers have their own ports and we just map them to the system ports as needed it's easy to move Dkron onto a different system port if you wish. For example running Dkron on port 12345:

```
docker run -d -p 12345:8080 --name dkron dkron/dkron agent --server --bootstrap-expect=1 --node-name=node1
```

### Docker Image Variants

Dkron offers two Docker image variants:

1. **Standard Image** (`dkron/dkron`): Includes all executors and processors
2. **Light Image** (`dkron/dkron:light`): Contains only the main binary with built-in shell and HTTP executors

### Mounting a mapped file storage volume

Dkron uses the local filesystem for storing the embedded database to store its own application data and the Raft protocol log. The end result is that your Dkron data will be on disk inside your container and lost if you ever remove the container.

To persist your data outside of the container and make it available for use between container launches we can mount a local path inside our container.

```
docker run -d -p 8080:8080 -v ~/dkron.data:/dkron.data --name dkron dkron/dkron agent --server --bootstrap-expect=1 --data-dir=/dkron.data
```

Now when you launch your container we are mounting that folder from our local filesystem into the container.

### Running a Dkron cluster in Docker Compose

For development or simple multi-node setups, Docker Compose offers an easy way to run a Dkron cluster:

```yaml
version: '3'
services:
  dkron-leader:
    image: dkron/dkron
    ports:
      - "8080:8080"
    command: agent --server --bootstrap-expect=3 --node-name=leader --retry-join=dkron-follower-1 --retry-join=dkron-follower-2 --data-dir=/dkron.data
    volumes:
      - leader-data:/dkron.data

  dkron-follower-1:
    image: dkron/dkron
    command: agent --server --node-name=follower1 --retry-join=dkron-leader --retry-join=dkron-follower-2 --data-dir=/dkron.data
    volumes:
      - follower1-data:/dkron.data
    depends_on:
      - dkron-leader

  dkron-follower-2:
    image: dkron/dkron
    command: agent --server --node-name=follower2 --retry-join=dkron-leader --retry-join=dkron-follower-1 --data-dir=/dkron.data
    volumes:
      - follower2-data:/dkron.data
    depends_on:
      - dkron-leader

volumes:
  leader-data:
  follower1-data:
  follower2-data:
```

## Production Deployment Considerations

When deploying Dkron in production, consider these best practices:

### System Requirements

- **CPU**: Minimum 2 cores recommended, more for heavy workloads
- **Memory**: Minimum 1GB RAM, 2GB+ recommended for production
- **Disk**: SSD recommended for better performance, especially for the data directory
- **Network**: Low-latency network between cluster nodes

### Cluster Size

- **Minimum Nodes**: 3 server nodes for high availability
- **Recommended**: Use an odd number of server nodes (3, 5, 7) to avoid split-brain scenarios
- **Region Distribution**: For multi-region deployments, use at least 3 servers per region

### Security Recommendations

1. **Network Security**:
   - Restrict access to Dkron API/UI using firewalls
   - Run behind a reverse proxy with TLS termination
   - Use VPN or private network for inter-node communication

2. **Authentication**:
   - Use Dkron Pro for built-in authentication features
   - Consider implementing a reverse proxy with authentication

3. **Encryption**:
   - Enable network encryption with `--encrypt` parameter using a key generated with `dkron keygen`

### Data Directory Considerations

- Place the data directory on persistent, high-performance storage
- Implement regular backups of this directory
- Monitor disk space usage

### Configuration Example for Production

Example production configuration file (`/etc/dkron/dkron.yml`):

```yaml
# Server settings
server: true
bootstrap-expect: 3
data-dir: /var/lib/dkron
log-level: info

# Network settings
bind-addr: "{{ GetPrivateIP }}:8946"
http-addr: ":8080"
advertise-addr: "{{ GetPublicIP }}"
encrypt: "YOUR_ENCRYPTION_KEY_GENERATED_WITH_DKRON_KEYGEN"

# Node identification
node-name: "dkron-prod-1"
datacenter: "dc1"
region: "us-west"

# Clustering
retry-join:
  - "10.0.1.10"
  - "10.0.1.11"
  - "10.0.1.12"
retry-interval: "30s"
raft-multiplier: 1

# Optional metrics for monitoring
statsd-addr: "localhost:8125"
enable-prometheus: true

# Tags for node selection
tags:
  role: "backend"
  env: "production"
```

### Monitoring and Operations

1. **Health Monitoring**:
   - Set up health checks against the API endpoint
   - Monitor node status and leader elections
   - Configure metrics collection with Prometheus or StatsD

2. **Log Management**:
   - Configure log collection and aggregation
   - Set appropriate log levels (`info` for production, `debug` for troubleshooting)

3. **Backup Strategy**:
   - Regular backups of the data directory
   - Document and test restoration procedures

## Next Steps

After installation, follow the [Getting Started guide](/docs/basics/getting-started) to learn how to create and manage jobs with Dkron.
