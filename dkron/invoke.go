package dkron

import (
	"errors"
	"math/rand"
	"time"

	"github.com/armon/circbuf"
	"github.com/hashicorp/serf/serf"
)

const (
	windows = "windows"

	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 256000
)

// invokeJob will execute the given job. Depending on the event.
func (a *Agent) invokeJob(job *Job, execution *Execution) error {
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
		out, err := executor.Execute(&ExecuteRequest{
			JobName: job.Name,
			Config:  exc,
		})
		if err != nil {
			log.WithError(err).WithField("job", job.Name).WithField("plugin", executor).Error("invoke: command error output")
			success = false
			output.Write([]byte(err.Error()))
		} else {
			success = true
		}

		output.Write(out)
	} else {
		log.WithField("executor", jex).Error("invoke: Specified executor is not present")
	}

	execution.FinishedAt = time.Now()
	execution.Success = success
	execution.Output = output.Bytes()

	rpcServer, err := a.getServerRPCAddresFromTags()
	if err != nil {
		return err
	}

	runningExecutions.Delete(execution.GetGroup())

	return a.GRPCClient.CallExecutionDone(rpcServer, execution)
}

func (a *Agent) selectServer() serf.Member {
	servers := a.listServers()
	server := servers[rand.Intn(len(servers))]

	return server
}

func (a *Agent) getServerRPCAddresFromTags() (string, error) {
	s := a.selectServer()

	if addr, ok := s.Tags["dkron_rpc_addr"]; ok {
		return addr, nil
	}
	return "", ErrNoRPCAddress
}
