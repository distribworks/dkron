package dkron

import (
	typesv1 "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/distribworks/dkron/v4/plugin"
	"github.com/sirupsen/logrus"
)

// processorStream feeds the output of a running execution to the job's processors
// as it arrives, so they can handle the complete output and not only the capped
// output stored with the finished execution.
type processorStream struct {
	jobName string
	key     string
	targets map[string]processorTarget
	logger  *logrus.Entry
}

type processorTarget struct {
	processor plugin.Processor
	config    plugin.Config
}

// newProcessorStream calls Begin on every processor of the job. Processors that
// fail to begin are left out of the stream, they still get Process when the
// execution is done.
func newProcessorStream(processors map[string]plugin.Processor, job *typesv1.Job, execution *typesv1.Execution, nodeName string, logger *logrus.Entry) *processorStream {
	ps := &processorStream{
		jobName: job.Name,
		key:     execution.Key(),
		targets: make(map[string]processorTarget),
		logger:  logger,
	}
	for name, pc := range job.Processors {
		processor, ok := processors[name]
		if !ok {
			continue
		}
		config := plugin.Config{"reporting_node": nodeName}
		for k, v := range pc.GetConfig() {
			config[k] = v
		}
		if err := processor.Begin(&plugin.ProcessorArgs{Execution: execution, Config: config}); err != nil {
			ps.logger.WithError(err).WithField("plugin", name).Error("processor: Error beginning output stream")
			continue
		}
		ps.targets[name] = processorTarget{processor: processor, config: config}
	}
	return ps
}

// write passes a chunk of output to every processor in the stream. A processor
// that fails to write is dropped from the stream, so it never sees output with
// holes in it.
func (ps *processorStream) write(chunk []byte, isError bool) {
	if len(chunk) == 0 {
		return
	}
	for name, t := range ps.targets {
		err := t.processor.Write(&plugin.ProcessorOutputArgs{
			JobName:      ps.jobName,
			ExecutionKey: ps.key,
			Output:       chunk,
			Error:        isError,
			Config:       t.config,
		})
		if err != nil {
			ps.logger.WithError(err).WithField("plugin", name).Error("processor: Error writing output, stopping output stream")
			delete(ps.targets, name)
		}
	}
}
