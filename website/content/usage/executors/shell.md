
---
title: Shell Executor
---

Shell executor runs a system command

## Configuration

Params

```
shell: Run this command using a shell environment
command: The command to run
env: Env vars separated by comma
cwd: Chdir before command run
timeout: Force kill job after specified time. Format: https://golang.org/pkg/time/#ParseDuration.
```

Example

```json
{
  "executor": "shell",
  "executor_config": {
      "shell": "true",
      "command": "my_command",
      "env": "ENV_VAR=va1,ANOTHER_ENV_VAR=var2",
      "cwd": "/app",
      "timeout": "24h"
  }
}
```

## Job execution prometheus metrics
Path: `/metrics`  
Port: 9422  
or configure via environment variable `SHELL_EXECUTOR_PROMETHEUS_PORT`

### Exposed metrics

|   Name                 | Type  |  Description                   |  Labels  |
|------------------------|:------|-------------------------------:|---------:|
| dkron_job_cpu_usage    | gauge | current CPU usage by job       | job_name |
| dkron_job_mem_usage_kb | gauge | current memory consumed by job | job_name |
