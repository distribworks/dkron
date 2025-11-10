# Raft FastLog

Dkron Pro includes support for the [Raft FastLog](https://github.com/tidwall/raft-fastlog) storage engine, which provides significant performance improvements over the default BoltDB storage.

## Overview

Raft FastLog is a high-performance storage engine designed specifically for Raft consensus logs. It offers:

- **Higher throughput** - Significantly faster write operations
- **Lower latency** - Reduced response times for log operations  
- **Better concurrency** - Improved performance under concurrent load
- **Memory efficiency** - Optimized memory usage patterns

## Performance Benefits

FastLog typically provides:
- 10-100x faster write performance compared to BoltDB
- Reduced memory allocation and garbage collection pressure
- Better scaling characteristics under high load
- Lower CPU utilization for log operations

## Configuration

To enable Raft FastLog, use the `--fast` command line option when starting Dkron Pro:

```bash
dkron agent --fast
```

## When to Use FastLog

FastLog is recommended for:
- High-frequency job scheduling scenarios
- Large clusters with many nodes
- Environments requiring low-latency job execution
- Production deployments with performance requirements

## Considerations

- FastLog requires Dkron Pro license
- Log files are not compatible between FastLog and BoltDB engines
- Ensure adequate disk space for log storage
- Monitor system resources during initial deployment

## Migration

When switching from BoltDB to FastLog:

1. Stop all Dkron nodes in the cluster
2. Backup existing data
3. Start nodes with `--fast` flag
4. The cluster will rebuild its state from scratch

**Note**: Migration requires cluster restart and state rebuild. Plan accordingly for production environments.
