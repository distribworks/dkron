package main

import (
	"github.com/victorcoder/dkron/dkron"
)

type LogOutput struct{}

func (l *LogOutput) Process(args *dkron.ExecutionProcessorArgs) dkron.Execution {
	return args.Execution
}
