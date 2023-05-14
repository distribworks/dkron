---
title: Job chaining
---

## Job chaining

You can set some jobs to run after other job is executed. To setup a job that will be executed after any other given job, just set the `parent_job` property when saving the new job.

The dependent job will be executed after the main job finished a successful execution.

Child jobs schedule property will be ignored if it's present.

Take into account that parent jobs must be created before any child job.

Example:

```json
{
  "name": "job1",
  "schedule": "@every 10s",
  "executor": "shell",
  "executor_config": {
    "command": "echo \"Hello from parent\""
  }
}

{
  "name": "child_job",
  "parent_job": "job1",
  "executor": "shell",
  "executor_config": {
    "command": "echo \"Hello from child\""
  }
}
```
