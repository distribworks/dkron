package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	types "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/distribworks/dkron/v4/plugin"
	log "github.com/sirupsen/logrus"
)

const defaultLogDir = "/var/log/dkron"

// FilesOutput plugin that saves each execution log
// in it's own file in the file system.
//
// The output is written to the file while the execution runs, so the file
// holds the complete output and can be followed live. If the output was not
// streamed to this node, the output of the finished execution is written instead.
type FilesOutput struct {
	forward bool
	logDir  string

	mu sync.Mutex
	// open holds the log files of running executions, by execution key.
	open map[string]*os.File
}

// Begin creates the log file of a starting execution
func (l *FilesOutput) Begin(args *plugin.ProcessorArgs) error {
	logDir := l.parseLogDir(args.Config)
	key := args.Execution.Key()
	filePath := fmt.Sprintf("%s/%s.log", logDir, key)

	log.WithField("file", filePath).Info("files: Streaming output to file")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.open == nil {
		l.open = make(map[string]*os.File)
	}
	l.open[key] = f
	return nil
}

// Write appends a chunk of output to the log file of a running execution
func (l *FilesOutput) Write(args *plugin.ProcessorOutputArgs) error {
	l.mu.Lock()
	f, ok := l.open[args.ExecutionKey]
	l.mu.Unlock()
	if !ok {
		return fmt.Errorf("files: No log file open for execution %s", args.ExecutionKey)
	}
	_, err := f.Write(args.Output)
	return err
}

// Process method writes the execution output to a file
func (l *FilesOutput) Process(args *plugin.ProcessorArgs) *types.Execution {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	l.parseConfig(args.Config)

	key := args.Execution.Key()
	filePath := fmt.Sprintf("%s/%s.log", l.logDir, key)

	l.mu.Lock()
	f, streamed := l.open[key]
	delete(l.open, key)
	l.mu.Unlock()

	if streamed {
		// The file already holds the complete output.
		log.WithField("file", filePath).Info("files: Closing streamed file")
		if err := f.Close(); err != nil {
			log.WithError(err).Error("Error closing log file")
		}
	} else {
		log.WithField("file", filePath).Info("files: Writing file")
		if err := os.WriteFile(filePath, args.Execution.Output, 0644); err != nil {
			log.WithError(err).Error("Error writing log file")
		}
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

	l.logDir = l.parseLogDir(config)
}

func (l *FilesOutput) parseLogDir(config plugin.Config) string {
	logDir := config["log_dir"]
	if logDir != "" {
		log.Infof("Log dir set to: %s", logDir)
		return logDir
	}

	log.WithField("param", "log_dir").Warning("Incorrect format or param not found.")
	if _, err := os.Stat(defaultLogDir); os.IsNotExist(err) {
		os.MkdirAll(defaultLogDir, os.ModePerm)
	}
	return defaultLogDir
}
