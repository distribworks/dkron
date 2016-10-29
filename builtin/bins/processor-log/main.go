package main

import (
	"github.com/victorcoder/dkron/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Processor: new(LogOutput),
	})
}
