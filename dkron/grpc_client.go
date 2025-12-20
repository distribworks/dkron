package dkron

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	typesv1 "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DkronGRPCClient defines the interface that any gRPC client for
// dkron should implement.
type DkronGRPCClient interface {
	Connect(string) (*grpc.ClientConn, error)
	ExecutionDone(string, *Execution) error
	GetJob(string, string) (*Job, error)
	SetJob(*Job) error
	DeleteJob(string) (*Job, error)
	DeleteExecutions(string) (*Job, error)
	Leave(string) error
	RunJob(string) (*Job, error)
	RaftGetConfiguration(string) (*typesv1.RaftGetConfigurationResponse, error)
	RaftRemovePeerByID(string, string) error
	GetActiveExecutions(string) ([]*typesv1.Execution, error)
	SetExecution(execution *typesv1.Execution) error
	AgentRun(addr string, job *typesv1.Job, execution *typesv1.Execution) error
}

// GRPCClient is the local implementation of the DkronGRPCClient interface.
type GRPCClient struct {
	dialOpt []grpc.DialOption
	agent   *Agent
	logger  *logrus.Entry
}

// NewGRPCClient returns a new instance of the gRPC client.
func NewGRPCClient(dialOpt grpc.DialOption, agent *Agent, logger *logrus.Entry) DkronGRPCClient {
	if dialOpt == nil {
		dialOpt = grpc.WithInsecure()
	}
	return &GRPCClient{
		dialOpt: []grpc.DialOption{
			dialOpt,
			grpc.WithBlock(),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // Add tracing to gRPC client
		},
		agent:  agent,
		logger: logger,
	}
}

// Connect dialing to a gRPC server
func (grpcc *GRPCClient) Connect(addr string) (*grpc.ClientConn, error) {
	// Initiate a connection with the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpcc.dialOpt...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// ExecutionDone calls the ExecutionDone gRPC method
func (grpcc *GRPCClient) ExecutionDone(addr string, execution *Execution) error {
	defer metrics.MeasureSince([]string{"grpc", "call_execution_done"}, time.Now())
	var conn *grpc.ClientConn

	conn, err := grpcc.Connect(addr)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "ExecutionDone",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	d := typesv1.NewDkronClient(conn)
	edr, err := d.ExecutionDone(context.Background(), &typesv1.ExecutionDoneRequest{Execution: execution.ToProto()})
	if err != nil {
		if err.Error() == fmt.Sprintf("rpc error: code = Unknown desc = %s", ErrNotLeader.Error()) {
			grpcc.logger.Info("grpc: ExecutionDone forwarded to the leader")
			return nil
		}

		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "ExecutionDone",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return err
	}
	grpcc.logger.WithFields(logrus.Fields{
		"method":      "ExecutionDone",
		"server_addr": addr,
		"from":        edr.From,
		"payload":     string(edr.Payload),
	}).Debug("grpc: Response from method")
	return nil
}

// GetJob calls GetJob gRPC method in the server
func (grpcc *GRPCClient) GetJob(addr, jobName string) (*Job, error) {
	defer metrics.MeasureSince([]string{"grpc", "get_job"}, time.Now())
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "GetJob",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := typesv1.NewDkronClient(conn)
	gjr, err := d.GetJob(context.Background(), &typesv1.GetJobRequest{JobName: jobName})
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "GetJob",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	return NewJobFromProto(gjr.Job, grpcc.logger), nil
}

// Leave calls Leave method on the gRPC server
func (grpcc *GRPCClient) Leave(addr string) error {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "Leave",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Synchronous call
	d := typesv1.NewDkronClient(conn)
	_, err = d.Leave(context.Background(), &emptypb.Empty{})
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "Leave",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return err
	}

	return nil
}

// SetJob calls the leader passing the job
func (grpcc *GRPCClient) SetJob(job *Job) error {
	var conn *grpc.ClientConn

	addr := grpcc.agent.raft.Leader()

	// Initiate a connection with the server
	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "SetJob",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Synchronous call
	d := typesv1.NewDkronClient(conn)
	_, err = d.SetJob(context.Background(), &typesv1.SetJobRequest{
		Job: job.ToProto(),
	})
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "SetJob",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return err
	}
	return nil
}

