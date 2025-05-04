---
title: Metrics
---

Dkron has the ability to send metrics to Statsd for dashboards and historical reporting or provide prometheus format metrics via the api. It sends job processing metrics, golang, and serf metrics.

## Configuration

### Statsd

Add this in your yaml config file to enable statsd metrics.

```yaml
statsd-addr: "localhost:8125"
# Or for datadog statsd
dog-statsd-addr: "localhost:8125"
```

### Prometheus

Add this to your yaml config file to enable serving prometheus metrics at the endpoint `/metrics`

```yaml
enable-prometheus: true
```

Additionally, in your Prometheus config file (prometheus.yml), add the following to link dkron metric endpoint
```yaml
scrape_configs:
  ... #initial configuration
  
  - job_name: "dkron_metrics"
    # metrics_path defaults to '/metrics'
    static_configs:
      - targets: ["localhost:8080"]
```

## Monitoring Dashboard Setup

Setting up a monitoring dashboard with metrics from Dkron is highly recommended for production deployments. You have several options:

1. **Grafana + Prometheus**: Create dashboards with job execution success rates, execution timing trends, and system health metrics
2. **DataDog**: Use built-in dashboard templates with the DogStatsD integration
3. **Custom StatsD Dashboards**: Configure with any StatsD-compatible visualization tool

Example Grafana dashboard panels for monitoring Dkron might include:
- Job execution success rate over time
- Average execution duration by job
- System resource utilization
- Cluster health status

## Metrics Reference

Dkron emits several categories of metrics to help you monitor both the application and your job executions:

### Agent Event Metrics

These metrics track internal events within Dkron agents:

| Metric | Description |
|--------|-------------|
| `dkron.agent.event_received.query_execution_done` | Count of completed job execution events received by the agent |
| `dkron.agent.event_received.query_run_job` | Count of job run requests received by the agent |

### Network Communication Metrics

These metrics help monitor the health of inter-node communication:

| Metric | Description |
|--------|-------------|
| `dkron.memberlist.gossip` | Count of gossip protocol messages exchanged |
| `dkron.memberlist.probeNode` | Count of node health probe checks |
| `dkron.memberlist.pushPullNode` | Count of anti-entropy sync operations |
| `dkron.memberlist.tcp.accept` | Count of accepted TCP connections |
| `dkron.memberlist.tcp.connect` | Count of initiated TCP connections |
| `dkron.memberlist.tcp.sent` | Count and bytes of TCP packets sent |
| `dkron.memberlist.udp.received` | Count and bytes of UDP packets received |
| `dkron.memberlist.udp.sent` | Count and bytes of UDP packets sent |

### gRPC Service Metrics

These metrics track the internal RPC communication:

| Metric | Description |
|--------|-------------|
| `dkron.grpc.call_execution_done` | Count and timing of execution completion RPC calls |
| `dkron.grpc.call_get_job` | Count and timing of job retrieval RPC calls |
| `dkron.grpc.execution_done` | Count of completed job executions |
| `dkron.grpc.get_job` | Count of job information retrievals |

### Runtime Metrics

These metrics provide insights into the Go runtime health:

| Metric | Description |
|--------|-------------|
| `dkron.runtime.alloc_bytes` | Current bytes allocated by the application |
| `dkron.runtime.free_count` | Count of memory free operations |
| `dkron.runtime.gc_pause_ns` | Duration of the last garbage collection pause in nanoseconds |
| `dkron.runtime.heap_objects` | Number of objects in the heap |
| `dkron.runtime.malloc_count` | Count of memory allocation operations |
| `dkron.runtime.num_goroutines` | Number of goroutines currently running |
| `dkron.runtime.sys_bytes` | Total bytes of memory obtained from the OS |
| `dkron.runtime.total_gc_pause_ns` | Total time spent in GC pauses |
| `dkron.runtime.total_gc_runs` | Total number of completed GC cycles |

### Serf Metrics

These metrics track the cluster membership and failure detection system:

| Metric | Description |
|--------|-------------|
| `dkron.serf.coordinate.adjustment_ms` | Time adjustment for coordinate system in milliseconds |
| `dkron.serf.msgs.received` | Count of messages received |
| `dkron.serf.msgs.sent` | Count of messages sent |
| `dkron.serf.queries` | Count of queries processed |
| `dkron.serf.queries.execution_done` | Count of execution completion queries |
| `dkron.serf.queries.run_job` | Count of job run queries |
| `dkron.serf.query_acks` | Count of query acknowledgments |
| `dkron.serf.query_responses` | Count of responses to queries |
| `dkron.serf.queue.Event` | Count of events in processing queue |
| `dkron.serf.queue.Intent` | Count of intent messages in queue |
| `dkron.serf.queue.Query` | Count of queries in processing queue |

## Alerting on Metrics

To set up effective alerting based on Dkron metrics, consider these recommendations:

1. **Job Failure Alerts**: Monitor `dkron.serf.queries.execution_done` with error status
2. **Cluster Health**: Set alerts on node count changes or consistent gossip failures
3. **Performance Degradation**: Watch for increases in job execution time trends
4. **Resource Constraints**: Monitor `dkron.runtime.gc_pause_ns` and other runtime metrics for signs of resource pressure

For Prometheus users, example alerting rules:

```yaml
groups:
- name: dkron_alerts
  rules:
  - alert: DkronJobFailures
    expr: increase(dkron_job_executions_failed_total[1h]) > 3
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Multiple job failures detected"
      description: "There have been more than 3 job failures in the last hour"
```
