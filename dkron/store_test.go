package dkron

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/buntdb"
	"go.opentelemetry.io/otel"
)

func TestStore(t *testing.T) {
	log := getTestLogger()
	s, err := NewStore(log, otel.Tracer("test"))
	require.NoError(t, err)
	defer s.Shutdown() // nolint: errcheck

	testJob := &Job{
		Name:           "test",
		Schedule:       "@every 2s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/false"},
		Disabled:       true,
	}

	testJob2 := &Job{
		Name:           "test2",
		Schedule:       "@every 2s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/false"},
		Disabled:       true,
	}

	ctx := context.Background()
	// Check that we still get an empty job list
	jobs, err := s.GetJobs(ctx, nil)
	assert.NoError(t, err)
	assert.NotNil(t, jobs, "jobs nil, expecting empty slice")
	assert.Empty(t, jobs)

	err = s.SetJob(ctx, testJob, true)
	assert.NoError(t, err)
	err = s.SetJob(ctx, testJob2, true)
	assert.NoError(t, err)

	jobs, err = s.GetJobs(ctx, nil)
	assert.NoError(t, err)
	assert.Len(t, jobs, 2)
	assert.Equal(t, "test", jobs[0].Name)

	testExecution := &Execution{
		JobName:    "test",
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(),
		Success:    true,
		Output:     "test",
		NodeName:   "testNode",
	}

	_, err = s.SetExecution(ctx, testExecution)
	require.NoError(t, err)

	testExecution2 := &Execution{
		JobName:    "test2",
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(),
		Success:    true,
		Output:     "test",
		NodeName:   "testNode",
	}
	_, err = s.SetExecution(ctx, testExecution2)
	require.NoError(t, err)

	execs, err := s.GetExecutions(ctx, "test", &ExecutionOptions{
		Sort:  "started_at",
		Order: "DESC",
	})
	assert.NoError(t, err)

	testExecution.Id = testExecution.Key()
	assert.Equal(t, testExecution, execs[0])
	assert.Len(t, execs, 1)

	// Test GetExecution method
	execution, err := s.GetExecution(ctx, "test", testExecution.Key())
	assert.NoError(t, err)
	assert.Equal(t, testExecution, execution)

	// Test GetExecution with non-existent execution
	_, err = s.GetExecution(ctx, "test", "non-existent")
	assert.EqualError(t, err, buntdb.ErrNotFound.Error())

	// Test GetExecution with non-existent job
	_, err = s.GetExecution(ctx, "non-existent-job", testExecution.Key())
	assert.EqualError(t, err, buntdb.ErrNotFound.Error())

	_, err = s.DeleteJob(ctx, "test")
	assert.NoError(t, err)

	_, err = s.DeleteJob(ctx, "test")
	assert.EqualError(t, err, buntdb.ErrNotFound.Error())

	_, err = s.DeleteJob(ctx, "test2")
	assert.NoError(t, err)
}

func TestStore_AddDependentJobToParent(t *testing.T) {
	s := setupStore(t)

	storeJob(t, s, "parent1")
	storeChildJob(t, s, "child1", "parent1")
	parent := loadJob(t, s, "parent1")

	assert.Equal(t, "child1", parent.DependentJobs[0])
}

func TestStore_ParentIsUpdatedAfterDeletingDependentJob(t *testing.T) {
	s := setupStore(t)

	storeJob(t, s, "parent1")
	storeChildJob(t, s, "child1", "parent1")
	parent := loadJob(t, s, "parent1")

	assert.Equal(t, "child1", parent.DependentJobs[0])

	deleteJob(t, s, "child1")
	parent = loadJob(t, s, "parent1")

	// Child has to have been removed from the parent (nr. of dependent jobs is 0)
	assert.Equal(t, 0, len(parent.DependentJobs))
}

func TestStore_DependentJobsUpdatedAfterSwappingParent(t *testing.T) {
	s := setupStore(t)

	storeJob(t, s, "parent1")
	storeChildJob(t, s, "child1", "parent1")
	parent1 := loadJob(t, s, "parent1")

	assert.Equal(t, parent1.DependentJobs[0], "child1")

	storeJob(t, s, "parent2")
	storeChildJob(t, s, "child1", "parent2")
	parent1 = loadJob(t, s, "parent1")

	assert.Equal(t, 0, len(parent1.DependentJobs))

	parent2 := loadJob(t, s, "parent2")

	assert.Equal(t, "child1", parent2.DependentJobs[0])
}

