---
title: Job metadata
---

## Job metadata

Jobs can have an optional extra property field called `metadata` that allows to set arbitrary tags to jobs and query the jobs using the API:

```json
{
    "name": "job_name",
    "command": "/bin/true",
    "schedule": "@every 2m",
    "metadata": {
        "user_id": "12345"
    }
}
```

And then query the API to get only the results needed:

```
$ curl http://localhost:8080/v1/jobs --data-urlencode "metadata[user_id]=12345"`
```
