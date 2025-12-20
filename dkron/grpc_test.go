package dkron

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	typesv1 "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/hashicorp/serf/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
		testExecution.Output = brokenStreamErrorMsg

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
		resp, err := a.GRPCServer.(*GRPCServer).ExecutionDone(ctx, &typesv1.ExecutionDoneRequest{
			Execution: testExecution.ToProto(),
		})
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, []byte("retry"), resp.Payload)
	})
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
