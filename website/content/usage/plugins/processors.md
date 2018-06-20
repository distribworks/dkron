---
title: Execution Processors
weight: 2
---

## Execution Processors

Processor plugins are called when an execution response has been received. They are passed the resulting execution data and configuration parameters, this plugins can perform a variety of operations with the execution and it's very flexible and per Job, examples of operations this plugins can do:

* Execution output storage, forwarding or redirection.
* Notification
* Monitoring

For example, Processor plugins can be used to redirect the output of a job execution to different targets.

Currently Dkron provides you with some built-in plugins but the list keeps growing. Some of the features previously implemented in the application will be progessively moved to plugins.

### Built-in processors

Dkron provides the following built-in processors:

0. not specified - Store the output in the key value store (Slow performance, good for testing, default method)
0. log - Output the execution log to Dkron stdout (Good performance, needs parsing)
0. syslog - Output to the syslog (Good performance, needs parsing)
0. files - Output to multiple files (Good performance, needs parsing)

[Dkro Pro](/products/pro/) provides you with several more processors.

All plugins accepts one configuration option: `forward` Indicated if the plugin must forward the original execution output. This allows for chaining plugins and sending output to different targets at the same time.

You can set more than one processor to a job. For example:

```
{
    "name": "job_name",
    "command": "/bin/true",
    "schedule": "@every 2m",
    "tags": {
        "role": "web"
    },
    "processors": {
        "files": {
            "forward": true
        },
        "syslog": {
            "forward": true
        },
    }
}
```