func TestStore_JobBecomesDependentJob(t *testing.T) {
	s := setupStore(t)

	storeJob(t, s, "child1")
	storeJob(t, s, "parent1")
	storeChildJob(t, s, "child1", "parent1")
	parent := loadJob(t, s, "parent1")

	assert.Equal(t, "child1", parent.DependentJobs[0])
}

func TestStore_JobBecomesIndependentJob(t *testing.T) {
	s := setupStore(t)

	storeJob(t, s, "parent1")
	storeChildJob(t, s, "child1", "parent1")
	storeJob(t, s, "child1")
	parent := loadJob(t, s, "parent1")

	assert.Equal(t, 0, len(parent.DependentJobs))
}

func TestStore_ChildIsUpdatedAfterDeletingParentJob(t *testing.T) {
	s := setupStore(t)
	ctx := context.Background()

	storeJob(t, s, "parent1")
	storeChildJob(t, s, "child1", "parent1")

	_, err := s.DeleteJob(ctx, "parent1")
	assert.EqualError(t, err, ErrDependentJobs.Error())

	deleteJob(t, s, "child1")
	_, err = s.DeleteJob(ctx, "parent1")
	assert.NoError(t, err)
}

func TestStore_GetJobsWithMetadata(t *testing.T) {
	s := setupStore(t)

	metadata := make(map[string]string)
	metadata["t1"] = "v1"
	storeJobWithMetadata(t, s, "job1", metadata)

	metadata["t2"] = "v2"
	storeJobWithMetadata(t, s, "job2", metadata)

	ctx := context.Background()

	var options JobOptions
	options.Metadata = make(map[string]string)
	options.Metadata["t1"] = "v1"
	jobs, err := s.GetJobs(ctx, &options)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(jobs))

	options.Metadata["t2"] = "v2"
	jobs, err = s.GetJobs(ctx, &options)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jobs))
	assert.Equal(t, "job2", jobs[0].Name)

	options.Metadata["t3"] = "v3"
	jobs, err = s.GetJobs(ctx, &options)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(jobs))
}

func Test_computeStatus(t *testing.T) {
	log := getTestLogger()
	s, err := NewStore(log, otel.Tracer("test"))
	require.NoError(t, err)

	n := time.Now()
	ctx := context.Background()

	// Prepare executions
	ex1 := &Execution{
		JobName:    "test",
		StartedAt:  n,
		FinishedAt: n,
		Success:    true,
		Output:     "test",
		NodeName:   "testNode1",
		Group:      1,
	}
	_, _ = s.SetExecution(ctx, ex1)

	ex2 := &Execution{
		JobName:    "test",
		StartedAt:  n.Add(10 * time.Millisecond),
		FinishedAt: n,
		Success:    false,
		Output:     "test",
		NodeName:   "testNode2",
		Group:      1,
	}
	_, _ = s.SetExecution(ctx, ex2)

	ex3 := &Execution{
		JobName:    "test",
		StartedAt:  n.Add(20 * time.Millisecond),
		FinishedAt: n,
		Success:    true,
		Output:     "test",
		NodeName:   "testNode1",
		Group:      2,
	}
	_, _ = s.SetExecution(ctx, ex3)

	ex4 := &Execution{
		JobName:    "test",
		StartedAt:  n.Add(30 * time.Millisecond),
		FinishedAt: n,
		Success:    true,
		Output:     "test",
		NodeName:   "testNode1",
		Group:      2,
	}
	_, _ = s.SetExecution(ctx, ex4)

	ex5 := &Execution{
		JobName:   "test",
		StartedAt: n.Add(40 * time.Millisecond),
		Success:   false,
		Output:    "test",
		NodeName:  "testNode1",
		Group:     3,
	}
	_, _ = s.SetExecution(ctx, ex5)

	ex6 := &Execution{
		JobName:  "test",
		Success:  false,
		Output:   "test",
		NodeName: "testNode1",
		Group:    4,
	}
	_, _ = s.SetExecution(ctx, ex6)

	// Tests status
	err = s.db.View(func(tx *buntdb.Tx) error {
		status, _ := s.computeStatus("test", 1, tx)
		assert.Equal(t, StatusPartiallyFailed, status)

		status, _ = s.computeStatus("test", 2, tx)
		assert.Equal(t, StatusSuccess, status)

		status, _ = s.computeStatus("test", 3, tx)
		assert.Equal(t, StatusFailed, status)

		status, _ = s.computeStatus("test", 4, tx)
		assert.Equal(t, StatusFailed, status)

		return nil
	})
	require.NoError(t, err)
}

