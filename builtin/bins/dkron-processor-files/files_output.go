package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/distribworks/dkron/v3/plugin"
	"github.com/distribworks/dkron/v3/plugin/types"
	log "github.com/sirupsen/logrus"
)

const defaultLogDir = "/var/log/dkron"

// FilesOutput plugin that saves each execution log
// in it's own file in the file system.
type FilesOutput struct {
	forward bool
	logDir  string
}

// Process method writes the execution output to a file
func (l *FilesOutput) Process(args *plugin.ProcessorArgs) types.Execution {
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

func (l *FilesOutput) parseConfig(config plugin.Config) {
	forward, err := strconv.ParseBool(config["forward"])
	if err != nil {
		l.forward = false
		log.WithField("param", "forward").Warning("Incorrect format or param not found.")
	} else {
		l.forward = forward
		log.Infof("Forwarding set to: %t", forward)
	}

	logDir := config["log_dir"]
	if logDir != "" {
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
