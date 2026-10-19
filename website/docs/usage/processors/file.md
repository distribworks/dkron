---
title: File Processor
---

File processor saves the execution output to a single log file in the specified directory

The output is written to the file while the job runs, so the file holds the complete output, however large, and can be followed live (e.g. with `tail -f`). The output stored with the execution is capped at 256000 bytes, the log file is not.

If the output of an execution was not streamed to the node that processes it, for example after a leader change during the execution, the stored (capped) output is written instead.

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
            "forward": "true"
        }
    }
}
```