func TestStore_GetExecutions_SortsWithinJobPrefixRange(t *testing.T) {
	log := getTestLogger()
	s, err := NewStore(log, otel.Tracer("test"))
	require.NoError(t, err)
	defer s.Shutdown() // nolint: errcheck

	ctx := context.Background()
	base := time.Now().UTC()

	older := &Execution{
		JobName:    "daily",
		StartedAt:  base,
		FinishedAt: base,
		Success:    true,
		Output:     "older",
		NodeName:   "node1",
	}
	newer := &Execution{
		JobName:    "daily",
		StartedAt:  base.Add(time.Second),
		FinishedAt: base.Add(time.Second),
		Success:    true,
		Output:     "newer",
		NodeName:   "node1",
	}
	otherJob := &Execution{
		JobName:    "daily_2",
		StartedAt:  base.Add(2 * time.Second),
		FinishedAt: base.Add(2 * time.Second),
		Success:    true,
		Output:     "other",
		NodeName:   "node2",
	}

	_, err = s.SetExecution(ctx, older)
	require.NoError(t, err)
	_, err = s.SetExecution(ctx, newer)
	require.NoError(t, err)
	_, err = s.SetExecution(ctx, otherJob)
	require.NoError(t, err)

	execs, err := s.GetExecutions(ctx, "daily", &ExecutionOptions{
		Sort:  "started_at",
		Order: "DESC",
	})
	require.NoError(t, err)
	require.Len(t, execs, 2)
	assert.Equal(t, "daily", execs[0].JobName)
	assert.Equal(t, "daily", execs[1].JobName)
	assert.Equal(t, newer.Key(), execs[0].Id)
	assert.Equal(t, older.Key(), execs[1].Id)
}

func TestStore_GetRunningExecutions(t *testing.T) {
	log := getTestLogger()
	s, err := NewStore(log, otel.Tracer("test"))
	require.NoError(t, err)
	defer s.Shutdown() // nolint: errcheck

	ctx := context.Background()

	// Create a test job
	testJob := &Job{
		Name:           "test",
		Schedule:       "@every 2s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/false"},
		Disabled:       true,
	}
	err = s.SetJob(ctx, testJob, true)
	require.NoError(t, err)

	// Test 1: No executions - should return empty slice
	runningExecs, err := s.GetRunningExecutions(ctx, "test")
	assert.NoError(t, err)
	assert.Empty(t, runningExecs)

	// Test 2: Add a running execution (StartedAt set, FinishedAt zero)
	runningExecution := &Execution{
		JobName:    "test",
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Time{}, // Zero time means not finished
		Success:    false,
		Output:     "running",
		NodeName:   "testNode1",
		Group:      time.Now().UnixNano(),
		Attempt:    1,
	}
	_, err = s.SetExecution(ctx, runningExecution)
	require.NoError(t, err)

	// Should find the running execution
	runningExecs, err = s.GetRunningExecutions(ctx, "test")
	assert.NoError(t, err)
	assert.Len(t, runningExecs, 1)
	assert.Equal(t, "test", runningExecs[0].JobName)
	assert.Equal(t, "testNode1", runningExecs[0].NodeName)

	// Test 3: Add a finished execution (both StartedAt and FinishedAt set)
	finishedExecution := &Execution{
		JobName:    "test",
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(), // Finished
		Success:    true,
		Output:     "finished",
		NodeName:   "testNode2",
		Group:      time.Now().UnixNano(),
		Attempt:    1,
	}
	_, err = s.SetExecution(ctx, finishedExecution)
	require.NoError(t, err)

	// Should still find only the running execution
	runningExecs, err = s.GetRunningExecutions(ctx, "test")
	assert.NoError(t, err)
	assert.Len(t, runningExecs, 1)
	assert.Equal(t, "testNode1", runningExecs[0].NodeName)

	// Test 4: Mark the running execution as done
	runningExecution.FinishedAt = time.Now().UTC()
	runningExecution.Success = true
	_, err = s.SetExecutionDone(ctx, runningExecution)
	require.NoError(t, err)

	// Should now find no running executions
	runningExecs, err = s.GetRunningExecutions(ctx, "test")
	assert.NoError(t, err)
	assert.Empty(t, runningExecs)

	// Test 5: Test with non-existent job - should return empty slice, not error
	runningExecs, err = s.GetRunningExecutions(ctx, "nonexistent")
	assert.NoError(t, err)
	assert.Empty(t, runningExecs)
}

