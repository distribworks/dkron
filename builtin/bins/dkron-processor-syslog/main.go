package main

import (
	"github.com/distribworks/dkron/v3/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Processor: new(SyslogOutput),
	})
}
