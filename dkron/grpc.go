package dkron

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"

	metrics "github.com/armon/go-metrics"
	pb "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"github.com/victorcoder/dkron/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	ErrExecutionDoneForDeletedJob = errors.New("rpc: Received execution done for a deleted job")
	ErrRPCDialing                 = errors.New("rpc: Error dialing, verify the network connection to the server")
	ErrNotLeader                  = errors.New("Error, server is not leader, this operation should be run on the leader")
)

type DkronGRPCServer interface {
	proto.DkronServer
	Serve(net.Listener) error
}

type GRPCServer struct {
	agent *Agent
}

// NewRPCServe creates and returns an instance of an RPCServer implementation
func NewGRPCServer(agent *Agent) DkronGRPCServer {
	return &GRPCServer{
		agent: agent,
	}
}

func (grpcs *GRPCServer) Serve(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	proto.RegisterDkronServer(grpcServer, grpcs)
	go grpcServer.Serve(lis)

	return nil
}

// Encode is used to encode a Protoc object with type prefix
func Encode(t MessageType, msg interface{}) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(uint8(t))
	m, err := pb.Marshal(msg.(pb.Message))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(m)
	return buf.Bytes(), err
}

// SetJob broadcast a state change to the cluster members that will store the job.
// Then restart the scheduler
// This only works on the leader
func (grpcs *GRPCServer) SetJob(ctx context.Context, setJobReq *proto.SetJobRequest) (*proto.SetJobResponse, error) {
	defer metrics.MeasureSince([]string{"grpc", "set_job"}, time.Now())
	log.WithFields(logrus.Fields{
		"job": setJobReq.Job.Name,
	}).Debug("grpc: Received SetJob")

	cmd, err := Encode(SetJobType, setJobReq.Job)
	if err != nil {
		return nil, err
	}
	af := grpcs.agent.raft.Apply(cmd, raftTimeout)
	if err := af.Error(); err != nil {
		return nil, err
	}
	res := af.Response()
	switch res {
	case ErrParentJobNotFound:
		return nil, ErrParentJobNotFound
	case ErrSameParent:
		return nil, ErrParentJobNotFound
	}

	// If everything is ok, restart the scheduler
	grpcs.agent.SchedulerRestart()

	return &proto.SetJobResponse{}, nil
}

// DeleteJob broadcast a state change to the cluster members that will delete the job.
// Then restart the scheduler
// This only works on the leader
func (grpcs *GRPCServer) DeleteJob(ctx context.Context, delJobReq *proto.DeleteJobRequest) (*proto.DeleteJobResponse, error) {
	defer metrics.MeasureSince([]string{"grpc", "delete_job"}, time.Now())
	log.WithField("job", delJobReq.GetJobName()).Debug("grpc: Received DeleteJob")

	cmd, err := Encode(DeleteJobType, delJobReq)
	if err != nil {
		return nil, err
	}
	af := grpcs.agent.raft.Apply(cmd, raftTimeout)
	if err := af.Error(); err != nil {
		return nil, err
	}
	res := af.Response()
	job, ok := res.(*Job)
	if !ok {
		return nil, fmt.Errorf("grpc: Error wrong response from apply in DeleteJob: %v", res)
	}
	jpb := job.ToProto()

	return &proto.DeleteJobResponse{Job: jpb}, nil
}

// GetJob loads the job from the datastore
func (grpcs *GRPCServer) GetJob(ctx context.Context, getJobReq *proto.GetJobRequest) (*proto.GetJobResponse, error) {
	defer metrics.MeasureSince([]string{"grpc", "get_job"}, time.Now())
	log.WithField("job", getJobReq.JobName).Debug("grpc: Received GetJob")

	j, err := grpcs.agent.Store.GetJob(getJobReq.JobName, nil)
	if err != nil {
		return nil, err
	}

	gjr := &proto.GetJobResponse{
		Job: &proto.Job{},
	}

	// Copy the data structure
	gjr.Job.Name = j.Name
	gjr.Job.Executor = j.Executor
	gjr.Job.ExecutorConfig = j.ExecutorConfig

	return gjr, nil
}

