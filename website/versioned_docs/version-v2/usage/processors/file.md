---
title: File Processor
---

File processor saves the execution output to a single log file in the specified directory

## Configuration

Parameters

```
log_dir: Path to the location where the log files will be saved
forward: Forward log output to the next processor
```

Example

```json
{
    "name": "job_name",
    "command": "echo 'Hello files'",
    "schedule": "@every 2m",
    "tags": {
        "role": "web"
    },
    "processors": {
        "files": {
            "log_dir": "/var/log/mydir",
            "forward": true
        }
    }
}
```