// DeleteJob calls the leader passing the job name
func (grpcc *GRPCClient) DeleteJob(jobName string) (*Job, error) {
	var conn *grpc.ClientConn

	addr := grpcc.agent.raft.Leader()

	// Initiate a connection with the server
	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "DeleteJob",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := typesv1.NewDkronClient(conn)
	res, err := d.DeleteJob(context.Background(), &typesv1.DeleteJobRequest{
		JobName: jobName,
	})
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "DeleteJob",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	job := NewJobFromProto(res.Job, grpcc.logger)

	return job, nil
}

// DeleteExecutions calls the leader to delete all executions for a job and reset counters
func (grpcc *GRPCClient) DeleteExecutions(jobName string) (*Job, error) {
	if jobName == "" {
		return nil, fmt.Errorf("job name cannot be empty")
	}

	var conn *grpc.ClientConn

	addr := grpcc.agent.raft.Leader()

	// Initiate a connection with the server
	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "DeleteExecutions",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := typesv1.NewDkronClient(conn)
	res, err := d.DeleteExecutions(context.Background(), &typesv1.DeleteExecutionsRequest{
		JobName: jobName,
	})
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "DeleteExecutions",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	job := NewJobFromProto(res.Job, grpcc.logger)

	return job, nil
}

// RunJob calls the leader passing the job name
func (grpcc *GRPCClient) RunJob(jobName string) (*Job, error) {
	var conn *grpc.ClientConn

	addr := grpcc.agent.raft.Leader()

	// Initiate a connection with the server
	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "RunJob",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := typesv1.NewDkronClient(conn)
	res, err := d.RunJob(context.Background(), &typesv1.RunJobRequest{
		JobName: jobName,
	})
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "RunJob",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	job := NewJobFromProto(res.Job, grpcc.logger)

	return job, nil
}

// RaftGetConfiguration get the current raft configuration of peers
func (grpcc *GRPCClient) RaftGetConfiguration(addr string) (*typesv1.RaftGetConfigurationResponse, error) {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "RaftGetConfiguration",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := typesv1.NewDkronClient(conn)
	res, err := d.RaftGetConfiguration(context.Background(), &emptypb.Empty{})
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "RaftGetConfiguration",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	return res, nil
}

// RaftRemovePeerByID remove a raft peer
func (grpcc *GRPCClient) RaftRemovePeerByID(addr, peerID string) error {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "RaftRemovePeerByID",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Synchronous call
	d := typesv1.NewDkronClient(conn)
	_, err = d.RaftRemovePeerByID(context.Background(),
		&typesv1.RaftRemovePeerByIDRequest{Id: peerID},
	)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "RaftRemovePeerByID",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return err
	}

	return nil
}

// GetActiveExecutions returns the active executions of a server node
func (grpcc *GRPCClient) GetActiveExecutions(addr string) ([]*typesv1.Execution, error) {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "GetActiveExecutions",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := typesv1.NewDkronClient(conn)
	gaer, err := d.GetActiveExecutions(context.Background(), &emptypb.Empty{})
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "GetActiveExecutions",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	return gaer.Executions, nil
}

// SetExecution calls the leader passing the execution
func (grpcc *GRPCClient) SetExecution(execution *typesv1.Execution) error {
	var conn *grpc.ClientConn

	addr := grpcc.agent.raft.Leader()

	// Initiate a connection with the server
	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "SetExecution",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Synchronous call
	d := typesv1.NewDkronClient(conn)
	_, err = d.SetExecution(context.Background(), execution)
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "SetExecution",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return err
	}
	return nil
}

