package main

import (
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/victorcoder/dkron/dkron"
)

type FilesOutput struct {
	forward bool
}

func (l *FilesOutput) Process(args *dkron.ExecutionProcessorArgs) dkron.Execution {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	out := args.Execution.Output
	tmp := "." //os.TempDir()
	filePath := fmt.Sprintf("%s/%s.log", tmp, args.Execution.Key())

	log.WithField("file", filePath).Info("files: Writing file")
	if err := ioutil.WriteFile(filePath, out, 0644); err != nil {
		log.WithError(err).Error("Error writting log file")
	}
	l.parseConfig(args.Config)
	if !l.forward {
		args.Execution.Output = []byte(filePath)
	}

	return args.Execution
}

func (l *FilesOutput) parseConfig(config dkron.PluginConfig) {
	forward, ok := config["forward"].(bool)
	if ok {
		l.forward = forward
		log.Infof("Forwarding set to: %s", forward)
	} else {
		log.Error("Incorrect format in forward param")
	}
}
