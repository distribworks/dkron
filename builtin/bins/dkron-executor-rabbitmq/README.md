### RabbitMQ Executor

#### Executor Configuration

The options names are inherited from the [RabbitMQ Publishers](https://www.rabbitmq.com/docs/publishers)

| Option                | Required | Description                   | Default    |
|-----------------------|----------|-------------------------------|------------|
| url                   | yes      | RabbitMQ connection string    | -          |
| exchange              | no       | RabbitMQ Exchange             | ""         |
| queue.name            | yes      | Queue name to send message to | -          |
| queue.create          | no       | Create queue if not exists    | false      |
| queue.durable         | no       | Durable queue                 | false      |
| queue.auto_delete     | no       | Auto delete queue             | false      |
| queue.exclusive       | no       | Exclusive queue               | false      |
| message.content_type  | no       | Message content type          | text/plain |
| message.delivery_mode | no       | Message delivery mode         | 0          |
| message.messageId     | no       | Message id                    | ""         |
| message.body          | yes      | Message body                  | -          |
| message.base64Body    | yes      | Base64 encoded message body   | -          |

#### Example

```shell
curl localhost:8080/v1/jobs -XPOST -d '{
  "name": "job1",
  "schedule": "@every 10s",
  "timezone": "Europe/Berlin",
  "owner": "Platform Team",
  "owner_email": "platform@example.com",
  "disabled": false,
  "tags": {
    "server": "true:1"
  },
  "metadata": {
    "user": "12345"
  },
  "concurrency": "allow",
  "executor": "rabbitmq",
  "executor_config": {
  		"url": "amqp://guest:guest@localhost:5672/",
      "exchange": "amq.default",
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
}'
```

