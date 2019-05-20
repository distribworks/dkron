package dkron

import (
	"errors"
	"time"

	"github.com/armon/circbuf"
	"github.com/golang/groupcache/consistenthash"
	"github.com/victorcoder/dkron/plugintypes"
)

const (
	windows = "windows"

	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 256000
)

// invokeJob will execute the given job. Depending on the event.
func (a *Agent) invokeJob(job *Job, execution *plugintypes.Execution) error {
	output, _ := circbuf.NewBuffer(maxBufSize)

	var success bool

	jex := job.Executor
	exc := job.ExecutorConfig
	if jex == "" {
		return errors.New("invoke: No executor defined, nothing to do")
	}

	// Check if executor is exists
	if executor, ok := a.ExecutorPlugins[jex]; ok {
		log.WithField("plugin", jex).Debug("invoke: calling executor plugin")
		runningExecutions.Store(execution.GetGroup(), execution)
		out, err := executor.Execute(&plugintypes.ExecuteRequest{
			JobName: job.Name,
			Config:  exc,
		})

		if err == nil && out.Error != "" {
			err = errors.New(out.Error)
		}
		if err != nil {
			log.WithError(err).WithField("job", job.Name).WithField("plugin", executor).Error("invoke: command error output")
			success = false
			output.Write([]byte(err.Error() + "\n"))
		} else {
			success = true
		}

		if out != nil {
			output.Write(out.Output)
		}
	} else {
		log.WithField("executor", jex).Error("invoke: Specified executor is not present")
	}

	execution.FinishedAt = time.Now()
	execution.Success = success
	execution.Output = output.Bytes()

	rpcServer, err := a.selectServerByKey(execution.Key())
	if err != nil {
		return err
	}
	log.WithField("server", rpcServer).Debug("invoke: Selected a server to send result")

	runningExecutions.Delete(execution.GetGroup())

	return a.GRPCClient.CallExecutionDone(rpcServer, execution)
}

// Select a server based on key using a consistent hash key
// like a cache store.
func (a *Agent) selectServerByKey(key string) (string, error) {
	ch := consistenthash.New(50, nil)
	ch.Add(a.GetPeers()...)
	peerAddress := ch.Get(key)

	return peerAddress, nil
}
