package dkron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAPITest(port string) (a *Agent) {
	c := DefaultConfig()
	c.BindAddr = testutil.GetBindAddr().String()
	c.HTTPAddr = fmt.Sprintf("127.0.0.1:%s", port)
	c.NodeName = "test"
	c.Server = true
	c.LogLevel = logLevel
	c.BootstrapExpect = 1
	c.DevMode = true
	c.DataDir = "dkron-test-" + port + ".data"

	a = NewAgent(c, nil)
	a.Start()

	for {
		if a.IsLeader() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(1 * time.Second)

	return
}

func TestAPIJobCreateUpdate(t *testing.T) {
	port := "8091"
	baseURL := fmt.Sprintf("http://localhost:%s/v1", port)
	setupAPITest(port)

	jsonStr := []byte(`{
		"name": "test_job",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "date"},
		"owner": "mec",
		"owner_email": "foo@bar.com",
		"disabled": true
	}`)

	resp, err := http.Post(baseURL+"/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
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
	resp, err = http.Post(baseURL+"/jobs", "encoding/json", bytes.NewBuffer(jsonStr1))
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
	//a.Stop()
}

func TestAPIJobCreateUpdateParentJob_SameParent(t *testing.T) {
	port := "8092"
	baseURL := fmt.Sprintf("http://localhost:%s/v1", port)
	setupAPITest(port)

	jsonStr := []byte(`{
		"name": "test_job",
		"schedule": "@every 1m",
		"command": "date",
		"owner": "mec",
		"owner_email": "foo@bar.com",
		"disabled": true,
		"parent_job": "test_job"
	}`)

	resp, err := http.Post(baseURL+"/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, 422, resp.StatusCode)
	errJSON, err := json.Marshal(ErrSameParent.Error())
	assert.Contains(t, string(errJSON)+"\n", string(body))

	// Send a shutdown request
	//a.Stop()
}

func TestAPIJobCreateUpdateParentJob_NoParent(t *testing.T) {
	port := "8093"
	baseURL := fmt.Sprintf("http://localhost:%s/v1", port)
	a := setupAPITest(port)

	jsonStr := []byte(`{
		"name": "test_job",
		"schedule": "@every 1m",
		"command": "date",
		"owner": "mec",
		"owner_email": "foo@bar.com",
		"disabled": true,
		"parent_job": "parent_test_job"
	}`)

	resp, err := http.Post(baseURL+"/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, 422, resp.StatusCode)
	errJSON, err := json.Marshal(ErrParentJobNotFound.Error())
	assert.Contains(t, string(errJSON)+"\n", string(body))
}

func TestAPIJobCreateUpdateValidationBadName(t *testing.T) {
	port := "8094"
	baseURL := fmt.Sprintf("http://localhost:%s/v1", port)
	dir, a := setupAPITest(t, port)
	defer os.RemoveAll(dir)
	defer a.Stop()

	jsonStr := []byte(`{
		"name": "BAD JOB NAME!",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "date"},
		"disabled": true
	}`)

	resp, err := http.Post(baseURL+"/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPIJobCreateUpdateValidationValidName(t *testing.T) {
	port := "8095"
	baseURL := fmt.Sprintf("http://localhost:%s/v1", port)
	dir, a := setupAPITest(t, port)
	defer os.RemoveAll(dir)
	defer a.Stop()

	jsonStr := []byte(`{
		"name": "abcdefghijklmnopqrstuvwxyz0123456789-_ßñëäïüøüáéíóýćàèìòùâêîôûæšłç",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "date"},
		"disabled": true
	}`)

	resp, err := http.Post(baseURL+"/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestAPIJobCreateUpdateValidationBadName(t *testing.T) {
	a := setupAPITest(t)

	jsonStr := []byte(`{"name": "BAD JOB NAME!", "schedule": "@every 1m", "executor": "shell", "executor_config": {"command": "date"}, "disabled": true}`)

	resp, err := http.Post("http://localhost:8090/v1/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Send a shutdown request
	a.Stop()
}

func TestAPIJobCreateUpdateValidationValidName(t *testing.T) {
	a := setupAPITest(t)

	jsonStr := []byte(`{"name": "abcdefghijklmnopqrstuvwxyz0123456789-_ßñëäïüøüáéíóýćàèìòùâêîôûæšłç", "schedule": "@every 1m", "executor": "shell", "executor_config": {"command": "date"}, "disabled": true}`)

	resp, err := http.Post("http://localhost:8090/v1/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Send a shutdown request
	a.Stop()
}
