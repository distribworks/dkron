# GRPC Executor

GRPC executor can send a request to a GRPC Server

## Requirements

In order to serialize protobufs, the server needs to have the reflection service active.
Without that we cannot get proto descriptors required for serialization.

## Configuration

Params:

```
"url": Required, Request url
"body": Optional, POST body
"timeout": Optional, Request timeout, unit seconds
"expectCode": Optional, One of https://grpc.github.io/grpc/core/md_doc_statuscodes.html
```

Example:

```json
{
  "executor": "http",
  "executor_config": {
    "url": "127.0.0.1:9000/test.TestService/Test",
    "body": "{\"key\": \"value\"}",
    "timeout": "30",
    "expectCode": "0"
  }
}
```
