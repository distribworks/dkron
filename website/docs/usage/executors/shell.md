# Shell Executor

The Shell executor is one of the most versatile executors in Dkron, allowing you to run any system command or script on target nodes.

## Overview

The shell executor runs commands on the target node's operating system. It can:

- Execute system commands
- Run scripts in various languages (bash, python, etc.)
- Perform file system operations
- Interact with local services and processes

## Configuration Parameters

| Parameter | Required | Description |
|-----------|:--------:|-------------|
| `shell` | No | When "true", runs the command in a shell environment (bash/cmd). Default: "false" |
| `command` | Yes | The command or script to execute |
| `env` | No | Environment variables in the format "KEY1=value1,KEY2=value2" |
| `cwd` | No | The working directory to run the command from |
| `timeout` | No | Maximum execution time after which the job is forcefully terminated |

## Basic Usage Examples

### Simple Command Execution

```json
{
  "executor": "shell",
  "executor_config": {
    "command": "echo Hello Dkron"
  }
}
```

### Using Shell Features

```json
{
  "executor": "shell",
  "executor_config": {
    "shell": "true",
    "command": "ps aux | grep nginx | wc -l"
  }
}
```

### Setting Environment Variables

```json
{
  "executor": "shell",
  "executor_config": {
    "shell": "true",
    "command": "python /scripts/process_data.py",
    "env": "DATA_PATH=/data/incoming,LOG_LEVEL=info,API_KEY=secret123"
  }
}
```

### Changing Working Directory

```json
{
  "executor": "shell",
  "executor_config": {
    "shell": "true",
    "command": "./run_backup.sh",
    "cwd": "/opt/backups"
  }
}
```

### Setting a Timeout

```json
{
  "executor": "shell",
  "executor_config": {
    "command": "/opt/scripts/long_running_task.sh",
    "timeout": "1h30m"
  }
}
```

## Advanced Examples

### Running a Multi-line Script

```json
{
  "executor": "shell",
  "executor_config": {
    "shell": "true",
    "command": "#!/bin/bash\necho 'Starting job'\ncd /tmp\ndate > last_run.txt\necho 'Job completed'"
  }
}
```

### Database Backup Example

```json
{
  "executor": "shell",
  "executor_config": {
    "shell": "true",
    "command": "pg_dump -U postgres -d mydb | gzip > /backups/mydb_$(date +%Y%m%d_%H%M%S).sql.gz",
    "env": "PGPASSWORD=securepassword",
    "timeout": "30m"
  }
}
```

### System Health Check with Exit Code

```json
{
  "executor": "shell",
  "executor_config": {
    "shell": "true",
    "command": "if [ $(df -h | grep '/dev/sda1' | awk '{print $5}' | tr -d '%') -gt 90 ]; then echo 'Disk space critical'; exit 1; else echo 'Disk space OK'; fi"
  }
}
```

## Monitoring and Metrics

The shell executor exposes Prometheus metrics on port 9422 (configurable with `SHELL_EXECUTOR_PROMETHEUS_PORT` environment variable):

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `dkron_job_cpu_usage` | gauge | Current CPU usage by job | job_name |
| `dkron_job_mem_usage_kb` | gauge | Current memory consumed by job (KB) | job_name |

## Exit Codes and Success/Failure

By default, an exit code of 0 indicates success, while any non-zero exit code indicates failure. You can customize which exit codes are considered successful in the job definition:

```json
{
  "success_count": "0,1,2",  // Exit codes 0, 1, and 2 will be considered successful
  "executor": "shell",
  "executor_config": {
    "command": "/opt/scripts/check.sh"
  }
}
```

## Security Considerations

The shell executor runs commands with the same permissions as the Dkron process. Consider the following security best practices:

1. **Principle of Least Privilege**: Run Dkron with a dedicated user that has only the permissions required for its jobs
2. **Avoid Sensitive Data in Commands**: Use environment variables for sensitive values instead of embedding them in commands
3. **Input Validation**: Validate any dynamic parts of commands, especially if they come from external sources
4. **Output Sanitization**: Be cautious when using command output, particularly if it's displayed in the UI or logs

## Troubleshooting

Common issues and solutions when working with the shell executor:

### Command Not Found

If you see "command not found" errors:
- Verify the command is installed on the target node
- Use absolute paths to executables (e.g., `/usr/bin/python` instead of just `python`)
- Check if the command is in the PATH of the user running Dkron

### Permission Issues

If you encounter permission errors:
- Check if the Dkron user has the necessary permissions
- Adjust file/directory permissions as needed
- Consider using `sudo` if appropriate (requires sudo configuration)

### Timeouts

If jobs are timing out:
- Adjust the `timeout` parameter to match the expected runtime
- Optimize the command to run more efficiently
- Consider breaking large jobs into smaller, chained jobs

### Environment Variables

If environment variables aren't working:
- Ensure the format is correct (`KEY=value,KEY2=value2`)
- Remember that values are passed as strings
- For complex values with spaces or special characters, consider using files instead

## Additional Resources

- [Process Output with Processors](/docs/usage/processors)
- [Job Chaining for Complex Workflows](/docs/usage/chaining)
- [Security Best Practices](/docs/usage/concepts#security-considerations)
