package dkron

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/hashicorp/serf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunQuery(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	advAddr := testutil.GetBindAddr().String()

	c := DefaultConfig()
	c.BindAddr = advAddr + ":5000"
	c.NodeName = "test1"
	c.Server = true
	c.Tags = map[string]string{"role": "test"}
	c.LogLevel = logLevel
	c.DevMode = true
	c.DataDir = dir
	c.BootstrapExpect = 1

	a := NewAgent(c)
	err = a.Start()
	require.NoError(t, err)
	time.Sleep(2 * time.Second)

	// Test error with no job
	_, err = a.RunQuery("foo", &Execution{})
	assert.True(t, errors.Is(err, badger.ErrKeyNotFound))

	j1 := &Job{
		Name:     "test_job",
		Schedule: "@daily",
	}
	err = a.Store.SetJob(j1, false)
	require.NoError(t, err)

	a.sched.Start([]*Job{j1}, a)

	_, err = a.RunQuery("test_job", &Execution{})
	assert.NoError(t, err)

	a.Stop()
}
