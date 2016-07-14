package dkron

import (
	"testing"

	s "github.com/docker/libkv/store"
	"github.com/stretchr/testify/assert"
)

func TestJobGetParent(t *testing.T) {
	store := NewStore("etcd", []string{etcdAddr}, nil, "dkron-test")
	a := &AgentCommand{
		store: store,
	}
	store.agent = a

	// Cleanup everything
	err := store.Client.DeleteTree("dkron-test")
	if err != s.ErrKeyNotFound {
		t.Logf("error cleaning up: %s", err)
	}

	parentTestJob := &Job{
		Name:     "parent_test",
		Schedule: "@every 2s",
	}

	if err := store.SetJob(parentTestJob); err != nil {
		t.Fatalf("error creating job: %s", err)
	}

	dependentTestJob := &Job{
		Name:      "dependent_test",
		ParentJob: "parent_test",
	}

	if err := store.SetJob(dependentTestJob); err != nil {
		t.Fatalf("error creating job: %s", err)
	}
	dependentTestJob.Agent = a

	parentTestJob, err = dependentTestJob.GetParent()
	assert.NoError(t, err)
	assert.Equal(t, []string{"dependent_test"}, parentTestJob.DependentJobs)

	ptj, err := dependentTestJob.GetParent()
	assert.NoError(t, err)
	assert.Equal(t, parentTestJob, ptj)

	// Remove the parent job
	dependentTestJob.ParentJob = ""
	err = store.SetJob(dependentTestJob)
	assert.NoError(t, err)

	dtj, err := store.GetJob("dependent_test")
	assert.NoError(t, err)
	assert.Equal(t, "", dtj.ParentJob)

	ptj, err = dtj.GetParent()
	assert.EqualError(t, ErrNoParent, err.Error())

	ptj, err = store.GetJob("parent_test")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, ptj.DependentJobs)
}
