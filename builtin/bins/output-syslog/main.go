package main

import (
	"github.com/victorcoder/dkron/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Outputter: new(SyslogOutput),
	})
}
