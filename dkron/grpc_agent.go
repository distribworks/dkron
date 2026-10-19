package dkron

import (
	"errors"
	"time"

	"github.com/armon/circbuf"
	typesv1 "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/hashicorp/go-metrics"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// maxBufSize limits how much data we collect from a handler.
	maxBufSize = 256000
)

type statusAgentHelper struct {
	execution *typesv1.Execution
	stream    typesv1.AgentService_AgentRunServer
}

func (s *statusAgentHelper) Update(b []byte, c bool) (int64, error) {
	// Send partial execution, the chunk travels on its own so it is not sent twice.
	if err := s.stream.Send(&typesv1.AgentRunStream{
		Execution:        s.execution,
		OutputChunk:      b,
		OutputChunkError: c,
	}); err != nil {
		return 0, err
	}
	return 0, nil
}

// GRPCAgentServer is the local implementation of the gRPC server interface.
type AgentServer struct {
	typesv1.AgentServiceServer
	agent  *Agent
	logger *logrus.Entry
}

// NewServer creates and returns an instance of a DkronGRPCServer implementation
func NewAgentServer(agent *Agent, logger *logrus.Entry) typesv1.AgentServiceServer {
	return &AgentServer{
		agent:  agent,
		logger: logger,
	}
}

// AgentRun is called when an agent starts running a job and lasts all execution,
// the agent will stream execution progress to the server.
func (as *AgentServer) AgentRun(req *typesv1.AgentRunRequest, stream typesv1.AgentService_AgentRunServer) error {
	defer metrics.MeasureSince([]string{"grpc_agent", "agent_run"}, time.Now())

	job := req.Job
	execution := req.Execution

	as.logger.WithFields(logrus.Fields{
		"job": job.Name,
	}).Info("grpc_agent: Starting job")

	output, _ := circbuf.NewBuffer(maxBufSize)

	var success bool

	jex := job.Executor
	exc := job.ExecutorConfig

	// Send the first update with the initial execution state to be stored in the server
	execution.StartedAt = timestamppb.Now()
	execution.NodeName = as.agent.config.NodeName

	if err := stream.Send(&typesv1.AgentRunStream{
		Execution: execution,
	}); err != nil {
		return err
	}

	if jex == "" {
		return errors.New("grpc_agent: No executor defined, nothing to do")
	}

	// Check if executor exists
	if executor, ok := as.agent.ExecutorPlugins[jex]; ok {
		as.logger.WithField("plugin", jex).Debug("grpc_agent: calling executor plugin")
		runningExecutions.Store(execution.GetGroup(), execution)
		defer runningExecutions.Delete(execution.GetGroup())
		status := &statusAgentHelper{
			stream:    stream,
			execution: execution,
		}
		out, err := executor.Execute(&typesv1.ExecuteRequest{
			JobName: job.Name,
			Config:  exc,
		}, status)

		if err == nil && out.Error != "" {
			err = errors.New(out.Error)
		}
		if err != nil {
			as.logger.WithError(err).WithField("job", job.Name).WithField("plugin", executor).Error("grpc_agent: command error output")
			success = false
			_, _ = output.Write([]byte(err.Error() + "\n"))
			// Stream it as well, so streaming processors see the complete output.
			_, _ = status.Update([]byte(err.Error()+"\n"), true)
		} else {
			success = true
		}

		if out != nil {
			_, _ = output.Write(out.Output)
		}
	} else {
		as.logger.WithField("executor", jex).Error("grpc_agent: Specified executor is not present")
		_, _ = output.Write([]byte("grpc_agent: Specified executor is not present"))
	}

	execution.FinishedAt = timestamppb.Now()
	execution.Success = success
	execution.Output = output.Bytes()

	// Send the final execution
	if err := stream.Send(&typesv1.AgentRunStream{
		Execution: execution,
	}); err != nil {
		// In case of error means that maybe the server is gone so fallback to ExecutionDone
		as.logger.WithError(err).WithField("job", job.Name).Error("grpc_agent: error sending the final execution, falling back to ExecutionDone")
		rpcServer, err := as.agent.checkAndSelectServer()
		if err != nil {
			return err
		}
		return as.agent.GRPCClient.ExecutionDone(rpcServer, NewExecutionFromProto(execution))
	}

	return nil
}
