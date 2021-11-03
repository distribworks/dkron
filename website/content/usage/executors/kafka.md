
---
title: Kafka Executor
---

A basic Kafka executor that produces a message on a Kafka broker.

## Configuration

Params

```
brokerAddress: "IP:port" of the broker
key:           The key of the message to produce
message:       The body of the message to produce
topic:         The Kafka topic for this message
debug:         Turns on debugging output if not empty
```

Example

```json
"executor": "kafka",
"executor_config": {
    "brokerAddress": "localhost:9092",
    "key": "My key",
    "message": "My message",
    "topic": "my_topic"
}
```
