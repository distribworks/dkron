package dkron

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorcoder/dkron/dkronpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"net"
	"time"
)

func Test_buildCmd(t *testing.T) {
	testJob1 := &Job{
		Command: "echo 'test1' && echo 'success'",
		Shell:   true,
	}

	cmd := buildCmd(testJob1)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, "test1\nsuccess\n", string(out))

	testJob2 := &Job{
		Command: "date && echo 'success'",
		Shell:   false,
	}
	cmd = buildCmd(testJob2)
	out, err = cmd.CombinedOutput()
	assert.Error(t, err)
}

type testServer struct{}

func (ts *testServer) Invoke(ctx context.Context, in *dkronpb.Execution) (*dkronpb.ExecutionResult, error) {
	switch in.JobName {
	case "success":
		return &dkronpb.ExecutionResult{
			Output: []byte(in.Payload),
		}, nil
	case "timeout":
		time.Sleep(time.Second * 2)
		return &dkronpb.ExecutionResult{
			Output: []byte(in.Payload),
		}, nil
	default:
		return nil, grpc.Errorf(codes.InvalidArgument, "unknown job")
	}
}

func Test_grpc(t *testing.T) {
	ts := &testServer{}
	lis, err := net.Listen("tcp", ":9001")
	assert.NoError(t, err)
	grpcServer := grpc.NewServer()
	dkronpb.RegisterDkronExecutorServer(grpcServer, ts)
	defer lis.Close()
	go grpcServer.Serve(lis)
	serverURL := "localhost:9001"
	tests := []struct {
		job *Job
		err error
	}{
		{
			job: &Job{
				Name: "success",
				Type: GrpcJob,
				Grpc: &GrpcCommand{URL: serverURL, Payload: "OK"},
			},
		},
		{
			job: &Job{
				Name: "nojob",
				Type: GrpcJob,
				Grpc: &GrpcCommand{URL: serverURL},
			},
			err: grpc.Errorf(codes.InvalidArgument, "unknown job"),
		},
		{
			job: &Job{
				Name: "timeout",
				Type: GrpcJob,
				Grpc: &GrpcCommand{URL: serverURL, Timeout: 1},
			},
			err: grpc.Errorf(codes.DeadlineExceeded, "context deadline exceeded"),
		},
	}
	for i, test := range tests {
		res, err := invokeGrpcJob(test.job)
		assert.Equal(t, err, test.err, "case %d", i)
		if err != nil {
			continue
		}
		assert.Equal(t, test.job.Grpc.Payload, string(res.Output), "case %d", i)
	}
}
