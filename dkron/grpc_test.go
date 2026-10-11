package dkron

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"sync/atomic"
	"testing"
	"time"

	typesv1 "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/distribworks/dkron/v4/plugin"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type executionDoneTestServer struct {
	typesv1.UnimplementedDkronServer
	err     error
	handler func(*typesv1.ExecutionDoneRequest) error
	calls   atomic.Int32
}

type countingProcessor struct {
	calls atomic.Int32
}

func (p *countingProcessor) Process(args *plugin.ProcessorArgs) *typesv1.Execution {
	p.calls.Add(1)
	return args.Execution
}

func (s *executionDoneTestServer) ExecutionDone(_ context.Context, req *typesv1.ExecutionDoneRequest) (*typesv1.ExecutionDoneResponse, error) {
	s.calls.Add(1)
	if s.err != nil {
		return nil, s.err
	}
	if s.handler != nil {
		if err := s.handler(req); err != nil {
			return nil, err
		}
	}
	return &typesv1.ExecutionDoneResponse{Payload: []byte("saved")}, nil
}

func startExecutionDoneTestServer(t *testing.T, server typesv1.DkronServer) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	grpcServer := grpc.NewServer()
	typesv1.RegisterDkronServer(grpcServer, server)
	go func() {
		_ = grpcServer.Serve(listener)
	}()
	t.Cleanup(func() {
		grpcServer.Stop()
		_ = listener.Close()
	})

	return listener.Addr().String()
}

type executionDoneClientMock struct {
	gRPCClientMock
	err   error
	calls int
}

func (m *executionDoneClientMock) ExecutionDone(string, *Execution) error {
	m.calls++
	return m.err
}