func TestSetExecutionDonePrometheusMetrics(t *testing.T) {
	s := setupStore(t)
	defer s.Shutdown() // nolint: errcheck
	ctx := context.Background()

	successJobName := fmt.Sprintf("success-%d", time.Now().UnixNano())
	failedJobName := fmt.Sprintf("failed-%d", time.Now().UnixNano())

	require.NoError(t, s.SetJob(ctx, &Job{
		Name:           successJobName,
		Schedule:       "@every 2s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/true"},
		Disabled:       true,
	}, true))

	require.NoError(t, s.SetJob(ctx, &Job{
		Name:           failedJobName,
		Schedule:       "@every 2s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/false"},
		Disabled:       true,
	}, true))

	successBefore := getPrometheusCounterValue(t, "dkron_job_executions_succeeded_total", "job_name", successJobName)
	failedBefore := getPrometheusCounterValue(t, "dkron_job_executions_failed_total", "job_name", failedJobName)

	_, err := s.SetExecutionDone(ctx, &Execution{
		JobName:    successJobName,
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(),
		Success:    true,
		NodeName:   "testNode",
		Group:      time.Now().UnixNano(),
		Attempt:    1,
	})
	require.NoError(t, err)

	_, err = s.SetExecutionDone(ctx, &Execution{
		JobName:    failedJobName,
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(),
		Success:    false,
		NodeName:   "testNode",
		Group:      time.Now().UnixNano(),
		Attempt:    1,
	})
	require.NoError(t, err)

	assert.Equal(t, successBefore+1, getPrometheusCounterValue(t, "dkron_job_executions_succeeded_total", "job_name", successJobName))
	assert.Equal(t, failedBefore+1, getPrometheusCounterValue(t, "dkron_job_executions_failed_total", "job_name", failedJobName))
}

func getPrometheusCounterValue(t *testing.T, metricName, labelName, labelValue string) float64 {
	t.Helper()

	mfs, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	for _, mf := range mfs {
		if mf.GetName() != metricName {
			continue
		}
		for _, metric := range mf.GetMetric() {
			for _, label := range metric.GetLabel() {
				if label.GetName() == labelName && label.GetValue() == labelValue {
					return metric.GetCounter().GetValue()
				}
			}
		}
	}

	return 0
}

// Following are supporting functions for the tests

func storeJob(t *testing.T, s *Store, jobName string) {
	job := scaffoldJob()
	job.Name = jobName
	require.NoError(t, s.SetJob(context.Background(), job, false))
}

func storeJobWithMetadata(t *testing.T, s *Store, jobName string, metadata map[string]string) {
	job := scaffoldJob()
	job.Name = jobName
	job.Metadata = metadata
	require.NoError(t, s.SetJob(context.Background(), job, false))
}

func storeChildJob(t *testing.T, s *Store, jobName string, parentName string) {
	job := scaffoldJob()
	job.Name = jobName
	job.ParentJob = parentName
	require.NoError(t, s.SetJob(context.Background(), job, false))
}

func scaffoldJob() *Job {
	return &Job{
		Name:           "test",
		Schedule:       "@every 1m",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/false"},
		Disabled:       true,
	}
}

func setupStore(t *testing.T) *Store {
	log := getTestLogger()
	s, err := NewStore(log, otel.Tracer("test"))
	require.NoError(t, err)
	return s
}

func loadJob(t *testing.T, s *Store, name string) *Job {
	job, err := s.GetJob(context.Background(), name, nil)
	require.NoError(t, err)
	return job
}

func deleteJob(t *testing.T, s *Store, name string) {
	_, err := s.DeleteJob(context.Background(), name)
	require.NoError(t, err)
}
