package dkron

import (
	"testing"
	"time"

	"github.com/abronan/valkeyrie/store"
	"github.com/stretchr/testify/assert"
)

func TestJobGetParent(t *testing.T) {
	s := NewStore(store.Backend(backend), []string{backendMachine}, nil, "dkron-test", nil)
	a := &Agent{
		Store: s,
	}
	s.agent = a

	// Cleanup everything
	err := s.Client().DeleteTree("dkron-test")
	if err != nil && err != store.ErrKeyNotFound {
		t.Logf("error cleaning up: %s", err)
	}

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

	parentTestJob, err = dependentTestJob.GetParent()
	assert.NoError(t, err)
	assert.Equal(t, []string{dependentTestJob.Name}, parentTestJob.DependentJobs)

	ptj, err := dependentTestJob.GetParent()
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

	ptj, err = dtj.GetParent()
	assert.EqualError(t, ErrNoParent, err.Error())

	ptj, err = s.GetJob(parentTestJob.Name, nil)
	assert.NoError(t, err)
	assert.Equal(t, []string{}, ptj.DependentJobs)
}

func TestJobGetNext(t *testing.T) {
	j := Job{
		Schedule: "@daily",
	}

	td := time.Now()
	tonight := time.Date(td.Year(), td.Month(), td.Day()+1, 0, 0, 0, 0, td.Location())
	n, err := j.GetNext()

	assert.NoError(t, err)
	assert.Equal(t, tonight, n)
}
