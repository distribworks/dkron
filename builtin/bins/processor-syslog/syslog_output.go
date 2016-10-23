package main

import (
	"log"
	"log/syslog"

	"github.com/victorcoder/dkron/dkron"
)

type SyslogOutput struct{}

func (l *SyslogOutput) Process(execution *dkron.Execution) dkron.Execution {
	logwriter, err := syslog.New(syslog.LOG_INFO, "dkron")
	if err == nil {
		log.SetOutput(logwriter)
	}

	log.Print(execution.Output)
	execution.Output = []byte("Output in syslog")
	return *execution
}
