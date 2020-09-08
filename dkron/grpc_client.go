package dkron

import (
	"fmt"
	"io"
	"time"

	metrics "github.com/armon/go-metrics"
	proto "github.com/distribworks/dkron/v3/plugin/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// DkronGRPCClient defines the interface that any gRPC client for
// dkron should implement.
type DkronGRPCClient interface {
	Connect(string) (*grpc.ClientConn, error)
	ExecutionDone(string, *Execution) error
	GetJob(string, string) (*Job, error)
	SetJob(*Job) error
	DeleteJob(string) (*Job, error)
	Leave(string) error
	RunJob(string) (*Job, error)
	RaftGetConfiguration(string) (*proto.RaftGetConfigurationResponse, error)
	RaftRemovePeerByID(string, string) error
	GetActiveExecutions(string) ([]*proto.Execution, error)
	SetExecution(execution *proto.Execution) error
	AgentRun(addr string, job *proto.Job, execution *proto.Execution) error
}

// GRPCClient is the local implementation of the DkronGRPCClient interface.
type GRPCClient struct {
	dialOpt []grpc.DialOption
	agent   *Agent
}

// NewGRPCClient returns a new instance of the gRPC client.
func NewGRPCClient(dialOpt grpc.DialOption, agent *Agent) DkronGRPCClient {
	if dialOpt == nil {
		dialOpt = grpc.WithInsecure()
	}
	return &GRPCClient{
		dialOpt: []grpc.DialOption{
			dialOpt,
			grpc.WithBlock(),
		},
		agent: agent,
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
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "ExecutionDone",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}

	d := proto.NewDkronClient(conn)
	edr, err := d.ExecutionDone(context.Background(), &proto.ExecutionDoneRequest{Execution: execution.ToProto()})
	if err != nil {
		if err.Error() == fmt.Sprintf("rpc error: code = Unknown desc = %s", ErrNotLeader.Error()) {
			log.Info("grpc: ExecutionDone forwarded to the leader")
			conn.Close()
			return nil
		}

		log.WithError(err).WithFields(logrus.Fields{
			"method":      "ExecutionDone",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return err
	}
	log.WithFields(logrus.Fields{
		"method":      "ExecutionDone",
		"server_addr": addr,
		"from":        edr.From,
		"payload":     string(edr.Payload),
	}).Debug("grpc: Response from method")
	conn.Close()
	return nil
}

// GetJob calls GetJob gRPC method in the server
func (grpcc *GRPCClient) GetJob(addr, jobName string) (*Job, error) {
	defer metrics.MeasureSince([]string{"grpc", "get_job"}, time.Now())
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "GetJob",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := proto.NewDkronClient(conn)
	gjr, err := d.GetJob(context.Background(), &proto.GetJobRequest{JobName: jobName})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "GetJob",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	return NewJobFromProto(gjr.Job), nil
}

// Leave calls Leave method on the gRPC server
func (grpcc *GRPCClient) Leave(addr string) error {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "Leave",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Synchronous call
	d := proto.NewDkronClient(conn)
	_, err = d.Leave(context.Background(), &empty.Empty{})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
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
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "SetJob",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Synchronous call
	d := proto.NewDkronClient(conn)
	_, err = d.SetJob(context.Background(), &proto.SetJobRequest{
		Job: job.ToProto(),
	})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
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
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "DeleteJob",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := proto.NewDkronClient(conn)
	res, err := d.DeleteJob(context.Background(), &proto.DeleteJobRequest{
		JobName: jobName,
	})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "DeleteJob",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	job := NewJobFromProto(res.Job)

	return job, nil
}

// RunJob calls the leader passing the job name
func (grpcc *GRPCClient) RunJob(jobName string) (*Job, error) {
	var conn *grpc.ClientConn

	addr := grpcc.agent.raft.Leader()

	// Initiate a connection with the server
	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "RunJob",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := proto.NewDkronClient(conn)
	res, err := d.RunJob(context.Background(), &proto.RunJobRequest{
		JobName: jobName,
	})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "RunJob",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	job := NewJobFromProto(res.Job)

	return job, nil
}

// RaftGetConfiguration get the current raft configuration of peers
func (grpcc *GRPCClient) RaftGetConfiguration(addr string) (*proto.RaftGetConfigurationResponse, error) {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "RaftGetConfiguration",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := proto.NewDkronClient(conn)
	res, err := d.RaftGetConfiguration(context.Background(), &empty.Empty{})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
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
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "RaftRemovePeerByID",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Synchronous call
	d := proto.NewDkronClient(conn)
	_, err = d.RaftRemovePeerByID(context.Background(),
		&proto.RaftRemovePeerByIDRequest{Id: peerID},
	)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "RaftRemovePeerByID",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return err
	}

	return nil
}

// GetActiveExecutions returns the active executions of a server node
func (grpcc *GRPCClient) GetActiveExecutions(addr string) ([]*proto.Execution, error) {
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(addr)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "GetActiveExecutions",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return nil, err
	}
	defer conn.Close()

	// Synchronous call
	d := proto.NewDkronClient(conn)
	gaer, err := d.GetActiveExecutions(context.Background(), &empty.Empty{})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "GetActiveExecutions",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return nil, err
	}

	return gaer.Executions, nil
}

// SetExecution calls the leader passing the execution
func (grpcc *GRPCClient) SetExecution(execution *proto.Execution) error {
	var conn *grpc.ClientConn

	addr := grpcc.agent.raft.Leader()

	// Initiate a connection with the server
	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "SetExecution",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Synchronous call
	d := proto.NewDkronClient(conn)
	_, err = d.SetExecution(context.Background(), execution)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "SetExecution",
			"server_addr": addr,
		}).Error("grpc: Error calling gRPC method")
		return err
	}
	return nil
}

// AgentRun runs a job in the given agent
func (grpcc *GRPCClient) AgentRun(addr string, job *proto.Job, execution *proto.Execution) error {
	defer metrics.MeasureSince([]string{"grpc_client", "agent_run"}, time.Now())
	var conn *grpc.ClientConn

	// Initiate a connection with the server
	conn, err := grpcc.Connect(string(addr))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"method":      "AgentRun",
			"server_addr": addr,
		}).Error("grpc: error dialing.")
		return err
	}
	defer conn.Close()

	// Streaming call
	a := proto.NewAgentClient(conn)
	stream, err := a.AgentRun(context.Background(), &proto.AgentRunRequest{
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

		// Error receiving from stream
		if err != nil {
			// At this point the execution status will be unknown, set the FinshedAt time and an explanatory message
			execution.FinishedAt = ptypes.TimestampNow()
			execution.Output = []byte(ErrBrokenStream.Error())

			log.WithError(err).Error(ErrBrokenStream)

			addr := grpcc.agent.raft.Leader()
			if err := grpcc.ExecutionDone(string(addr), NewExecutionFromProto(execution)); err != nil {
				return err
			}
			return err
		}

		// Registers an active stream
		grpcc.agent.activeExecutions.Store(ars.Execution.Key(), ars.Execution)
		log.WithField("key", ars.Execution.Key()).Debug("grpc: received execution stream")

		execution = ars.Execution
		defer grpcc.agent.activeExecutions.Delete(execution.Key())

		// Store the received execution in the raft log and store
		if !first {
			if err := grpcc.SetExecution(ars.Execution); err != nil {
				return err
			}
			first = true
		}
	}
}
