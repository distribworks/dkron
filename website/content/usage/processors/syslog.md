---
title: Syslog Processor
---

Syslog processor writes the execution output to the system syslog daemon

Note: Only work on linux systems

## Configuration

Parameters

`forward: Forward the output to the next processor`

Example

```json
{
    "name": "job_name",
    "command": "echo 'Hello syslog'",
    "schedule": "@every 2m",
    "tags": {
        "role": "web"
    },
    "processors": {
        "syslog": {
            "forward": true
        }
    }
}
```