// ExecutionDone saves the execution to the store
func (grpcs *GRPCServer) ExecutionDone(ctx context.Context, execDoneReq *proto.ExecutionDoneRequest) (*proto.ExecutionDoneResponse, error) {
	defer metrics.MeasureSince([]string{"grpc", "execution_done"}, time.Now())
	log.WithFields(logrus.Fields{
		"group": execDoneReq.Execution.Group,
		"job":   execDoneReq.Execution.JobName,
		"from":  execDoneReq.Execution.NodeName,
	}).Debug("grpc: Received execution done")

	// Get the leader address and compare with the current node address.
	// Forward the request to the leader in case current node is not the leader.
	if !grpcs.agent.IsLeader() {
		addr := grpcs.agent.raft.Leader()
		grpcs.agent.GRPCClient.CallExecutionDone(string(addr), NewExecutionFromProto(execDoneReq.Execution))
		return nil, ErrNotLeader
	}

	// This is the leader at this point, so process the execution, encode the value and apply the log to the cluster.
	// Get the defined output types for the job, and call them
	job, err := grpcs.agent.Store.GetJob(execDoneReq.Execution.JobName, nil)
	if err != nil {
		return nil, err
	}
	origExec := *NewExecutionFromProto(execDoneReq.Execution)
	execution := origExec
	for k, v := range job.Processors {
		log.WithField("plugin", k).Info("grpc: Processing execution with plugin")
		if processor, ok := grpcs.agent.ProcessorPlugins[k]; ok {
			v["reporting_node"] = grpcs.agent.config.NodeName
			e := processor.Process(&ExecutionProcessorArgs{Execution: origExec, Config: v})
			execution = e
		} else {
			log.WithField("plugin", k).Error("grpc: Specified plugin not found")
		}
	}

	execDoneReq.Execution = execution.ToProto()
	cmd, err := Encode(ExecutionDoneType, execDoneReq)
	if err != nil {
		return nil, err
	}
	af := grpcs.agent.raft.Apply(cmd, raftTimeout)
	if err := af.Error(); err != nil {
		return nil, err
	}

	if err := af.Error(); err != nil {
		return nil, err
	}

	// If the execution failed, retry it until retries limit (default: don't retry)
	if !execution.Success && execution.Attempt < job.Retries+1 {
		execution.Attempt++

		// Keep all execution properties intact except the last output
		// as it could exceed serf query limits.
		execution.Output = []byte{}

		log.WithFields(logrus.Fields{
			"attempt":   execution.Attempt,
			"execution": execution,
		}).Debug("grpc: Retrying execution")

		grpcs.agent.RunQuery(job, &execution)
		return nil, nil
	}

	exg, err := grpcs.agent.Store.GetExecutionGroup(&execution)
	if err != nil {
		log.WithError(err).WithField("group", execution.Group).Error("grpc: Error getting execution group.")
		return nil, err
	}

	// Send notification
	Notification(grpcs.agent.config, &execution, exg, job).Send()

	// Jobs that have dependent jobs are a bit more expensive because we need to call the Status() method for every execution.
	// Check first if there's dependent jobs and then check for the job status to begin execution dependent jobs on success.
	if len(job.DependentJobs) > 0 && job.GetStatus() == StatusSuccess {
		for _, djn := range job.DependentJobs {
			dj, err := grpcs.agent.Store.GetJob(djn, nil)
			if err != nil {
				return nil, err
			}
			log.WithField("job", djn).Debug("grpc: Running dependent job")
			dj.Run()
		}
	}

	return &proto.ExecutionDoneResponse{
		From:    grpcs.agent.config.NodeName,
		Payload: []byte("saved"),
	}, nil
}

func (grpcs *GRPCServer) Leave(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	return in, grpcs.agent.Stop()
}

// RunJob runs a job in the cluster
func (grpcs *GRPCServer) RunJob(ctx context.Context, req *proto.RunJobRequest) (*proto.RunJobResponse, error) {
	job, err := grpcs.agent.Store.GetJob(req.JobName, nil)
	if err != nil {
		return nil, err
	}

	ex := NewExecution(job.Name)
	grpcs.agent.RunQuery(job, ex)

	jpb := job.ToProto()

	return &proto.RunJobResponse{Job: jpb}, nil
}

// ToggleJob toggle the enablement of a job
func (grpcs *GRPCServer) ToggleJob(ctx context.Context, getJobReq *proto.ToggleJobRequest) (*proto.ToggleJobResponse, error) {
	return nil, nil
}
