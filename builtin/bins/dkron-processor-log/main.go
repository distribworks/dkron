package main

import (
	"github.com/distribworks/dkron/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Processor: new(LogOutput),
	})
}
