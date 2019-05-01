---
title: AWS ECS Executor
---

The ECS exeutor is capable of launching tasks in ECS clusters, then listen to a stream of CloudWatch Logs and return the output.

To configure a job to be run in ECS, the executor needs a JSON Task definition template or an already defined task in ECS.

To allow the ECS Task runner to run tasks, the machine running Dkron needs to have the appropriate permissions configured in AWS IAM:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1460720941000",
            "Effect": "Allow",
            "Action": [
                "ecs:RunTask",
                "ecs:DescribeTasks",
                "ecs:DescribeTaskDefinition",
                "logs:FilterLogEvents",
                "logs:DescribeLogStreams",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```

To configure a job to be run with the ECS executor:

Example using an existing taskdef

```json
{
  "executor": "ecs",
  "executor_config": {
    "taskdefName": "mytaskdef-family",
    "region": "eu-west-1",
    "cluster": "default",
    "env": "ENVIRONMENT=variable",
    "service": "mycontainer",
    "overrides": "echo,\"Hello from dkron\""
  }
}
```

Example using a provided taskdef

```json
{
  "executor": "ecs",
  "executor_config": {
    "taskdefBody": "{\"containerDefinitions\": [{\"essential\": true,\"image\": \"hello-world\",\"memory\": 100,\"name\": \"hello-world\"}],\"family\": \"helloworld\"}",
    "region": "eu-west-1",
    "cluster": "default",
    "fargate": "yes",
    "env": "ENVIRONMENT=variable",
    "maxAttempts": "5000"
  }
}
```

This is the complete list of configuration parameters of the plugin:

```
taskdefBody
taskdefName
region
cluster
logGroup
fargate
securityGroup
subnet
env
service
overrides
maxAttempts // Defaults to 2000, will perform a check every 6s * 2000 times waiting a total of 12000s or 3.3h
```