func TestGRPCExecutionDone(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	viper.Reset()

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()
	aAddr := ip1.String()

	c := DefaultConfig()
	c.BindAddr = aAddr
	c.NodeName = "test-grpc"
	c.Server = true
	c.LogLevel = logLevel
	c.BootstrapExpect = 1
	c.DevMode = true
	c.DataDir = dir
	c.HTTPAddr = "127.0.0.1:0"

	a := NewAgent(c)
	_ = a.Start()

	for {
		if a.IsLeader() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	testJob := &Job{
		Name:           "test",
		Schedule:       "@manually",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/true"},
		Disabled:       true,
	}

	ctx := context.Background()

	err = a.Store.SetJob(ctx, testJob, true)
	require.NoError(t, err)

	testChildJob := &Job{
		Name:           "child-test",
		ParentJob:      testJob.Name,
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/true"},
		Disabled:       false,
	}

	err = a.Store.SetJob(ctx, testChildJob, true)
	require.NoError(t, err)

	testExecution := &Execution{
		JobName:    testJob.Name,
		Group:      time.Now().UnixNano(),
		StartedAt:  time.Now(),
		NodeName:   "testNode",
		FinishedAt: time.Now(),
		Success:    true,
		Output:     "test",
	}

	log := getTestLogger()
	rc := NewGRPCClient(nil, a, log)

	t.Run("Should run job", func(t *testing.T) {
		err = rc.ExecutionDone(a.advertiseRPCAddr(), testExecution)
		require.NoError(t, err)

		execs, err := a.Store.GetExecutions(ctx, "test", &ExecutionOptions{})
		require.NoError(t, err)

		assert.Len(t, execs, 1)
		assert.Equal(t, string(testExecution.Output), string(execs[0].Output))
	})

	t.Run("Should run a dependent job", func(t *testing.T) {
		execs, err := a.Store.GetExecutions(ctx, "child-test", &ExecutionOptions{})
		require.NoError(t, err)

		assert.Len(t, execs, 1)
	})

	t.Run("Should store execution on a deleted job", func(t *testing.T) {
		// Test job with dependents no delete
		_, err = a.Store.DeleteJob(ctx, testJob.Name)
		require.Error(t, err)

		// Remove dependents and parent
		_, err = a.Store.DeleteJob(ctx, testChildJob.Name)
		require.NoError(t, err)
		_, err = a.Store.DeleteJob(ctx, testJob.Name)
		require.NoError(t, err)

		// Test store execution on a deleted job
		testExecution.FinishedAt = time.Now()
		err = rc.ExecutionDone(a.advertiseRPCAddr(), testExecution)

		assert.Error(t, err, ErrExecutionDoneForDeletedJob)
	})

	t.Run("Test ephemeral jobs", func(t *testing.T) {
		testJob.Ephemeral = true

		err = a.Store.SetJob(ctx, testJob, true)
		require.NoError(t, err)

		err = rc.ExecutionDone(a.advertiseRPCAddr(), testExecution)
		assert.NoError(t, err)

		j, err := a.Store.GetJob(ctx, "test", nil)
		assert.Error(t, err)
		assert.Nil(t, j)
	})

	t.Run("Test job with non-existent dependent", func(t *testing.T) {
		testJob.Name = "test2"
		testJob.DependentJobs = []string{"non-existent"}
		testExecution.JobName = testJob.Name

		err = a.Store.SetJob(ctx, testJob, true)
		require.NoError(t, err)

		err = rc.ExecutionDone(a.advertiseRPCAddr(), testExecution)
		assert.Error(t, err)
	})

	t.Run("Test job retry with broken stream error", func(t *testing.T) {
		// Use the actual error format that would be returned when a broken stream occurs
		brokenStreamErrorMsg := ErrBrokenStream.Error() + ": rpc error: code = Internal desc = grpc: error while marshaling"

		testJob.Name = "test-retry"
		testJob.Schedule = "0 * * * * *" // Every minute at 0 seconds (6-field format)
		testJob.Retries = 2
		testJob.DependentJobs = nil
		testJob.Ephemeral = false
		testJob.Disabled = false
		testExecution.JobName = testJob.Name
		testExecution.Success = false
		testExecution.Attempt = 1
		testExecution.NodeName = a.config.NodeName // Use the agent's node name
		testExecution.FinishedAt = time.Time{}
		testExecution.Output = ""

		err = a.Store.SetJob(ctx, testJob, true)
		require.NoError(t, err)

		// Add job to scheduler so it can be retrieved for retry
		job := NewJobFromProto(testJob.ToProto(), a.logger)
		job.Agent = a
		err = a.sched.AddJob(job)
		require.NoError(t, err)

		// Store initial execution to establish group
		_, err = a.Store.SetExecution(ctx, testExecution)
		require.NoError(t, err)

		// Call ExecutionDone with a failed execution that has a broken stream error
		// This should trigger a retry since Retries > 0
		testExecution.FinishedAt = time.Now().UTC()
		testExecution.Output = brokenStreamErrorMsg
		resp, err := a.GRPCServer.(*GRPCServer).ExecutionDone(ctx, &typesv1.ExecutionDoneRequest{
			Execution: testExecution.ToProto(),
		})
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, []byte("retry"), resp.Payload)
	})
	t.Run("Test stop_on_failure disables job on failure", func(t *testing.T) {
		testJob.Name = "test-stop-on-failure"
		testJob.Schedule = "0 * * * * *"
		testJob.Retries = 0
		testJob.DependentJobs = nil
		testJob.Ephemeral = false
		testJob.Disabled = false
		testJob.StopOnFailure = true
		testExecution.JobName = testJob.Name
		testExecution.Success = false
		testExecution.Attempt = 1
		testExecution.NodeName = a.config.NodeName
		testExecution.Output = "failure output"

		err = a.Store.SetJob(ctx, testJob, true)
		require.NoError(t, err)

		// Add job to scheduler
		job := NewJobFromProto(testJob.ToProto(), a.logger)
		job.Agent = a
		err = a.sched.AddJob(job)
		require.NoError(t, err)

		// Store initial execution to establish group
		_, err = a.Store.SetExecution(ctx, testExecution)
		require.NoError(t, err)

		// Call ExecutionDone with a failed execution - retries are exhausted, stop_on_failure is true
		resp, err := a.GRPCServer.(*GRPCServer).ExecutionDone(ctx, &typesv1.ExecutionDoneRequest{
			Execution: testExecution.ToProto(),
		})
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, []byte("saved"), resp.Payload)

		// Verify that the job was disabled
		disabledJob, err := a.Store.GetJob(ctx, testJob.Name, nil)
		require.NoError(t, err)
		assert.True(t, disabledJob.Disabled, "job should be disabled after stop_on_failure")
	})
}

