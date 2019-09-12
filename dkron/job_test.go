package dkron

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/distribworks/dkron/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobGetParent(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	a := &Agent{}
	s, err := NewStore(a, dir)
	defer s.Shutdown()
	require.NoError(t, err)
	a.Store = s

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
	assert.Nil(t, ptj.DependentJobs)
}

func TestNewJobFromProto(t *testing.T) {
	testConfig := map[string]PluginConfig{
		"test_processor": {
			"config_key": "config_value",
		},
	}

	in := &proto.Job{}
	proc := map[string]*proto.PluginConfig{
		"test_processor": {
			Config: map[string]string{"config_key": "config_value"},
		},
	}
	in.Processors = proc

	j := NewJobFromProto(in)
	assert.Equal(t, testConfig, j.Processors)
}

func TestToProto(t *testing.T) {
	j := &Job{
		Processors: map[string]PluginConfig{
			"test_processor": {
				"config_key": "config_value",
			},
		},
	}
	proc := map[string]*proto.PluginConfig{
		"test_processor": {
			Config: map[string]string{"config_key": "config_value"},
		},
	}

	jpb := j.ToProto()
	assert.Equal(t, jpb.Processors, proc)
}
