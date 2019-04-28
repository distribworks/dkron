package dkron

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/armon/circbuf"
	"github.com/golang/groupcache/consistenthash"
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
	log.WithField("server", rpcServer).Debug("invoke: *************** Select server to send result")

	runningExecutions.Delete(execution.GetGroup())

	return a.GRPCClient.CallExecutionDone(rpcServer, execution)
}

// Select a server based on key using a consistent hash key
// like a cache store.
func (a *Agent) selectServerByKey(key string) (string, error) {
	ch := consistenthash.New(50, nil)
	a.UpdatePeers(ch.Add)
	su := ch.Get(key)
	u, err := url.Parse(su)
	if err != nil {
		return "", err
	}
	var server serf.Member
	ss := a.ListServers()

	uPort, err := strconv.Atoi(u.Port())
	if err != nil {
		return "", err
	}

	for _, s := range ss {
		if s.Addr.String() == u.Hostname() && s.Port == uint16(uPort-1000) {
			server = s
		}
	}

	if addr, ok := server.Tags["dkron_rpc_addr"]; ok {
		return addr, nil
	}
	return "", ErrNoRPCAddress
}
