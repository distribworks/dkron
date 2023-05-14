package dkron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/buntdb"
)

func TestStore(t *testing.T) {
	log := getTestLogger()
	s, err := NewStore(log)
	require.NoError(t, err)
	defer s.Shutdown()

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

	// Check that we still get an empty job list
	jobs, err := s.GetJobs(nil)
	assert.NoError(t, err)
	assert.NotNil(t, jobs, "jobs nil, expecting empty slice")
	assert.Empty(t, jobs)

	err = s.SetJob(testJob, true)
	assert.NoError(t, err)
	err = s.SetJob(testJob2, true)
	assert.NoError(t, err)

	jobs, err = s.GetJobs(nil)
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

	_, err = s.SetExecution(testExecution)
	require.NoError(t, err)

	testExecution2 := &Execution{
		JobName:    "test2",
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(),
		Success:    true,
		Output:     "test",
		NodeName:   "testNode",
	}
	_, err = s.SetExecution(testExecution2)
	require.NoError(t, err)

	execs, err := s.GetExecutions("test", &ExecutionOptions{
		Sort:  "started_at",
		Order: "DESC",
	})
	assert.NoError(t, err)

	testExecution.Id = testExecution.Key()
	assert.Equal(t, testExecution, execs[0])
	assert.Len(t, execs, 1)

	_, err = s.DeleteJob("test")
	assert.NoError(t, err)

	_, err = s.DeleteJob("test")
	assert.EqualError(t, err, buntdb.ErrNotFound.Error())

	_, err = s.DeleteJob("test2")
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

	storeJob(t, s, "parent1")
	storeChildJob(t, s, "child1", "parent1")

	_, err := s.DeleteJob("parent1")
	assert.EqualError(t, err, ErrDependentJobs.Error())

	deleteJob(t, s, "child1")
	_, err = s.DeleteJob("parent1")
	assert.NoError(t, err)
}

func TestStore_GetJobsWithMetadata(t *testing.T) {
	s := setupStore(t)

	metadata := make(map[string]string)
	metadata["t1"] = "v1"
	storeJobWithMetadata(t, s, "job1", metadata)

	metadata["t2"] = "v2"
	storeJobWithMetadata(t, s, "job2", metadata)

	var options JobOptions
	options.Metadata = make(map[string]string)
	options.Metadata["t1"] = "v1"
	jobs, err := s.GetJobs(&options)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(jobs))

	options.Metadata["t2"] = "v2"
	jobs, err = s.GetJobs(&options)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jobs))
	assert.Equal(t, "job2", jobs[0].Name)

	options.Metadata["t3"] = "v3"
	jobs, err = s.GetJobs(&options)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(jobs))
}

func Test_computeStatus(t *testing.T) {
	log := getTestLogger()
	s, err := NewStore(log)
	require.NoError(t, err)

	n := time.Now()

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
	s.SetExecution(ex1)

	ex2 := &Execution{
		JobName:    "test",
		StartedAt:  n.Add(10 * time.Millisecond),
		FinishedAt: n,
		Success:    false,
		Output:     "test",
		NodeName:   "testNode2",
		Group:      1,
	}
	s.SetExecution(ex2)

	ex3 := &Execution{
		JobName:    "test",
		StartedAt:  n.Add(20 * time.Millisecond),
		FinishedAt: n,
		Success:    true,
		Output:     "test",
		NodeName:   "testNode1",
		Group:      2,
	}
	s.SetExecution(ex3)

	ex4 := &Execution{
		JobName:    "test",
		StartedAt:  n.Add(30 * time.Millisecond),
		FinishedAt: n,
		Success:    true,
		Output:     "test",
		NodeName:   "testNode1",
		Group:      2,
	}
	s.SetExecution(ex4)

	ex5 := &Execution{
		JobName:   "test",
		StartedAt: n.Add(40 * time.Millisecond),
		Success:   false,
		Output:    "test",
		NodeName:  "testNode1",
		Group:     3,
	}
	s.SetExecution(ex5)

	ex6 := &Execution{
		JobName:  "test",
		Success:  false,
		Output:   "test",
		NodeName: "testNode1",
		Group:    4,
	}
	s.SetExecution(ex6)

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

// Following are supporting functions for the tests

func storeJob(t *testing.T, s *Store, jobName string) {
	job := scaffoldJob()
	job.Name = jobName
	require.NoError(t, s.SetJob(job, false))
}

func storeJobWithMetadata(t *testing.T, s *Store, jobName string, metadata map[string]string) {
	job := scaffoldJob()
	job.Name = jobName
	job.Metadata = metadata
	require.NoError(t, s.SetJob(job, false))
}

func storeChildJob(t *testing.T, s *Store, jobName string, parentName string) {
	job := scaffoldJob()
	job.Name = jobName
	job.ParentJob = parentName
	require.NoError(t, s.SetJob(job, false))
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
	s, err := NewStore(log)
	require.NoError(t, err)
	return s
}

func loadJob(t *testing.T, s *Store, name string) *Job {
	job, err := s.GetJob(name, nil)
	require.NoError(t, err)
	return job
}

func deleteJob(t *testing.T, s *Store, name string) {
	_, err := s.DeleteJob(name)
	require.NoError(t, err)
}
