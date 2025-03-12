# RabbitMQ Executor

A basic RabbitMQ executor that produces a message on a RabbitMQ server/cluster.

## Configuration

Params

```
url:                         URL of the RabbitMQ server/cluster with VHost. Example: amqp://guest:guest@localhost:5672/
queue.name:                  The name of the queue to publish the message to. Required
queue.durable:               Whether the queue is durable or not. Optional, defaults to false
queue.auto_delete:           Whether the queue is auto-deleted after publishing the message. Optional, defaults to false
queue.exclusive:             Whether the queue is exclusive or not. Optional, defaults to false
message.content_type:        The content type of the message
message.delivery_mode:       Message delivery mode. 1 for non-persistent, 2 for persistent. Optional, defaults to 0
message.messageId:           The message ID
message.body:                The actual message body in desired format that will be sent to the queue
message.base64:              Encoded message body in base64 format. Optional, but should not be set if message.body is set.
```

Example:

```json
{
  "executor": "rabbitmq",
  "executor_config": {
    "url": "amqp://guest:guest@localhost:5672/",
    "queue.name": "test",
    "queue.create": "true",
    "queue.durable": "true",
    "queue.auto_delete": "false",
    "queue.exclusive": "false",
    "message.content_type": "application/json",
    "message.delivery_mode": "2",
    "message.messageId": "4373732772",
    "message.body": "{\"key\":\"value\"}"
  }
}
```