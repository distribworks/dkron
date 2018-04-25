package dkron

import (
	"testing"

	s "github.com/abronan/valkeyrie/store"
	"github.com/stretchr/testify/assert"
)

func TestJobGetParent(t *testing.T) {
	store := NewStore("etcd", []string{etcdAddr}, nil, "dkron-test", nil)
	a := &Agent{
		Store: store,
	}
	store.agent = a

	// Cleanup everything
	err := store.Client.DeleteTree("dkron-test")
	if err != s.ErrKeyNotFound {
		t.Logf("error cleaning up: %s", err)
	}

	parentTestJob := &Job{
		Name:     "parent_test",
		Command:  "/bin/false",
		Schedule: "@every 2s",
	}

	if err := store.SetJob(parentTestJob, nil); err != nil {
		t.Fatalf("error creating job: %s", err)
	}

	dependentTestJob := &Job{
		Name:      "dependent_test",
		Command:   "/bin/false",
		ParentJob: "parent_test",
	}

	err = store.SetJob(dependentTestJob, nil)
	assert.NoError(t, err)

	err = store.SetJobDependencyTree(dependentTestJob, nil)
	assert.NoError(t, err)

	parentTestJob, err = dependentTestJob.GetParent()
	assert.NoError(t, err)
	assert.Equal(t, []string{dependentTestJob.Name}, parentTestJob.DependentJobs)

	ptj, err := dependentTestJob.GetParent()
	assert.NoError(t, err)
	assert.Equal(t, parentTestJob, ptj)

	// Remove the parent job
	ej, _ := store.GetJob(dependentTestJob.Name)

	dependentTestJob.ParentJob = ""
	dependentTestJob.Schedule = "@every 2m"
	err = store.SetJob(dependentTestJob, nil)
	assert.NoError(t, err)

	err = store.SetJobDependencyTree(dependentTestJob, ej)
	assert.NoError(t, err)

	dtj, _ := store.GetJob(dependentTestJob.Name)
	assert.NoError(t, err)
	assert.Equal(t, "", dtj.ParentJob)

	ptj, err = dtj.GetParent()
	assert.EqualError(t, ErrNoParent, err.Error())

	ptj, err = store.GetJob(parentTestJob.Name)
	assert.NoError(t, err)
	assert.Equal(t, []string{}, ptj.DependentJobs)
}
