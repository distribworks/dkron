# Google Cloud PubSub Executor

A basic GC PubSub executor that produces a message on a Google Cloud PubSub topic.

## Configuration

Params

```
project:                     Google Cloud project ID. Required
topic:                       The topic name for the message. Required
data:                        The actual message body in base64 format. Required
attributes:                  The attributes of the message in JSON format. Optional
```

Example:

```json
{
  "executor": "gcppubsub",
  "executor_config": {
    "project": "project-id",
    "topic": "topic-name",
    "data": "aGVsbG8gd29ybGQ=",
    "attributes": "{\"hello\":\"world\",\"waka\":\"paka\"}"
  }
}
```