func TestGRPCSetExecution_returns_error_when_raft_unavailable(t *testing.T) {
	// Given
	server := &GRPCServer{
		agent: &Agent{
			isLeaderFn: func() bool { return true },
		},
		logger: getTestLogger(),
	}
	execution := &typesv1.Execution{
		JobName:  "test",
		NodeName: "testNode",
	}

	// When
	var err error
	require.NotPanics(t, func() {
		_, err = server.SetExecution(context.Background(), execution)
	})

	// Then
	require.EqualError(t, err, "raft apply unavailable")
}

func TestGRPCSetExecution_returns_error_when_not_leader(t *testing.T) {
	// Given
	server := &GRPCServer{
		agent:  &Agent{},
		logger: getTestLogger(),
	}
	execution := &typesv1.Execution{
		JobName:  "test",
		NodeName: "testNode",
	}

	// When
	_, err := server.SetExecution(context.Background(), execution)

	// Then
	require.ErrorIs(t, err, ErrNotLeader)
}

func TestGRPCExecutionDoneFollowerOnlyAcknowledgesSuccessfulForward(t *testing.T) {
	execution := &Execution{
		JobName:    "forwarded-completion",
		NodeName:   "agent-1",
		StartedAt:  time.Now().UTC().Add(-time.Minute),
		FinishedAt: time.Now().UTC(),
		Group:      time.Now().UnixNano(),
	}

	t.Run("forward failure is returned", func(t *testing.T) {
		forwardErr := status.Error(codes.Unavailable, "leader changed")
		client := &executionDoneClientMock{err: forwardErr}
		server := &GRPCServer{
			agent: &Agent{
				GRPCClient: client,
				isLeaderFn: func() bool { return false },
				leaderFn:   func() raft.ServerAddress { return "old-leader:6868" },
			},
			logger: getTestLogger(),
		}

		response, err := server.ExecutionDone(context.Background(), &typesv1.ExecutionDoneRequest{Execution: execution.ToProto()})

		require.Error(t, err)
		assert.Equal(t, codes.Unavailable, status.Code(err))
		assert.Contains(t, err.Error(), forwardErr.Error())
		assert.Nil(t, response)
		assert.Equal(t, 1, client.calls)
	})

	t.Run("successful forward is acknowledged", func(t *testing.T) {
		client := &executionDoneClientMock{}
		server := &GRPCServer{
			agent: &Agent{
				GRPCClient: client,
				isLeaderFn: func() bool { return false },
				leaderFn:   func() raft.ServerAddress { return "current-leader:6868" },
			},
			logger: getTestLogger(),
		}

		response, err := server.ExecutionDone(context.Background(), &typesv1.ExecutionDoneRequest{Execution: execution.ToProto()})

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, []byte("forwarded"), response.Payload)
		assert.Equal(t, 1, client.calls)
	})
}

func TestGRPCExecutionDoneDuplicateSkipsCompletionSideEffects(t *testing.T) {
	ctx := context.Background()
	store, err := NewStore(getTestLogger(), otel.Tracer("test"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = store.Shutdown() })

	processor := &countingProcessor{}
	job := scaffoldJob()
	job.Name = "duplicate-completion"
	job.Processors = map[string]plugin.Config{"counting": {}}
	require.NoError(t, store.SetJob(ctx, job, false))

	execution := &Execution{
		JobName:    job.Name,
		NodeName:   "agent-1",
		StartedAt:  time.Now().UTC().Add(-time.Minute),
		FinishedAt: time.Now().UTC(),
		Success:    true,
		Group:      time.Now().UnixNano(),
		Attempt:    1,
	}
	_, err = store.SetExecutionDone(ctx, execution)
	require.NoError(t, err)

	server := &GRPCServer{
		agent: &Agent{
			Store:            store,
			ProcessorPlugins: map[string]plugin.Processor{"counting": processor},
			config:           DefaultConfig(),
			isLeaderFn:       func() bool { return true },
		},
		logger: getTestLogger(),
	}

	response, err := server.ExecutionDone(ctx, &typesv1.ExecutionDoneRequest{Execution: execution.ToProto()})

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, []byte("duplicate"), response.Payload)
	assert.Zero(t, processor.calls.Load())
}

