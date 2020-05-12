package main

import (
	"strconv"

	"github.com/distribworks/dkron/v3/plugin"
	"github.com/distribworks/dkron/v3/plugin/types"
	log "github.com/sirupsen/logrus"
)

// Process sends log to Fluent
func (l *FluentOutput) Process(args *plugin.ProcessorArgs) types.Execution {

	l.parseConfig(args.Config)

	var data = map[string]interface{}{
		"host":     args.Execution.NodeName,
		"job_name": args.Execution.JobName,
		"message":  args.Execution.Output,
	}

	go l.sendLog(data)

	if !l.forward {
		args.Execution.Output = []byte("Output sent to Fluent")
	}

	return args.Execution
}

func (l *FluentOutput) parseConfig(config plugin.Config) {
	forward, err := strconv.ParseBool(config["forward"])
	if err != nil {
		l.forward = false
		log.WithField("param", "forward").Warning("Incorrect format or param not found.")
	} else {
		l.forward = forward
		log.Infof("Forwarding set to: %t", forward)
	}
}

func (l *FluentOutput) sendLog(data map[string]interface{}) {
	err := l.fluent.Post(l.tag, data)
	if err != nil {
		log.WithError(err).Error("Error sending to Fluent")
	}
}
