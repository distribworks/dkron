package main

import (
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/victorcoder/dkron/dkron"
)

type FilesOutput struct{}

func (l *FilesOutput) Process(args *dkron.ExecutionProcessorArgs) dkron.Execution {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	out := args.Execution.Output
	tmp := "." //os.TempDir()
	filePath := fmt.Sprintf("%s/%s.log", tmp, args.Execution.Key())

	log.WithField("file", filePath).Info("files: Writing file")
	if err := ioutil.WriteFile(filePath, out, 0644); err != nil {
		log.WithError(err).Error("Error writting log file")
	}
	args.Execution.Output = []byte(filePath)

	return args.Execution
}
