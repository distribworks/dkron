package main

import (
	"log"
	"log/syslog"

	"github.com/victorcoder/dkron/dkron"
)

type SyslogOutput struct{}

func (l *SyslogOutput) Process(args *dkron.ExecutionProcessorArgs) dkron.Execution {
	logwriter, err := syslog.New(syslog.LOG_INFO, "dkron")
	if err == nil {
		log.SetOutput(logwriter)
	}
	args.Execution.Output = []byte("Output in syslog")
	return args.Execution
}
