package dkron

import (
	"errors"
	"time"

	"github.com/armon/circbuf"
	metrics "github.com/armon/go-metrics"
	"github.com/distribworks/dkron/v3/plugin/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
)

const (
	// maxBufSize limits how much data we collect from a handler.
	maxBufSize = 256000
)

type statusAgentHelper struct {
	execution *types.Execution
	stream    types.Agent_AgentRunServer
}

func (s *statusAgentHelper) Update(b []byte, c bool) (int64, error) {
	s.execution.Output = b
	// Send partial execution
	if err := s.stream.Send(&types.AgentRunStream{
		Execution: s.execution,
	}); err != nil {
		return 0, err
	}
	return 0, nil
}

// GRPCAgentServer is the local implementation of the gRPC server interface.
type AgentServer struct {
	agent *Agent
}

// NewServer creates and returns an instance of a DkronGRPCServer implementation
func NewAgentServer(agent *Agent) types.AgentServer {
	return &AgentServer{
		agent: agent,
	}
}

// AgentRun is called when an agent starts running a job and lasts all execution,
// the agent will stream execution progress to the server.
func (as *AgentServer) AgentRun(req *types.AgentRunRequest, stream types.Agent_AgentRunServer) error {
	defer metrics.MeasureSince([]string{"grpc_agent", "agent_run"}, time.Now())

	job := req.Job
	execution := req.Execution

	log.WithFields(logrus.Fields{
		"job": job.Name,
	}).Info("grpc_agent: Starting job")

	output, _ := circbuf.NewBuffer(maxBufSize)

	var success bool

	jex := job.Executor
	exc := job.ExecutorConfig
	if jex == "" {
		return errors.New("grpc_agent: No executor defined, nothing to do")
	}

	// Send the first update with the initial execution state to be stored in the server
	execution.StartedAt = ptypes.TimestampNow()
	execution.NodeName = as.agent.config.NodeName

	if err := stream.Send(&types.AgentRunStream{
		Execution: execution,
	}); err != nil {
		return err
	}

	// Check if executor exists
	if executor, ok := as.agent.ExecutorPlugins[jex]; ok {
		log.WithField("plugin", jex).Debug("grpc_agent: calling executor plugin")
		runningExecutions.Store(execution.GetGroup(), execution)
		out, err := executor.Execute(&types.ExecuteRequest{
			JobName: job.Name,
			Config:  exc,
		}, &statusAgentHelper{
			stream:    stream,
			execution: execution,
		})

		if err == nil && out.Error != "" {
			err = errors.New(out.Error)
		}
		if err != nil {
			log.WithError(err).WithField("job", job.Name).WithField("plugin", executor).Error("grpc_agent: command error output")
			success = false
			output.Write([]byte(err.Error() + "\n"))
		} else {
			success = true
		}

		if out != nil {
			output.Write(out.Output)
		}
	} else {
		log.WithField("executor", jex).Error("grpc_agent: Specified executor is not present")
		output.Write([]byte("grpc_agent: Specified executor is not present"))
	}

	execution.FinishedAt = ptypes.TimestampNow()
	execution.Success = success
	execution.Output = output.Bytes()

	runningExecutions.Delete(execution.GetGroup())

	// Send the final execution
	if err := stream.Send(&types.AgentRunStream{
		Execution: execution,
	}); err != nil {
		// In case of error means that maybe the server is gone so fallback to ExecutionDone
		log.WithError(err).WithField("job", job.Name).Error("grpc_agent: error sending the final execution, falling back to ExecutionDone")
		rpcServer, err := as.agent.checkAndSelectServer()
		if err != nil {
			return err
		}
		return as.agent.GRPCClient.ExecutionDone(rpcServer, NewExecutionFromProto(execution))
	}

	return nil
}
