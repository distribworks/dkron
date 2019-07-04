---
title: Use with AWS ECS
---

{{% notice note %}}
[Dkron Pro](/products/pro) comes with a [native ECS executor](/pro/ecs) out of the box.
{{% /notice %}}

## Use with Amazon ECS

To use Dkron to schedule jobs that run in containers, a wrapper ECS script is needed.

Install the following snippet in the node that will run the call to ECS

<script src="https://gist.github.com/distribworks/3ac4aae9279d7c68c486fecccc2546cc.js"></script>

### Prerequisites

The node that will run the call to ECS will need to have installed

* AWS cli
* jq

### Example

`ecs-run --cluster cron --task-definition cron-taskdef --container-name cron --region us-east-1 --command "rake foo"`
