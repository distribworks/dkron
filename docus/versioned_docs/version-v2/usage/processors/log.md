---
title: Log Processor
---

Log processor writes the execution output to stdout/stderr

## Configuration

Parameters

`forward: Forward the output to the next processor`

Example

```json
{
    "name": "job_name",
    "command": "echo 'Hello log'",
    "schedule": "@every 2m",
    "tags": {
        "role": "web"
    },
    "processors": {
        "log": {
            "forward": true
        }
    }
}
```
