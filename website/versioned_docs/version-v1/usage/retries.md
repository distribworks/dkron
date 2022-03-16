---
title: Job retries
---

Jobs can be configured to retry in case of failure.

## Configuration

```json
{
  "name": "job1",
  "schedule": "@every 10s",
  "executor": "shell",
  "executor_config": {
    "command": "echo \"Hello from parent\""
  },
  "retries": 5
}
```

In case of failure to run the job in one node, it will try to run the job again in that node until the retries count reaches the limit.

