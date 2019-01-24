package dkron

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/stretchr/testify/assert"
)

func setupAPITest(t *testing.T) (a *Agent) {
	c := DefaultConfig()
	c.BindAddr = testutil.GetBindAddr().String()
	c.HTTPAddr = "127.0.0.1:8090"
	c.NodeName = "test"
	c.Server = true
	c.LogLevel = logLevel
	c.Keyspace = "dkron-test"
	c.Backend = "etcdv3"
	c.BackendMachines = []string{os.Getenv("DKRON_BACKEND_MACHINE")}

	a = NewAgent(c, nil)
	a.Start()

	for {
		if a.ready {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(1 * time.Second)

	// clean up the keyspace to ensure clean runs
	a.Store.Client.DeleteTree("dkron-test")

	return
}

func TestAPIJobCreateUpdate(t *testing.T) {
	a := setupAPITest(t)

	jsonStr := []byte(`{"name": "test_job", "schedule": "@every 1m", "executor": "shell", "executor_config": {"command": "date"}, "owner": "mec", "owner_email": "foo@bar.com", "disabled": true}`)

	resp, err := http.Post("http://localhost:8090/v1/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var origJob Job
	if err := json.Unmarshal(body, &origJob); err != nil {
		t.Fatal(err)
	}

	jsonStr1 := []byte(`{"name": "test_job", "schedule": "@every 1m", "executor": "shell", "executor_config": {"command": "test"}, "disabled": false}`)
	resp, err = http.Post("http://localhost:8090/v1/jobs", "encoding/json", bytes.NewBuffer(jsonStr1))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var overwriteJob Job
	if err := json.Unmarshal(body, &overwriteJob); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, origJob.Name, overwriteJob.Name)
	assert.False(t, overwriteJob.Disabled)
	assert.NotEqual(t, origJob.ExecutorConfig["command"], overwriteJob.ExecutorConfig["command"])
	assert.Equal(t, "test", overwriteJob.ExecutorConfig["command"])

	// Send a shutdown request
	a.Stop()
}

func TestAPIJobCreateUpdateParentJob_SameParent(t *testing.T) {
	a := setupAPITest(t)

	jsonStr := []byte(`{
		"name": "test_job",
		"schedule": "@every 1m",
		"command": "date",
		"owner": "mec",
		"owner_email":
		"foo@bar.com",
		"disabled": true,
		"parent_job": "test_job"
	}`)

	resp, err := http.Post("http://localhost:8090/v1/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, 422, resp.StatusCode)
	errJSON, err := json.Marshal(ErrSameParent.Error())
	assert.Contains(t, string(errJSON)+"\n", string(body))

	// Send a shutdown request
	a.Stop()
}

func TestAPIJobCreateUpdateParentJob_NoParent(t *testing.T) {
	a := setupAPITest(t)

	jsonStr := []byte(`{
		"name": "test_job",
		"schedule": "@every 1m",
		"command": "date",
		"owner": "mec",
		"owner_email":
		"foo@bar.com",
		"disabled": true,
		"parent_job": "parent_test_job"
	}`)

	resp, err := http.Post("http://localhost:8090/v1/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, 422, resp.StatusCode)
	errJSON, err := json.Marshal(ErrParentJobNotFound.Error())
	assert.Contains(t, string(errJSON)+"\n", string(body))

	// Send a shutdown request
	a.Stop()
}
