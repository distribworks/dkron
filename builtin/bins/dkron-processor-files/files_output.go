package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/victorcoder/dkron/plugintypes"
)

const defaultLogDir = "/var/log/dkron"

// FilesOutput plugin that saves each execution log
// in it's own file in the file system.
type FilesOutput struct {
	forward bool
	logDir  string
}

// Process method of the plugin
func (l *FilesOutput) Process(args *plugintypes.ExecutionProcessorArgs) plugintypes.Execution {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	l.parseConfig(args.Config)

	out := args.Execution.Output
	filePath := fmt.Sprintf("%s/%s.log", l.logDir, args.Execution.Key())

	log.WithField("file", filePath).Info("files: Writing file")
	if err := ioutil.WriteFile(filePath, out, 0644); err != nil {
		log.WithError(err).Error("Error writting log file")
	}

	if !l.forward {
		args.Execution.Output = []byte(filePath)
	}

	return args.Execution
}

func (l *FilesOutput) parseConfig(config plugintypes.PluginConfig) {
	forward, ok := config["forward"].(bool)
	if ok {
		l.forward = forward
		log.Infof("Forwarding set to: %t", forward)
	} else {
		l.forward = false
		log.WithField("param", "forward").Warning("Incorrect format or param not found.")
	}

	logDir, ok := config["log_dir"].(string)
	if ok {
		l.logDir = logDir
		log.Infof("Log dir set to: %s", logDir)
	} else {
		l.logDir = defaultLogDir
		log.WithField("param", "log_dir").Warning("Incorrect format or param not found.")
		if _, err := os.Stat(defaultLogDir); os.IsNotExist(err) {
			os.MkdirAll(defaultLogDir, os.ModePerm)
		}
	}
}
