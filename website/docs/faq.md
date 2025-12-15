# Frequently Asked Questions (FAQ)

This FAQ addresses the most common issues and questions reported by Dkron users based on GitHub issues and community discussions.

## Performance and Scalability

### Why is my Dkron cluster experiencing performance degradation with high-frequency jobs?

**Issue**: Leader node crashes or becomes unresponsive when running many jobs with frequent schedules (e.g., every second or minute).

**Solutions**:

- Avoid scheduling jobs too frequently (overlapping on the second). Consider batching operations or using a queue system for very high-frequency tasks.
- Offload logs to stdout/file using processors.
- Monitor memory usage and ensure adequate resources are allocated to the leader node.
- Distribute jobs across multiple agent nodes using tags to reduce load on the leader.
- Review the number of concurrent jobs and adjust concurrency policies accordingly.
- Consider upgrading to a more powerful server for the leader node in large-scale deployments.
- Consider upgrading to Dkron Pro using Raft fastlog

**Related**: [Issue #1828](https://github.com/distribworks/dkron/issues/1828), [Issue #1842](https://github.com/distribworks/dkron/issues/1842)

### Does Dkron provide performance benchmark data?

Dkron does not currently provide official benchmark data. Performance varies significantly based on:

- Job frequency and complexity
- Cluster size and configuration
- Hardware resources
- Network latency
- Type of executor used

For high-frequency scheduling scenarios, it's recommended to conduct your own benchmarking in an environment similar to your production setup.

**Related**: [Issue #1839](https://github.com/distribworks/dkron/issues/1839)

## Memory Issues

### Why is Dkron consuming excessive memory or experiencing OOM kills?

**Common causes**:

1. **Memory leaks**: Older versions had memory leak issues that have been addressed in recent releases.
2. **High execution count**: Storing too many execution records (Dkron keeps up to 100 executions per job).
3. **Large job output**: Jobs that produce large amounts of output can consume significant memory.
4. **Too many jobs**: Running thousands of jobs simultaneously increases memory usage.

**Solutions**:

- Update to the latest version of Dkron (4.x series has several memory improvements).
- Use the `files` processor to write job output to disk instead of storing it in memory.
- Limit job output size in executor configurations.
- Monitor memory usage with Prometheus metrics.
- Set appropriate memory limits in Docker/Kubernetes deployments.
- Reduce the number of concurrent jobs or distribute them across more agents.

**Related**: [Issue #1753](https://github.com/distribworks/dkron/issues/1753), [Issue #1566](https://github.com/distribworks/dkron/issues/1566), [Issue #1553](https://github.com/distribworks/dkron/issues/1553)

## Clustering and High Availability

### Why is my agent unable to connect to the cluster?

**Common causes**:

- Network connectivity issues
- Firewall blocking required ports (default: 8946 for Serf, 8190 for Raft)
- Incorrect `retry-join` configuration
- DNS resolution problems
- Encryption key mismatch

**Solutions**:

- Verify network connectivity between nodes using `telnet` or `nc`.
- Check firewall rules allow traffic on ports 8946 (Serf) and 8190 (Raft).
- Use IP addresses instead of hostnames in `retry-join` if DNS is unreliable.
- Ensure all nodes use the same encryption key if encryption is enabled.
- Check logs for specific error messages: `level=error ts=<timestamp> caller=agent.go:xxx`

**Related**: [Issue #1727](https://github.com/distribworks/dkron/issues/1727), [Issue #1702](https://github.com/distribworks/dkron/issues/1702)

### Why does retry-join fail after a server restart?

If using DNS names in `retry-join`, Dkron may cache the old IP addresses. Solutions:

- Use static IP addresses in `retry-join` configuration.
- Ensure DNS updates propagate before restarting nodes.
- Set appropriate DNS TTL values for dynamic environments.

**Related**: [Issue #1253](https://github.com/distribworks/dkron/issues/1253), [Issue #1213](https://github.com/distribworks/dkron/issues/1213)

### How do I recover from "no leader" state after restart?

In single-node mode:

```bash
# Check Raft peers
dkron raft list-peers

# If cluster state is corrupted, backup and remove the data directory
cp -r /var/lib/dkron /var/lib/dkron.backup
rm -rf /var/lib/dkron/*
# Restart dkron
```

In multi-node mode:

- Ensure at least `bootstrap-expect` number of servers are running.
- Check that all servers can communicate on the Raft port (8190).
- Verify there's no network partition.
- Review logs for Raft election errors.

**Related**: [Issue #1529](https://github.com/distribworks/dkron/issues/1529), [Issue #1212](https://github.com/distribworks/dkron/issues/1212)

### Why do servers have inconsistent job data?

**Causes**:

- Network partitions causing split-brain scenarios
- Raft replication issues
- Corrupted data on one or more nodes

**Solutions**:

- Ensure stable network connectivity between all server nodes.
- Verify Raft replication is working: `dkron raft list-peers`
- Check logs for Raft errors: `grep -i raft /var/log/dkron.log`
- In extreme cases, restore from backup or rebuild the cluster.

**Related**: [Issue #1371](https://github.com/distribworks/dkron/issues/1371), [Issue #1337](https://github.com/distribworks/dkron/issues/1337)

## Kubernetes and Container Deployments

### Why does one Dkron pod fail with "panic: log not found" in Kubernetes?

This typically occurs during Raft initialization in StatefulSets when persistent volumes are not properly configured.

**Solutions**:

- Ensure persistent volumes are correctly configured and bound.
- Verify `bootstrap-expect` matches the number of server replicas.
- Check that all pods can resolve each other's DNS names.
- Use `StatefulSet` with headless service for stable network identities.
- Ensure proper initialization order with `podManagementPolicy: Parallel` or readiness probes.

**Related**: [Issue #1557](https://github.com/distribworks/dkron/issues/1557)

### How do I configure log rotation in Docker?

Dkron writes logs to stdout/stderr by default in Docker. To manage logs:

**Docker Compose**:

```yaml
services:
    dkron:
        logging:
            driver: "json-file"
            options:
                max-size: "10m"
                max-file: "3"
```

**Docker CLI**:

```bash
docker run --log-driver json-file \
  --log-opt max-size=10m \
  --log-opt max-file=3 \
  distribworks/dkron:latest
```

To mount logs to the host:

```yaml
volumes:
    - ./logs:/var/log/dkron
```

Then configure Dkron to write to files via `log-level` and redirect output.

**Related**: [Issue #1556](https://github.com/distribworks/dkron/issues/1556)

## Job Execution

### Why does my shell job fail to execute?

**Common causes**:

1. Shell not available in the container/environment
2. Incorrect command syntax
3. Missing environment variables
4. Permission issues

**Solutions**:

- Verify shell is available: `which sh` or `which bash`
- Use absolute paths for commands: `/bin/bash -c "command"`
- Set the `shell` option explicitly in executor config:
    ```json
    {
        "executor": "shell",
        "executor_config": {
            "command": "your-command",
            "shell": "true"
        }
    }
    ```
- Check file permissions and ownership
- Review execution logs for specific error messages

**Related**: [Issue #1775](https://github.com/distribworks/dkron/issues/1775)

### Why does "forbid" concurrency policy not prevent concurrent executions?

In distributed environments, there can be a race condition where multiple executions start before the concurrency check completes. This is a known limitation of distributed scheduling.

**Workarounds**:

- Implement application-level locking (e.g., using Redis, etcd, or database locks).
- Increase job duration to reduce the likelihood of overlapping schedules.
- Use idempotent job designs that can handle concurrent executions safely.
- Add a small delay at job start to check for other running instances.

**Related**: [Issue #1569](https://github.com/distribworks/dkron/issues/1569)

### Why does `cmd.Process.Kill()` not clean up child processes?

On Unix systems, killing the parent process doesn't automatically kill child processes. This can leave orphaned processes running.

**Solution**:
Use process groups to ensure all child processes are terminated:

```bash
# In your script, create a process group
set -m  # Enable job control
trap 'kill -- -$$' EXIT INT TERM

# Your commands here
```

Or use the `timeout` command:

```bash
timeout 300 your-long-running-command
```

**Related**: [Issue #1385](https://github.com/distribworks/dkron/issues/1385)

## Job Configuration

### Why does `runoncreate` block the API call?

When `runoncreate` is set to `true`, the job creation API waits for the first execution to complete before returning. This is intentional behavior for validation purposes.

**Workaround**:
If you need immediate API response:

1. Create the job with `runoncreate: false`
2. Make a separate API call to trigger execution: `POST /v1/jobs/{job_name}/run`

**Related**: [Issue #1268](https://github.com/distribworks/dkron/issues/1268)

## Web UI

### Why does the UI show "Empty output" for running jobs?

The UI displays execution output after the job completes. For long-running jobs, output is not streamed in real-time.

**Solutions**:

- Wait for job completion to see output
- Use the `log` or `files` processor to write output to a log aggregation system
- Use the API to fetch execution details: `GET /v1/executions/{execution_id}`
- Check agent logs directly on the executing node

**Related**: [Issue #1742](https://github.com/distribworks/dkron/issues/1742)

### Why does a running job show "failed" status with strange finished_at date?

This is a UI display bug that has been addressed in recent versions. The job is still running, but the UI incorrectly shows it as failed.

**Solution**:

- Update to the latest version (4.x)
- Check actual job status via API: `GET /v1/jobs/{job_name}`
- Monitor execution status: `GET /v1/executions/{execution_id}`

**Related**: [Issue #1817](https://github.com/distribworks/dkron/issues/1817)

### Why is the job ID missing on small screens?

This was a responsive design issue in the Web UI that has been fixed in version 4.x.

**Solution**:
Update to the latest version of Dkron.

**Related**: [Issue #1583](https://github.com/distribworks/dkron/issues/1583)

## Plugins and Executors

### Why doesn't the GRPC executor close connections?

Older versions had connection leak issues. Solutions:

- Update to the latest version
- Implement connection pooling in your gRPC service
- Set appropriate timeout values in executor config

**Related**: [Issue #1236](https://github.com/distribworks/dkron/issues/1236)

## API

### Why does the restore API not work?

The restore endpoint requires proper backup format and permissions.

**Usage**:

```bash
# Create backup
curl http://localhost:8080/v1/backup > backup.json

# Restore backup
curl -X POST http://localhost:8080/v1/restore \
  -H "Content-Type: application/json" \
  -d @backup.json
```

**Note**: Restore should be performed on a clean cluster or during maintenance windows.

**Related**: [Issue #1376](https://github.com/distribworks/dkron/issues/1376)

## Configuration and Setup

### How do I set up Dkron in a dual data center configuration?

Dkron uses Raft which requires a majority of nodes to be available. For multi-datacenter setups:

**Option 1: Single Raft cluster across DCs** (not recommended)

- High latency can cause leader election issues
- Network partitions can cause outages

**Option 2: Separate clusters per DC** (recommended)

- Run independent Dkron clusters in each DC
- Use external synchronization for job definitions
- Configure failover at the application level

**Option 3: Use Pro version features**

- Enhanced replication capabilities
- Better multi-datacenter support

**Related**: [Issue #1363](https://github.com/distribworks/dkron/issues/1363)

### What are the resource requirements for Dkron?

Requirements vary based on:

- Number of jobs
- Job execution frequency
- Cluster size
- Execution history retention

**Minimum recommendations**:

- **Server nodes**: 1-2 CPU cores, 512MB-1GB RAM
- **Agent nodes**: 1 CPU core, 256MB-512MB RAM

**Production recommendations**:

- **Server nodes**: 2-4 CPU cores, 2-4GB RAM
- **Agent nodes**: 2 CPU cores, 1-2GB RAM
- **Storage**: 10GB+ for execution history

Scale vertically for the leader node and horizontally for agents.

**Related**: [Issue #1560](https://github.com/distribworks/dkron/issues/1560)

## Troubleshooting

### How do I debug "invalid memory address or nil pointer dereference" errors?

This panic typically indicates a bug in Dkron or a plugin. Steps to debug:

1. **Update to the latest version** - Many nil pointer issues have been fixed
2. **Check logs** for the full stack trace
3. **Identify the context** - What operation was being performed?
4. **Reproduce** - Can you consistently trigger the issue?
5. **Report** - If it persists, file a GitHub issue with:
    - Dkron version
    - Configuration
    - Steps to reproduce
    - Full logs and stack trace

**Related**: [Issue #1388](https://github.com/distribworks/dkron/issues/1388), [Issue #1265](https://github.com/distribworks/dkron/issues/1265)

### Where can I find more help?

- **Documentation**: https://dkron.io/docs
- **GitHub Issues**: https://github.com/distribworks/dkron/issues
- **Discussions**: https://github.com/distribworks/dkron/discussions
- **Commercial Support**: Available for Dkron Pro - see https://dkron.io/pro

## Best Practices

### General recommendations for production deployments:

1. **High Availability**: Run at least 3 server nodes (odd number for Raft quorum)
2. **Monitoring**: Enable Prometheus metrics and set up alerts
3. **Backups**: Regularly backup job configurations using the backup API
4. **Updates**: Keep Dkron updated to get bug fixes and performance improvements
5. **Resource Limits**: Set appropriate memory/CPU limits in container environments
6. **Job Design**:
    - Use idempotent operations
    - Implement proper error handling
    - Keep job execution time reasonable (< 1 hour)
    - Avoid very high-frequency schedules (< 1 minute)
7. **Security**:
    - Enable TLS for API communication
    - Use firewall rules to restrict access
    - Consider using ACLs (Pro version)
8. **Logging**: Use processors (log, files, syslog) to persist job output
9. **Testing**: Test job configurations in a staging environment first
10. **Documentation**: Document your job configurations and dependencies
