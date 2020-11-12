---
title: Target nodes spec
weight: 10
---

## Target nodes spec

Dkron has the ability to run jobs in specific nodes by leveraging the use of tags. You can choose whether a job is run on a node or group of nodes by specifying tags and a count of target nodes having this tag do you want a job to run.

The target node syntax is:
    
    "tag": "value[:count]"

To achieve this Nodes and Jobs have tags, for example, having a node with the following tags:

```json
{
    "tags": {
        "dc": "dc1",
        "expect": "3",
        "port": "6868",
        "region": "global",
        "role": "dkron",
        "rpc_addr": "10.88.94.129:6868",
        "server": "true",
        "version": "devel",
        "my_role": "web"
    }
}
```

{{% alert info %}}**Tip:** You can specify tags for nodes in the dkron config file or in the command line using `--tags` parameter{{% /alert %}}

Following some examples using different tag combinations:

#### Target all nodes with a tag

```json
{
    "name": "job_name",
    "command": "/bin/true",
    "schedule": "@every 2m",
    "tags": {
        "my_role": "web"
    }
}
```

{{<mermaid align="left">}}
graph LR;
    J("Job tags: #quot;my_role#quot;: #quot;web#quot;") -->|Run Job|N1["Node1 tags: #quot;my_role#quot;: #quot;web#quot;"]
    J -->|Run Job|N2["Node2 tags: #quot;my_role#quot;: #quot;web#quot;"]
    J -->|Run Job|N3["Node2 tags: #quot;my_role#quot;: #quot;web#quot;"]
{{</mermaid>}}

#### Target only one nodes of a group of nodes with a tag

```json
{
    "name": "job_name",
    "command": "/bin/true",
    "schedule": "@every 2m",
    "tags": {
        "my_role": "web:1"
    }
}
```

{{<mermaid align="left">}}
graph LR;
    J("Job tags: #quot;my_role#quot;: #quot;web:1#quot;") -->|Run Job|N1["Node1 tags: #quot;my_role#quot;: #quot;web#quot;"]
    J -.- N2["Node2 tags: #quot;my_role#quot;: #quot;web#quot;"]
    J -.- N3["Node2 tags: #quot;my_role#quot;: #quot;web#quot;"]
{{</mermaid>}}

Dkron will try to run the job in the amount of nodes indicated by that count having that tag.

### Details and limitations

* Tags specified in a Job are combined using `AND`, therefore a job that specifies several tags like:

```json
{
    "tags": {
        "my_role": "web",
        "role": "dkron"
    }
}
```

Will try to run the job in nodes that have all speciefied tags.

There is no limit in the tags that a job can have but having a Job with several tags with count like:

```json
{
    "tags": {
        "my_role": "web:1",
        "role": "dkron:2"
    }
}
```

Will try to run the job in nodes that have all specified tags and using the lowest count. In the last example, it will run in **one** node having `"my_role": "web"` and `"role": "dkron"` tag, even if there is more than one node with these tags.

* In case there is no matching nodes with the specified tags, the job will not run
* In case no tags are specified for a job it will run in all nodes in the cluster
