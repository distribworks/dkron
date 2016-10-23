package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/victorcoder/dkron/dkron"
)

type LogOutput struct{}

func (l *LogOutput) Process(args *dkron.ExecutionProcessorArgs) dkron.Execution {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	log.Info(args.Execution.Output)

	return args.Execution
}
