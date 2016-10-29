package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/victorcoder/dkron/dkron"
)

type LogOutput struct {
	forward bool
}

func (l *LogOutput) Process(args *dkron.ExecutionProcessorArgs) dkron.Execution {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	l.parseConfig(args.Config)
	if !l.forward {
		args.Execution.Output = []byte("Output in dkron log")
	}

	return args.Execution
}

func (l *LogOutput) parseConfig(config dkron.PluginConfig) {
	forward, ok := config["forward"].(bool)
	if ok {
		l.forward = forward
		log.Infof("Forwarding set to: %s", forward)
	} else {
		log.Error("Incorrect format in forward param")
	}
}
