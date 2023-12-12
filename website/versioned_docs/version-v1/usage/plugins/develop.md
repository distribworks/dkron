---
title: Developing plugins
weight: 99
---

## Developing a Plugin

{{% alert theme="warning" %}}**Advanced topic!** Plugin development is a highly advanced topic, and is not required knowledge for day-to-day usage. If you don't plan on writing any plugins, we recommend not reading the following section of the documentation.{{% /alert %}}

Developing a plugin is simple. The only knowledge necessary to write a plugin is basic command-line skills and basic knowledge of the Go programming language.

Note: A common pitfall is not properly setting up a `$GOPATH`. This can lead to strange errors. You can read more about this here to familiarize yourself.

Create a new Go project somewhere in your `$GOPATH`. If you're a GitHub user, we recommend creating the project in the directory `$GOPATH/src/github.com/USERNAME/dkron-NAME-TYPE`, where `USERNAME` is your GitHub username and `NAME` is the name of the plugin you're developing. This structure is what Go expects and simplifies things down the road.

With the directory made, create a main.go file. This project will be a binary so the package is "main":

```go
package main

import (
	"github.com/distribworks/dkron/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Processor: new(MyPlugin),
	})
}
```

And that's basically it! You'll have to change the argument given to plugin.Serve to be your actual plugin, but that is the only change you'll have to make. The argument should be a structure implementing one of the plugin interfaces (depending on what sort of plugin you're creating).

Dkron plugins must follow a very specific naming convention of `dkron-TYPE-NAME`. For example, `dkron-processor-files`, which tells Dkron that the plugin is a processor that can be referenced as "files".
