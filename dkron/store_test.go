package dkron

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	s, err := NewStore(nil, dir)
	require.NoError(t, err)
	defer s.Shutdown()

	testJob := &Job{
		Name:           "test",
		Schedule:       "@every 2s",
		Executor:       "shell",
		ExecutorConfig: map[string]string{"command": "/bin/false"},
		Disabled:       true,
	}

	// Check that we still get an empty job list
	jobs, err := s.GetJobs(nil)
	if err != nil {
		t.Fatalf("error getting jobs: %s", err)
	}
	assert.NotNil(t, jobs, "jobs nil, expecting empty slice")

	if err := s.SetJob(testJob, true); err != nil {
		t.Fatalf("error creating job: %s", err)
	}

	jobs, err = s.GetJobs(nil)
	if err != nil {
		t.Fatalf("error getting jobs: %s", err)
	}
	assert.Len(t, jobs, 1)
	assert.Equal(t, "test", jobs[0].Name)

	if _, err := s.DeleteJob("test"); err != nil {
		t.Fatalf("error deleting job: %s", err)
	}

	if _, err := s.DeleteJob("test"); err == nil {
		t.Fatalf("error job deletion should fail: %s", err)
	}

	testExecution := &Execution{
		JobName:    "test",
		StartedAt:  time.Now(),
		FinishedAt: time.Now(),
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode",
	}

	_, err = s.SetExecution(testExecution)
	if err != nil {
		t.Fatalf("error setting the execution: %s", err)
	}

	execs, err := s.GetExecutions("test")
	if err != nil {
		t.Fatalf("error getting executions: %s", err)
	}

	if !execs[0].StartedAt.Equal(testExecution.StartedAt) {
		t.Fatalf("error on retrieved excution expected: %s got: %s", testExecution.StartedAt, execs[0].StartedAt)
	}

	if len(execs) != 1 {
		t.Fatalf("error in number of expected executions: %v", execs)
	}
}

func TestStore_AddDependentJobToParent(t *testing.T) {
	s, dir := setupStore(t)
	defer cleanupStore(dir, s)

	storeJob(t, s, "parent1")
	storeChildJob(t, s, "child1", "parent1")
	parent := loadJob(t, s, "parent1")

	assert.Equal(t, "child1", parent.DependentJobs[0])
}

func TestStore_ParentIsUpdatedAfterDeletingDependentJob(t *testing.T) {
	s, dir := setupStore(t)
	defer cleanupStore(dir, s)

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
	s, dir := setupStore(t)
	defer cleanupStore(dir, s)

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
	s, dir := setupStore(t)
	defer cleanupStore(dir, s)

	storeJob(t, s, "child1")
	storeJob(t, s, "parent1")
	storeChildJob(t, s, "child1", "parent1")
	parent := loadJob(t, s, "parent1")

	assert.Equal(t, "child1", parent.DependentJobs[0])
}

func TestStore_JobBecomesIndependentJob(t *testing.T) {
	s, dir := setupStore(t)
	defer cleanupStore(dir, s)

	storeJob(t, s, "parent1")
	storeChildJob(t, s, "child1", "parent1")
	storeJob(t, s, "child1")
	parent := loadJob(t, s, "parent1")

	assert.Equal(t, 0, len(parent.DependentJobs))
}

func TestStore_GetLastExecutionGroup(t *testing.T) {
	// This can not use time.Now() because that will include monotonic information
	// that will cause the unmarshalled execution to differ from our generated version
	// See `go doc time`
	earlyTime := time.Date(2000, 01, 01, 12, 05, 00, 00, time.UTC)
	middleTime := earlyTime.Add(1 * time.Minute)
	lateTime := earlyTime.Add(1 * time.Hour)

	executionSingleEarly := &Execution{
		JobName:    "test",
		StartedAt:  earlyTime,
		FinishedAt: earlyTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode1",
		Group:      1,
	}
	executionSingleMiddle := &Execution{
		JobName:    "test",
		StartedAt:  middleTime,
		FinishedAt: middleTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode1",
		Group:      2,
	}
	executionGroupMiddle1 := &Execution{
		JobName:    "test",
		StartedAt:  middleTime,
		FinishedAt: middleTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode1",
		Group:      3,
	}
	executionGroupMiddle2 := &Execution{
		JobName:    "test",
		StartedAt:  middleTime,
		FinishedAt: middleTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode2",
		Group:      3,
	}
	executionGroupLater1 := &Execution{
		JobName:    "test",
		StartedAt:  lateTime,
		FinishedAt: lateTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode1",
		Group:      4,
	}
	executionGroupLater2 := &Execution{
		JobName:    "test",
		StartedAt:  lateTime,
		FinishedAt: lateTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode2",
		Group:      4,
	}

	tests := []struct {
		name          string
		jobName       string
		addExecutions []*Execution
		want          []*Execution
		wantErr       bool
	}{
		{
			"Test with one",
			"test",
			[]*Execution{executionSingleEarly},
			[]*Execution{executionSingleEarly},
			false,
		}, {
			"Test with two",
			"test",
			[]*Execution{executionSingleEarly, executionSingleMiddle},
			[]*Execution{executionSingleMiddle},
			false,
		}, {
			"Test with three",
			"test",
			[]*Execution{executionSingleEarly, executionSingleMiddle, executionGroupMiddle1},
			[]*Execution{executionGroupMiddle1},
			false,
		}, {
			"Test with one group",
			"test",
			[]*Execution{executionSingleEarly, executionGroupMiddle1, executionGroupMiddle2},
			[]*Execution{executionGroupMiddle1, executionGroupMiddle2},
			false,
		}, {
			"Test with two groups",
			"test",
			[]*Execution{executionSingleEarly, executionGroupMiddle1, executionGroupMiddle2, executionGroupLater1, executionGroupLater2},
			[]*Execution{executionGroupLater1, executionGroupLater2},
			false,
		}, {
			"Test with none",
			"test",
			[]*Execution{},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "dkron-test")
			require.NoError(t, err)
			s, err := NewStore(nil, dir)
			require.NoError(t, err)

			for _, e := range tt.addExecutions {
				s.SetExecution(e)
			}

			got, err := s.GetLastExecutionGroup(tt.jobName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GetLastExecutionGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, w := range tt.want {
				assert.Contains(t, got, w)
			}

			err = s.Shutdown()
			require.NoError(t, err)
			err = os.RemoveAll(dir)
			require.NoError(t, err)
		})
	}
}

// Following are supporting functions for the tests

func storeJob(t *testing.T, s *Store, jobName string) {
	job := scaffoldJob()
	job.Name = jobName
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

func setupStore(t *testing.T) (*Store, string) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)

	a := NewAgent(nil, nil)
	s, err := NewStore(a, dir)
	require.NoError(t, err)
	a.Store = s

	return s, dir
}

func cleanupStore(dir string, s *Store) {
	s.Shutdown()
	os.RemoveAll(dir)
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
