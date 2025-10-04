# AGENTS instructions

This file provides guidance to agents when working with code in this repository.

## Project Overview

Dkron is a distributed, fault-tolerant job scheduling system (distributed cron) written in Go. It uses:
- **Raft protocol** (HashiCorp Raft) for consensus and leader election
- **Serf** (HashiCorp) for cluster membership and failure detection via gossip protocol
- **BuntDB** for embedded storage
- **gRPC** for inter-node communication
- **Gin** for HTTP/REST API
- **React Admin** for the web UI

The system is inspired by Google's "Reliable Cron across the Planet" whitepaper and Airbnb's Chronos.

## Development Commands

### Building
```bash
# Build the main binary (requires UI to be built first)
make main

# Build main binary manually
go build -tags=hashicorpmetrics main.go

# Install built-in plugins (executors and processors)
GOBIN=`pwd` go install ./builtin/...

# Clean build artifacts
make clean
```

### Testing
```bash
# Run full test suite in Docker
make test

# Run tests locally
make localtest
# Or directly:
go test -v ./...

# Run specific test
go test -v ./dkron -run TestJobName
```

### Frontend Development
```bash
# Install dependencies (uses Bun)
cd ui && bun install

# Start development server
cd ui && npm start

# Build UI for production (outputs to dkron/ui-dist)
cd ui && yarn build --out-dir ../dkron/ui-dist
# Or use make:
make ui
```

### Docker Development
```bash
# Start development cluster with live reload
docker compose -f docker-compose.dev.yml up

# Start production-like cluster
docker compose up -d

# Scale the cluster
docker compose up -d --scale dkron-server=4
docker compose up -d --scale dkron-agent=10

# Access UI at http://localhost:8080/ui
```

### Protocol Buffers
```bash
# Generate Go code from proto files
make proto
# Or manually:
protoc -I proto/ --go_out=types --go_opt=paths=source_relative \
  --go-grpc_out=types --go-grpc_opt=paths=source_relative proto/dkron.proto
```

### Code Generation
```bash
# Generate API client from OpenAPI spec
make client
```

## Architecture

### Core Components

**Agent** (`dkron/agent.go`)
- Main orchestrator that combines all subsystems
- Manages both server (Raft leader/follower) and agent (execution) roles
- Contains: Serf cluster, Raft consensus, gRPC server/client, HTTP API, scheduler, storage, plugins

**Storage** (`dkron/store.go`)
- Interface: `Storage` defines operations (Get/SetJob, Get/SetExecution, etc.)
- Implementation: `Store` uses BuntDB (embedded key-value store)
- Stores jobs and executions with max 100 executions per job
- Key prefixes: `jobs/` and `executions/`

**Leader Election** (`dkron/leader.go`)
- Raft-based leader election
- Only the leader schedules jobs (prevents duplicate executions)
- Leader runs the `Scheduler` which manages cron timers

**Job Scheduling** (`dkron/scheduler.go`, `dkron/job.go`)
- Cron-based scheduling using `robfig/cron/v3`
- Extended cron syntax support via `extcron/` package
- Jobs can have parent-child dependencies
- Concurrency policies: `allow` or `forbid`
- Hash symbol (`~`) in schedules is replaced with hash of job name for distribution

**Execution** (`dkron/run.go`, `dkron/execution.go`)
- Leader determines target nodes based on job tags
- Nodes execute jobs via executor plugins
- Results streamed back to leader via gRPC
- Execution status: success, failed, running, partially_failed

**Cluster Communication**
- `dkron/serf.go`: Serf for membership, failure detection, and event propagation
- `dkron/grpc.go`: gRPC server for job execution coordination
- `dkron/grpc_client.go`: gRPC client for communicating with other nodes
- `dkron/api.go`: REST API using Gin framework

**Plugin System** (`plugin/`)
- Two plugin types: **Executors** (run jobs) and **Processors** (process outputs)
- Uses HashiCorp go-plugin for RPC-based plugins
- Built-in executors: shell, http, grpc, kafka, nats, rabbitmq, gcppubsub
- Built-in processors: log, files, syslog, fluent
- Plugins live in `builtin/bins/`

### Key Data Structures

**Job** (`dkron/job.go`)
- Name (unique ID), Schedule (cron expression), Timezone
- Executor and ExecutorConfig
- Parent job for dependencies
- Tags for targeting specific nodes
- Concurrency policy, retries, timeout

**Execution** (`dkron/execution.go`)
- Links to Job, records start/finish times
- NodeName (where it ran), Output, Success/Error status

### Configuration

Config loaded via Viper from:
- Config files: `/etc/dkron/dkron.yaml`, `$HOME/.dkron/dkron.yaml`, `./config/dkron.yaml`
- Environment variables: `DKRON_*` (e.g., `DKRON_NODE_NAME`)
- CLI flags (via Cobra in `cmd/`)

Key settings: `node-name`, `bind-addr`, `server`, `bootstrap-expect`, `data-dir`, `tags`

### Module Structure

- `cmd/`: CLI commands (agent, keygen, version, etc.)
- `dkron/`: Core library (agent, storage, API, scheduler, clustering)
- `builtin/bins/`: Built-in executor and processor plugin binaries
- `plugin/`: Plugin interfaces and serving logic
- `proto/`, `types/`: Protobuf definitions and generated code
- `extcron/`: Extended cron parser supporting simple syntax
- `ntime/`: Nullable time type
- `logging/`: Custom logging setup
- `client/`: Auto-generated API client
- `ui/`: React Admin frontend (TypeScript/Vite)

## Code Patterns

### Working with Storage
Always use the `Storage` interface. Wrap operations in context and check for `buntdb.ErrNotFound` when getting jobs/executions.

### Leader Operations
Check `a.IsLeader()` before performing leader-only operations like job scheduling or applying Raft logs.

### gRPC Communication
Use `GRPCClient.ExecutionDone()` to report execution results to the leader. The leader uses `AgentRun()` to trigger job execution on target nodes.

### Plugin Development
Implement `plugin.Executor` or `plugin.Processor` interfaces. Use `plugin.Serve()` to serve the plugin. See existing plugins in `builtin/bins/` for examples.

## Important Notes

- The system uses both Raft (for consistency) and Serf (for membership)
- Jobs are only scheduled by the leader to avoid duplicates
- Executions are tracked with a sliding window (max 100 per job)
- The UI is embedded in the binary via Go embed (`dkron/ui-dist`)
- Plugin communication uses HashiCorp's go-plugin over gRPC
- Module path is `github.com/distribworks/dkron/v4` (note the v4 major version)
