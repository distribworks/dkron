package dkron

import (
	"testing"
	"time"

	"github.com/distribworks/dkron/v2/plugin"
	proto "github.com/distribworks/dkron/v2/plugin/types"
	"github.com/hashicorp/serf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestJobGetParent(t *testing.T) {
	s, err := NewStore()
	defer s.Shutdown()
	require.NoError(t, err)

	parentTestJob := &Job{
		Name:           "parent_test",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/false"},
		Schedule:       "@every 2s",
	}

	if err := s.SetJob(parentTestJob, true); err != nil {
		t.Fatalf("error creating job: %s", err)
	}

	dependentTestJob := &Job{
		Name:           "dependent_test",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/false"},
		ParentJob:      "parent_test",
	}

	err = s.SetJob(dependentTestJob, true)
	assert.NoError(t, err)

	parentTestJob, err = dependentTestJob.GetParent(s)
	assert.NoError(t, err)
	assert.Equal(t, []string{dependentTestJob.Name}, parentTestJob.DependentJobs)

	ptj, err := dependentTestJob.GetParent(s)
	assert.NoError(t, err)
	assert.Equal(t, parentTestJob.Name, ptj.Name)

	// Remove the parent job
	dependentTestJob.ParentJob = ""
	dependentTestJob.Schedule = "@every 2m"
	err = s.SetJob(dependentTestJob, true)
	assert.NoError(t, err)

	dtj, _ := s.GetJob(dependentTestJob.Name, nil)
	assert.NoError(t, err)
	assert.Equal(t, "", dtj.ParentJob)

	ptj, err = dtj.GetParent(s)
	assert.EqualError(t, ErrNoParent, err.Error())

	ptj, err = s.GetJob(parentTestJob.Name, nil)
	assert.NoError(t, err)
	assert.Nil(t, ptj.DependentJobs)
}

func TestNewJobFromProto(t *testing.T) {
	testConfig := map[string]plugin.Config{
		"test_processor": {
			"config_key": "config_value",
		},
	}

	in := &proto.Job{}
	proc := map[string]*proto.PluginConfig{
		"test_processor": {
			Config: map[string]string{"config_key": "config_value"},
		},
	}
	in.Processors = proc

	j := NewJobFromProto(in)
	assert.Equal(t, testConfig, j.Processors)
}

func TestToProto(t *testing.T) {
	j := &Job{
		Processors: map[string]plugin.Config{
			"test_processor": {
				"config_key": "config_value",
			},
		},
	}
	proc := map[string]*proto.PluginConfig{
		"test_processor": {
			Config: map[string]string{"config_key": "config_value"},
		},
	}

	jpb := j.ToProto()
	assert.Equal(t, jpb.Processors, proc)
}

func Test_isRunnable(t *testing.T) {
	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()

	c := DefaultConfig()
	c.BindAddr = ip1.String()
	c.NodeName = "test1"
	c.Server = true
	c.LogLevel = logLevel
	c.DevMode = true

	a := NewAgent(c)
	a.GRPCClient = &gRPCClientMock{}
	a.Start()
	time.Sleep(2 * time.Second)

	testCases := []struct {
		name string
		job  *Job
		want bool
	}{
		{
			name: "global lock",
			job: &Job{
				Name: "test_job",
				Agent: &Agent{
					GlobalLock: true,
				},
			},
			want: false,
		},
		{
			name: "running forbid",
			job: &Job{
				Name:        "test_job",
				Agent:       a,
				Concurrency: ConcurrencyForbid,
			},
			want: false,
		},
		{
			name: "running true",
			job: &Job{
				Name:    "test_job",
				Agent:   a,
				running: true,
			},
			want: false,
		},
		{
			name: "should run",
			job: &Job{
				Name:  "test_job",
				Agent: a,
			},
			want: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.job.isRunnable())
		})
	}
}

type gRPCClientMock struct {
}

func (gRPCClientMock) Connect(s string) (*grpc.ClientConn, error) { return nil, nil }
func (gRPCClientMock) ExecutionDone(s string, e *Execution) error { return nil }
func (gRPCClientMock) GetJob(s string, a string) (*Job, error)    { return nil, nil }
func (gRPCClientMock) SetJob(j *Job) error                        { return nil }
func (gRPCClientMock) DeleteJob(s string) (*Job, error)           { return nil, nil }
func (gRPCClientMock) Leave(s string) error                       { return nil }
func (gRPCClientMock) RunJob(s string) (*Job, error)              { return nil, nil }
func (gRPCClientMock) RaftGetConfiguration(s string) (*proto.RaftGetConfigurationResponse, error) {
	return nil, nil
}
func (gRPCClientMock) RaftRemovePeerByID(s string, a string) error { return nil }

func (gRPCClientMock) GetActiveExecutions(s string) ([]*proto.Execution, error) {
	return []*proto.Execution{
		&proto.Execution{
			JobName: "test_job",
		},
	}, nil
}
