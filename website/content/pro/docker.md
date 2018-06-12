---
title: Docker executor
---

Docker executor can run docker jobs

## Configuration

To run a docker job create a job config with the following executor:

Example:

```json
...
"executor": "docker"
"executor_config: {
    "image": "alpine",
    "network": "BRIDGE",
    "volumes": [
      {
        "containerPath": "/var/log/",
        "hostPath": "/logs/",
        "mode": "RW"
      }
    ]
  },
  "cpus": "0.5",
  "mem": "512",
  "fetch": [],
  "command": "while sleep 10; do date =u %T; done"
}
...
```
