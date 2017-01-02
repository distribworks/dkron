package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-syslog"
	"github.com/victorcoder/dkron/dkron"
)

type SyslogOutput struct {
	forward bool
}

func (l *SyslogOutput) Process(args *dkron.ExecutionProcessorArgs) dkron.Execution {
	logger, err := gsyslog.NewLogger(gsyslog.LOG_INFO, "CRON", "[dkron]")
	if err != nil {
<<<<<<< HEAD:builtin/bins/dkron-processor-syslog/syslog_output.go
		log.Error("Error creating logger")
=======
		log.WithError(err).Error("Error creating logger")
>>>>>>> master:builtin/bins/dkron-processor-syslog/syslog_output.go
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
		log.Infof("Forwarding set to: %s", forward)
	} else {
		log.Error("Incorrect format in forward param")
	}
}
