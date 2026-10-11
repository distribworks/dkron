package dkron

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/buntdb"
	"go.opentelemetry.io/otel"
)

func TestStoreDeleteExecutions(t *testing.T) {
	log := getTestLogger()
	s, err := NewStore(log, otel.Tracer("test"))
	require.NoError(t, err)
	defer s.Shutdown()

	ctx := context.Background()

	// Create a test job
	testJob := &Job{
		Name:           "test_delete_exec",
		Schedule:       "@every 1m",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "date"},
	}
	err = s.SetJob(ctx, testJob, true)
	require.NoError(t, err)

	// Add executions
	exec1 := &Execution{
		JobName:    "test_delete_exec",
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(),
		Success:    true,
		Output:     "success output",
		NodeName:   "node1",
	}
	exec2 := &Execution{
		JobName:    "test_delete_exec",
		StartedAt:  time.Now().UTC().Add(1 * time.Second),
		FinishedAt: time.Now().UTC().Add(1 * time.Second),
		Success:    false,
		Output:     "error output",
		NodeName:   "node1",
	}

	_, err = s.SetExecution(ctx, exec1)
	require.NoError(t, err)
	_, err = s.SetExecution(ctx, exec2)
	require.NoError(t, err)

	// Mark executions as done to update counters
	_, err = s.SetExecutionDone(ctx, exec1)
	require.NoError(t, err)
	_, err = s.SetExecutionDone(ctx, exec2)
	require.NoError(t, err)

	// Verify executions exist and counters are set
	executions, err := s.GetExecutions(ctx, "test_delete_exec", &ExecutionOptions{})
	require.NoError(t, err)
	assert.Equal(t, 2, len(executions))

	job, err := s.GetJob(ctx, "test_delete_exec", nil)
	require.NoError(t, err)
	assert.Equal(t, 1, job.SuccessCount)
	assert.Equal(t, 1, job.ErrorCount)
	assert.True(t, job.LastSuccess.HasValue())
	assert.True(t, job.LastError.HasValue())

	// Delete executions
	err = s.DeleteExecutions(ctx, "test_delete_exec")
	require.NoError(t, err)

	// Verify executions are deleted
	executions, err = s.GetExecutions(ctx, "test_delete_exec", &ExecutionOptions{})
	// When all executions are deleted, GetExecutions returns ErrNotFound
	if err == nil {
		assert.Equal(t, 0, len(executions))
	} else {
		// Alternatively, it might return ErrNotFound which is also valid
		assert.Error(t, err)
	}

	// Verify counters are reset
	job, err = s.GetJob(ctx, "test_delete_exec", nil)
	require.NoError(t, err)
	assert.Equal(t, 0, job.SuccessCount)
	assert.Equal(t, 0, job.ErrorCount)
	assert.False(t, job.LastSuccess.HasValue())
	assert.False(t, job.LastError.HasValue())
}

func TestStoreDeleteExecutionsNonExistentJob(t *testing.T) {
	log := getTestLogger()
	s, err := NewStore(log, otel.Tracer("test"))
	require.NoError(t, err)
	defer s.Shutdown()

	ctx := context.Background()

	// Try to delete executions for a non-existent job
	err = s.DeleteExecutions(ctx, "non_existent_job")
	assert.Error(t, err) // Should error because job doesn't exist
}

func TestStoreDeleteExecutionsDoesNotDeletePrefixedJobExecutions(t *testing.T) {
	log := getTestLogger()
	s, err := NewStore(log, otel.Tracer("test"))
	require.NoError(t, err)
	defer s.Shutdown()

	ctx := context.Background()

	for _, jobName := range []string{"daily", "daily_2"} {
		err = s.SetJob(ctx, &Job{
			Name:           jobName,
			Schedule:       "@every 1m",
			Executor:       "shell",
			ExecutorConfig: map[string]string{"command": "date"},
		}, true)
		require.NoError(t, err)
	}

	for _, exec := range []*Execution{
		{
			JobName:    "daily",
			StartedAt:  time.Now().UTC(),
			FinishedAt: time.Now().UTC(),
			Success:    true,
			Output:     "daily",
			NodeName:   "node1",
		},
		{
			JobName:    "daily_2",
			StartedAt:  time.Now().UTC().Add(time.Second),
			FinishedAt: time.Now().UTC().Add(time.Second),
			Success:    true,
			Output:     "daily_2",
			NodeName:   "node1",
		},
	} {
		_, err = s.SetExecution(ctx, exec)
		require.NoError(t, err)
	}

	err = s.DeleteExecutions(ctx, "daily")
	require.NoError(t, err)

	_, err = s.GetExecutions(ctx, "daily", &ExecutionOptions{})
	assert.ErrorIs(t, err, buntdb.ErrNotFound)

	executions, err := s.GetExecutions(ctx, "daily_2", &ExecutionOptions{})
	require.NoError(t, err)
	require.Len(t, executions, 1)
	assert.Equal(t, "daily_2", executions[0].JobName)
}
