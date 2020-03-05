package dkron

import (
	"errors"
	"time"

	"github.com/armon/circbuf"
	"github.com/distribworks/dkron/v2/plugin/types"
	"github.com/golang/groupcache/consistenthash"
)

const (
	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 256000
)

// ErrNoSuitableServer returns an error in case no suitable server to send the request is found.
var ErrNoSuitableServer = errors.New("no suitable server found to send the request, aborting")

// invokeJob will execute the given job. Depending on the event.
func (a *Agent) invokeJob(job *Job, execution *Execution) error {
	output, _ := circbuf.NewBuffer(maxBufSize)

	var success bool

	jex := job.Executor
	exc := job.ExecutorConfig
	if jex == "" {
		return errors.New("invoke: No executor defined, nothing to do")
	}

	// Check if executor exists
	if executor, ok := a.ExecutorPlugins[jex]; ok {
		log.WithField("plugin", jex).Debug("invoke: calling executor plugin")
		runningExecutions.Store(execution.GetGroup(), execution)
		out, err := executor.Execute(&types.ExecuteRequest{
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

	return a.GRPCClient.ExecutionDone(rpcServer, execution)
}

// Select a server based on key using a consistent hash key
// like a cache store.
func (a *Agent) selectServerByKey(key string) (string, error) {
	ch := consistenthash.New(50, nil)
	for _, p := range a.LocalServers() {
		ch.Add(p.RPCAddr.String())
	}

	peerAddress := ch.Get(key)
	if peerAddress == "" {
		return "", ErrNoSuitableServer
	}

	return peerAddress, nil
}
