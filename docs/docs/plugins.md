# Intro

Plugins in Dkron allow you to add funcionality that integrates with the workflow of the job execution in Dkron. It's a powerful system that allows you to extend and adapt Dkron to your special needs.

This page documents the basics of how the plugin system in Dkron works, and how to setup a basic development environment for plugin development if you're writing a Dkron plugin.

Advanced topic! Plugin development is a highly advanced topic, and is not required knowledge for day-to-day usage. If you don't plan on writing any plugins, we recommend not reading this section of the documentation.

## How it Works

Dkron execution execution processors are provided via plugins. Each plugin exposes functionality for modifying the execution. Plugins are executed as a separate process and communicate with the main Dkron binary over an RPC interface.

The code within the binaries must adhere to certain interfaces. The network communication and RPC is handled automatically by higher-level libraries. The exact interface to implement is documented in its respective documentation section.

## Installing a Plugin

Dkron searches for plugins at startup, to install a plugin just drop the binary in one of the following locations:

1. /etc/dkron/plugins
2. Dkron executable directory

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

## Using plugins

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


## Developing a Plugin

Developing a plugin is simple. The only knowledge necessary to write a plugin is basic command-line skills and basic knowledge of the Go programming language.

Note: A common pitfall is not properly setting up a $GOPATH. This can lead to strange errors. You can read more about this here to familiarize yourself.

Create a new Go project somewhere in your $GOPATH. If you're a GitHub user, we recommend creating the project in the directory $GOPATH/src/github.com/USERNAME/dkron-NAME, where USERNAME is your GitHub username and NAME is the name of the plugin you're developing. This structure is what Go expects and simplifies things down the road.

With the directory made, create a main.go file. This project will be a binary so the package is "main":

package main

import (
    "github.com/victorcoder/dkron/plugin"
)

func main() {
    plugin.Serve(new(MyPlugin))
}
And that's basically it! You'll have to change the argument given to plugin.Serve to be your actual plugin, but that is the only change you'll have to make. The argument should be a structure implementing one of the plugin interfaces (depending on what sort of plugin you're creating).

Dkron plugins must follow a very specific naming convention of dkron-TYPE-NAME. For example, dkron-processor-files, which tells Dkron that the plugin is a processor that can be referenced as "files".

