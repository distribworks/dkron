package dkron

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/armon/circbuf"
	"github.com/distribworks/dkron/v2/plugin/types"
	proto "github.com/distribworks/dkron/v2/plugin/types"
	"github.com/sirupsen/logrus"
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
	client := proto.NewDkronClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.AgentRun(ctx)
	if err != nil {
		return err
	}

	if err := stream.Send(&proto.AgentRunStream{
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

	runningExecutions.Delete(execution.GetGroup())

	// Send the final execution
	if err := stream.Send(&proto.AgentRunStream{
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
	log.Debugf("agent: AgentRun reply: %v", reply)

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