// AgentRun runs a job in the given agent
func (grpcc *GRPCClient) AgentRun(addr string, job *typesv1.Job, execution *typesv1.Execution) error {
	defer metrics.MeasureSince([]string{"grpc_client", "agent_run"}, time.Now())

	maxRetries := grpcc.agent.config.AgentRunMaxRetries
	initialInterval := grpcc.agent.config.AgentRunRetryInitialInterval
	maxInterval := grpcc.agent.config.AgentRunRetryMaxInterval

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Calculate exponential backoff, preventing overflow
			// For attempt 1: initialInterval * 1 (2^0)
			// For attempt 2: initialInterval * 2 (2^1)
			// For attempt 3: initialInterval * 4 (2^2), etc.
			backoff := initialInterval
			for i := 1; i < attempt; i++ {
				backoff *= 2
				// Cap early to prevent overflow
				if backoff > maxInterval {
					backoff = maxInterval
					break
				}
			}
			if backoff > maxInterval {
				backoff = maxInterval
			}
			
			grpcc.logger.WithError(lastErr).WithFields(logrus.Fields{
				"attempt":        attempt + 1,
				"total_attempts": maxRetries + 1,
				"backoff":        backoff,
				"job":            job.Name,
				"node":           addr,
			}).Warn("grpc: Retrying AgentRun after failure")
			
			time.Sleep(backoff)
		}

		err := grpcc.agentRunAttempt(addr, job, execution)
		if err == nil {
			// Success
			if attempt > 0 {
				grpcc.logger.WithFields(logrus.Fields{
					"attempt": attempt + 1,
					"job":     job.Name,
					"node":    addr,
				}).Info("grpc: AgentRun succeeded after retry")
			}
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableError(err) {
			grpcc.logger.WithError(err).WithFields(logrus.Fields{
				"job":  job.Name,
				"node": addr,
			}).Error("grpc: Non-retryable error in AgentRun")
			break
		}
	}

	// All retries exhausted
	grpcc.logger.WithError(lastErr).WithFields(logrus.Fields{
		"job":            job.Name,
		"node":           addr,
		"total_attempts": maxRetries + 1,
	}).Error("grpc: AgentRun failed after all retry attempts")

	return lastErr
}

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Try to extract gRPC status code first
	if st, ok := status.FromError(err); ok {
		code := st.Code()
		// Retry on common transient gRPC status codes
		return code == codes.Unavailable ||
			code == codes.DeadlineExceeded ||
			code == codes.ResourceExhausted ||
			code == codes.Aborted ||
			code == codes.Internal // Internal errors may be transient
	}

	// Fall back to string matching for non-gRPC errors
	errStr := err.Error()
	return strings.Contains(errStr, "transport is closing") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "context deadline exceeded")
}

// agentRunAttempt performs a single attempt of AgentRun
func (grpcc *GRPCClient) agentRunAttempt(addr string, job *typesv1.Job, execution *typesv1.Execution) error {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		grpcc.logger.WithError(err).WithFields(logrus.Fields{
			"method":      "AgentRun",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Streaming call
	a := typesv1.NewAgentServiceClient(conn)
	stream, err := a.AgentRun(context.Background(), &typesv1.AgentRunRequest{
		Job:       job,
		Execution: execution,
	})
	if err != nil {
		return err
	}

	var first bool
	for {
		ars, err := stream.Recv()

		// Stream ends
		if err == io.EOF {
			addr := grpcc.agent.raft.Leader()
			if err := grpcc.ExecutionDone(string(addr), NewExecutionFromProto(execution)); err != nil {
				return err
			}
			return nil
		}

		// Error received from the stream
		if err != nil {
			// At this point the execution status will be unknown, set the FinishedAt time and an explanatory message
			execution.FinishedAt = timestamppb.Now()
			execution.Success = false
			execution.Output = []byte(ErrBrokenStream.Error() + ": " + err.Error())

			grpcc.logger.WithError(err).Error(ErrBrokenStream)

			addr := grpcc.agent.raft.Leader()
			if err := grpcc.ExecutionDone(string(addr), NewExecutionFromProto(execution)); err != nil {
				return err
			}
			return err
		}

		// Registers an active stream
		grpcc.agent.activeExecutions.Store(ars.Execution.Key(), ars.Execution)
		grpcc.logger.WithField("key", ars.Execution.Key()).Debug("grpc: received execution stream")

		execution = ars.Execution
		defer grpcc.agent.activeExecutions.Delete(execution.Key())

		// Store the received execution in the raft log and store
		if !first {
			if err := grpcc.SetExecution(ars.Execution); err != nil {
				return err
			}
			first = true
		}

		// Notify the starting of the execution
		if err := SendPreNotifications(grpcc.agent.config, NewExecutionFromProto(execution), nil, NewJobFromProto(job, grpcc.logger), grpcc.logger); err != nil {
			grpcc.logger.WithFields(map[string]interface{}{
				"job_name": job.Name,
				"node":     grpcc.agent.config.NodeName,
			}).Error("agent: Error sending start notification")
		}
	}
}
