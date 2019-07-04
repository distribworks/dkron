package main

import (
	"github.com/distribworks/dkron/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Processor: new(SyslogOutput),
	})
}
