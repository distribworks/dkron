---
title: Target nodes spec
weight: 10
---

## Target nodes spec

You can choose whether a job is run on a node or nodes by specifying tags and a count of target nodes having this tag do you want a job to run.

{{% notice note %}}
The target node syntax: `[tag-value]:[count]`
{{% /notice %}}

### Examples:

Target all nodes with a tag:

```json
{
    "name": "job_name",
    "command": "/bin/true",
    "schedule": "@every 2m",
    "tags": {
        "role": "web"
    }
}
```

{{<mermaid align="left">}}
graph LR;
    J("Job tags: #quot;role#quot;: #quot;web#quot;") -->|Run Job|N1["Node1 tags: #quot;role#quot;: #quot;web#quot;"]
    J -->|Run Job|N2["Node2 tags: #quot;role#quot;: #quot;web#quot;"]
    J -->|Run Job|N3["Node2 tags: #quot;role#quot;: #quot;web#quot;"]
{{< /mermaid >}}

Target only one nodes of a group of nodes with a tag:

```json
{
    "name": "job_name",
    "command": "/bin/true",
    "schedule": "@every 2m",
    "tags": {
        "role": "web:1"
    }
}
```

{{<mermaid align="left">}}
graph LR;
    J("Job tags: #quot;role#quot;: #quot;web:1#quot;") -->|Run Job|N1["Node1 tags: #quot;role#quot;: #quot;web#quot;"]
    J -.- N2["Node2 tags: #quot;role#quot;: #quot;web#quot;"]
    J -.- N3["Node2 tags: #quot;role#quot;: #quot;web#quot;"]
{{< /mermaid >}}

Dkron will try to run the job in the amount of nodes indicated by that count having that tag.

