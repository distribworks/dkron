package main

import (
	"github.com/victorcoder/dkron/dkron"
)

type LogOutput struct{}

func (l *LogOutput) Output(execution *dkron.Execution) []byte {
	return execution.Output
}
