package main

import (
	"github.com/distribworks/dkron/dkron"
	"github.com/hashicorp/go-syslog"
	log "github.com/sirupsen/logrus"
)

type SyslogOutput struct {
	forward bool
}

func (l *SyslogOutput) Process(args *dkron.ExecutionProcessorArgs) dkron.Execution {
	logger, err := gsyslog.NewLogger(gsyslog.LOG_INFO, "CRON", "[dkron]")
	if err != nil {
		log.WithError(err).Error("Error creating logger")
		return args.Execution
	}
	logger.WriteLevel(gsyslog.LOG_INFO, args.Execution.Output)

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
		log.Infof("Forwarding set to: %t", forward)
	} else {
		log.Error("Incorrect format in forward param")
	}
}
