---
title: Execution Processors
weight: 2
---

## Execution Processors

Processor plugins are called when an execution response has been received. They are passed the resulting execution data and configuration parameters, this plugins can perform a variety of operations with the execution and it's very flexible and per Job, examples of operations this plugins can do:

* Execution output storage, forwarding or redirection.
* Notification
* Monitoring

Currently Dkron provides you with some stock plugins but the list keeps growing. Some of the features previously implemented in the application will be progessively moved to plugins.

## Logging processors

Logging output of each job execution can be modified by using processor plugins.

Processor plugins can be used to redirect the output of a job execution to different targets.

Processors are set per job using the `processors` property. `processors` is an object of processor plugins to use and it's corresponding configuration. To know what parameters each plugin accepts refer to the plugin documentation.

### Built-in logging processors

Depending on your needs the execution log can be redirected using the following plugins:

0. not specified - Store the output in the key value store (Slow performance, good for testing, default method)
0. log - Output the execution log to Dkron stdout (Good performance, needs parsing)
0. syslog - Output to the syslog (Good performance, needs parsing)
0. files - Output to multiple files (Good performance, needs parsing)

All plugins accepts one configuration option: `forward` Indicated if the plugin must forward the original execution output. This allows for chaining plugins and sending output to different targets at the same time.

## Using processors

For each job you can configure an arbitrary number of plugins.

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
        }
    }
}
```
