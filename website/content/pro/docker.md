---
title: Docker executor
---

Docker executor can launch docker based cron jobs using the docker command of the target node.

This executor needs a recent version of docker to be available and configured in the target node.

## Configuration

To run a docker job create a job config with the docker executor as in this example:

```json
{
  "executor": "docker",
  "executor_config": {
    "image": "alpine", //docker image to use
    "volumes": "/logs:/var/log/", //comma separated list of volume mappings
    "command": "echo \"Hello from dkron\"", //command to pass to run on container
    "env": "ENVIRONMENT=variable" //environment variables to pass to the container
  }
}
```
