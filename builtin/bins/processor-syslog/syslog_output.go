package main

import (
	"log"
	"log/syslog"

	"github.com/victorcoder/dkron/dkron"
)

type SyslogOutput struct {
	forward bool
}

func (l *SyslogOutput) Process(args *dkron.ExecutionProcessorArgs) dkron.Execution {
	logwriter, err := syslog.New(syslog.LOG_INFO, "dkron")
	if err == nil {
		log.SetOutput(logwriter)
	}

	l.parseConfig(args.Config)
	if !l.forward {
		args.Execution.Output = []byte("Output in syslog")
	}

	return args.Execution
}

func (l *SyslogOutput) parseConfig(config dkron.PluginConfig) {
	forward, ok := config["forward"].(bool)
	if ok {
		l.forward = forward
		log.Infof("Forwarding set to: %s", forward)
	} else {
		log.Error("Incorrect format in forward param")
	}
}
