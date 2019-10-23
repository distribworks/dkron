---
title: Concurrency
toc: true
---

## Concurrency

Jobs can be configured to allow overlapping executions or forbid them. 

Concurrency property accepts two option: 

* **allow** (default): Allow concurrent job executions.
* **forbid**: If the job is already running don't send the execution, it will skip the executions until the next schedule.

Example:

```json
{
  "name": "job1",
  "schedule": "@every 10s",
  "executor": "shell",
  "executor_config": {
    "command": "echo \"Hello from parent\""
  },
  "concurrency": "forbid"
}
```
