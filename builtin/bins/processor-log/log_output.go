package main

import (
	"github.com/victorcoder/dkron/dkron"
)

type LogOutput struct{}

func (l *LogOutput) Process(execution *dkron.Execution) dkron.Execution {
	return *execution
}
