---
title: Job retries
---

Jobs can be configured to retry in case of failure.

## Configuration

```json
{
  ...
  retries: 2,
  ...
}
```

In case of failure to run the job in one node, it will try to run the job again in that node until the retries count reaches the limit.

