package main

import (
	"strconv"

	"github.com/distribworks/dkron/v3/plugin"
	"github.com/distribworks/dkron/v3/plugin/types"
	gsyslog "github.com/hashicorp/go-syslog"
	log "github.com/sirupsen/logrus"
)

type SyslogOutput struct {
	forward bool
}

func (l *SyslogOutput) Process(args *plugin.ProcessorArgs) types.Execution {
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

func (l *SyslogOutput) parseConfig(config plugin.Config) {
	forward, err := strconv.ParseBool(config["forward"])
	if err != nil {
		l.forward = false
		log.WithField("param", "forward").Warning("Incorrect format or param not found.")
	} else {
		l.forward = forward
		log.Infof("Forwarding set to: %t", forward)
	}
}
