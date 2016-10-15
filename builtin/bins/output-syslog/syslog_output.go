package main

import (
	"log"
	"log/syslog"

	"github.com/victorcoder/dkron/dkron"
)

type SyslogOutput struct{}

func (l *SyslogOutput) Output(execution *dkron.Execution) []byte {
	logwriter, err := syslog.New(syslog.LOG_INFO, "dkron")
	if err == nil {
		log.SetOutput(logwriter)
	}

	log.Print(execution.Output)
	return []byte("Output in syslog")
}
