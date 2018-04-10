---
title: Target nodes spec
toc: true
---

## Target nodes spec

You can choose whether a job is run on a node or nodes by specifying tags and a count of target nodes having this tag do you want a job to run.

### Examples:

Target all nodes with a tag:

```
{
    "name": "job_name",
    "command": "/bin/true",
    "schedule": "@every 2m",
    "tags": {
        "role": "web"
    }
}
```

Target only two nodes of a group of nodes with a tag:

```
{
    "name": "job_name",
    "command": "/bin/true",
    "schedule": "@every 2m",
    "tags": {
        "role": "web:2"
    }
}
```

Dkron will try to run the job in the amount of nodes indicated by that count having that tag.
