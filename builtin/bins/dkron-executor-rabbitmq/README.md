### Create a job to send rabbitmq text Message
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
    "text": "hello world!",
    "queue": "test"
  }
}'
