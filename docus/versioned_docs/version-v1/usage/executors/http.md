---
title: HTTP Executor
---

HTTP executor can send a request to an HTTP endpoint

## Configuration

Params:

```
method: Request method in uppercase
url: Request url
headers: Json string, such as "[\"Content-Type: application/json\"]"
body: POST body
timeout: Request timeout, unit seconds
expectCode: Expect response code, such as 200,206
expectBody: Expect response body, support regexp, such as /success/
debug: Debug option, will log everything when this option is not empty
```

Example

```json
{
  "executor": "http",
  "executor_config": {
      "method": "GET",
      "url": "http://example.com",
      "headers": "[]",
      "body": "",
      "timeout": "30",
      "expectCode": "200",
      "expectBody": "",
      "debug": "true"
  }
}
```
