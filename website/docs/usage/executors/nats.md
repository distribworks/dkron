# NATS Executor

The NATS executor sends a message to a NATS server/cluster.

Currently, only username/password authentication is supported.

## Configuration

Params:

```
url:      Comma separated list of NATS server URLs
message:  The message to send
subject:  The subject to send the message to
userName: username for authentication
password: password for authentication
debug:    If not empty, turns on debugging. Will log the NATS specific job config and the request sent.
```

Example:

```json
{
  "executor": "nats",
  "executor_config": {
    "url": "tls://nats.demo.io:4443",
    "message": "the message",
    "subject": "myfavoritesubject",
    "userName": "someusername",
    "password": "somepassword"
  }
}
```
