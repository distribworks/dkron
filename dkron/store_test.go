package dkron

import (
	"reflect"
	"testing"
	"time"

	"github.com/abronan/valkeyrie/store"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	s := NewStore("etcd", []string{etcdAddr}, nil, "dkron-test", nil)

	// Cleanup everything
	if err := cleanTestKVSpace(s); err != nil {
		t.Logf("error cleaning up: %v", err)
	}

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

func TestStore_GetLastExecutionGroup(t *testing.T) {
	s := createTestStore()

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
	}
	executionSingleMiddle := &Execution{
		JobName:    "test",
		StartedAt:  middleTime,
		FinishedAt: middleTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode1",
	}
	executionGroupMiddle1 := &Execution{
		JobName:    "test",
		StartedAt:  middleTime,
		FinishedAt: middleTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode1",
	}
	executionGroupMiddle2 := &Execution{
		JobName:    "test",
		StartedAt:  middleTime,
		FinishedAt: middleTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode2",
	}
	executionGroupLater1 := &Execution{
		JobName:    "test",
		StartedAt:  lateTime,
		FinishedAt: lateTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode1",
	}
	executionGroupLater2 := &Execution{
		JobName:    "test",
		StartedAt:  lateTime,
		FinishedAt: lateTime,
		Success:    true,
		Output:     []byte("type"),
		NodeName:   "testNode2",
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
			if err := cleanTestKVSpace(s); err != nil {
				t.Logf("error cleaning up: %v", err)
			}
			for _, e := range tt.addExecutions {
				s.SetExecution(e)
			}

			got, err := s.GetLastExecutionGroup(tt.jobName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GetLastExecutionGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.GetLastExecutionGroup() = %+v, want %+v", got, tt.want)
				for i, e := range got {
					t.Logf("Got %d: %v", i, e)
				}
				for i, e := range tt.want {
					t.Logf("Want %d: %v", i, e)
				}
			}
		})
	}
}

func createTestStore() *Store {
	return NewStore("etcd", []string{etcdAddr}, nil, "dkron-test", nil)
}

func cleanTestKVSpace(s *Store) error {
	err := s.Client.DeleteTree("dkron-test")
	if err != nil && err != store.ErrKeyNotFound {
		return err
	}
	return nil
}