func TestExecutionDoneRetriesCurrentLeaderAfterNotLeaderResponse(t *testing.T) {
	ctx := context.Background()
	store, err := NewStore(getTestLogger(), otel.Tracer("test"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = store.Shutdown() })

	job := scaffoldJob()
	job.Name = "leader-change-completion"
	require.NoError(t, store.SetJob(ctx, job, false))

	execution := &Execution{
		JobName:    job.Name,
		NodeName:   "ephemeral-agent",
		StartedAt:  time.Now().UTC().Add(-time.Minute),
		FinishedAt: time.Time{},
		Group:      time.Now().UnixNano(),
		Attempt:    1,
	}
	_, err = store.SetExecution(ctx, execution)
	require.NoError(t, err)

	staleFollower := &executionDoneTestServer{err: ErrNotLeader}
	staleFollowerAddr := startExecutionDoneTestServer(t, staleFollower)
	currentLeader := &executionDoneTestServer{
		handler: func(req *typesv1.ExecutionDoneRequest) error {
			_, err := store.SetExecutionDone(ctx, NewExecutionFromProto(req.Execution))
			return err
		},
	}
	currentLeaderAddr := startExecutionDoneTestServer(t, currentLeader)

	clientConfig := DefaultConfig()
	clientConfig.AgentRunMaxRetries = 1
	clientConfig.AgentRunRetryInitialInterval = 0
	clientConfig.AgentRunRetryMaxInterval = 0
	routingAgent := &Agent{
		config: clientConfig,
		leaderFn: func() raft.ServerAddress {
			return raft.ServerAddress(currentLeaderAddr)
		},
	}
	client := NewGRPCClient(nil, routingAgent, getTestLogger())

	execution.FinishedAt = time.Now().UTC()
	execution.Success = true
	require.NoError(t, client.ExecutionDone(staleFollowerAddr, execution))
	assert.Equal(t, int32(1), staleFollower.calls.Load())
	assert.Equal(t, int32(1), currentLeader.calls.Load())

	running, err := store.GetRunningExecutions(ctx, job.Name)
	require.NoError(t, err)
	assert.Empty(t, running)

	stored, err := store.GetExecution(ctx, job.Name, execution.Key())
	require.NoError(t, err)
	assert.False(t, stored.FinishedAt.IsZero())
	assert.True(t, stored.Success)
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "gRPC Unavailable status code",
			err:      status.Error(codes.Unavailable, "transport is closing"),
			expected: true,
		},
		{
			name:     "gRPC DeadlineExceeded status code",
			err:      status.Error(codes.DeadlineExceeded, "deadline exceeded"),
			expected: true,
		},
		{
			name:     "gRPC ResourceExhausted status code",
			err:      status.Error(codes.ResourceExhausted, "quota exceeded"),
			expected: true,
		},
		{
			name:     "gRPC Aborted status code",
			err:      status.Error(codes.Aborted, "transaction aborted"),
			expected: true,
		},
		{
			name:     "gRPC Internal status code",
			err:      status.Error(codes.Internal, "internal error"),
			expected: true,
		},
		{
			name:     "gRPC InvalidArgument status code",
			err:      status.Error(codes.InvalidArgument, "bad request"),
			expected: false,
		},
		{
			name:     "gRPC NotFound status code",
			err:      status.Error(codes.NotFound, "not found"),
			expected: false,
		},
		{
			name:     "transport is closing",
			err:      errors.New("transport is closing"),
			expected: true,
		},
		{
			name:     "connection refused",
			err:      errors.New("connection refused"),
			expected: true,
		},
		{
			name:     "connection reset",
			err:      errors.New("connection reset by peer"),
			expected: true,
		},
		{
			name:     "broken pipe",
			err:      errors.New("broken pipe"),
			expected: true,
		},
		{
			name:     "context deadline exceeded",
			err:      errors.New("context deadline exceeded"),
			expected: true,
		},
		{
			name:     "non-retryable error",
			err:      errors.New("some other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryableError(tt.err)
			assert.Equal(t, tt.expected, result, "isRetryableError(%v) = %v, want %v", tt.err, result, tt.expected)
		})
	}
}
