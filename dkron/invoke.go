package dkron

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/armon/circbuf"
	"github.com/distribworks/dkron/v2/plugin/types"
	"github.com/sirupsen/logrus"
)

const (
	// maxBufSize limits how much data we collect from a handler.
	maxBufSize = 256000
)

// ErrNoSuitableServer returns an error in case no suitable server to send the request is found.
var ErrNoSuitableServer = errors.New("no suitable server found to send the request, aborting")

type statusHelper struct {
	execution *Execution
	stream    types.Dkron_AgentRunClient
}

func (s *statusHelper) Update(b []byte, c bool) (int64, error) {
	s.execution.Output = string(b)
	// Send partial execution
	if err := s.stream.Send(&types.AgentRunStream{
		Execution: s.execution.ToProto(),
	}); err != nil {
		return 0, err
	}
	return 0, nil
}

// invokeJob will execute the given job. Depending on the event.
func (a *Agent) invokeJob(job *Job, execution *Execution) error {
	output, _ := circbuf.NewBuffer(maxBufSize)

	var success bool

	jex := job.Executor
	exc := job.ExecutorConfig
	if jex == "" {
		return errors.New("invoke: No executor defined, nothing to do")
	}

	// Connect to a server to stream the execution
	rpcServer, err := a.checkAndSelectServer()
	if err != nil {
		return err
	}
	log.WithField("server", rpcServer).Debug("invoke: Selected a server to send result")

	conn, err := a.GRPCClient.Connect(rpcServer)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "Invoke",
			"server_addr": rpcServer,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()
	client := types.NewDkronClient(conn)

	stream, err := client.AgentRun(context.Background())
	if err != nil {
		return err
	}

	// Send the first update with the initial execution state to be stored in the server
	if err := stream.Send(&types.AgentRunStream{
		Execution: execution.ToProto(),
	}); err != nil {
		return err
	}

	// Check if executor exists
	if executor, ok := a.ExecutorPlugins[jex]; ok {
		log.WithField("plugin", jex).Debug("invoke: calling executor plugin")
		runningExecutions.Store(execution.GetGroup(), execution)
		out, err := executor.Execute(&types.ExecuteRequest{
			JobName: job.Name,
			Config:  exc,
		}, &statusHelper{
			stream:    stream,
			execution: execution,
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
	execution.Output = output.String()

	runningExecutions.Delete(execution.GetGroup())

	// Send the final execution
	if err := stream.Send(&types.AgentRunStream{
		Execution: execution.ToProto(),
	}); err != nil {
		// In case of error means that maybe the server is gone so fallback to ExecutionDone
		return a.GRPCClient.ExecutionDone(rpcServer, execution)
	}

	// Close the stream
	reply, err := stream.CloseAndRecv()
	if err != nil {
		// In case of error means that maybe the server is gone so fallback to ExecutionDone
		return a.GRPCClient.ExecutionDone(rpcServer, execution)
	}
	log.WithField("from", reply.From).Debug("agent: AgentRun reply")

	return nil
}

// Check if the server is alive and select it
func (a *Agent) checkAndSelectServer() (string, error) {
	var peers []string
	for _, p := range a.LocalServers() {
		peers = append(peers, p.RPCAddr.String())
	}

	for _, peer := range peers {
		log.Debugf("Checking peer: %v", peer)
		conn, err := net.DialTimeout("tcp", peer, 1*time.Second)
		if err == nil {
			conn.Close()
			log.Debugf("Found good peer: %v", peer)
			return peer, nil
		}
	}
	return "", ErrNoSuitableServer
}
