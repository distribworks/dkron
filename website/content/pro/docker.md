---
title: Docker executor
---

Docker executor can run docker jobs

## Configuration

To run a docker job create a job config with the following executor:

Example:

```json
  "executor": "docker",
  "executor_config": {
    "image": "alpine", //docker image to use
    "volumes": "/logs:/var/log/", //comma separated list of volume mappings
    "command": "echo \"Hello from dkron\"", //command to pass to run on container
    "env": "ENVIRONMENT=variable" //environment variables to pass to the container
  }
```
