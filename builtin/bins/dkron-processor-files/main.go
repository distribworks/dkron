package main

import (
	"github.com/distribworks/dkron/v4/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Processor: new(FilesOutput),
	})
}
