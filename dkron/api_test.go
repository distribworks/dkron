package dkron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAPITest(t *testing.T, port string) (dir string, a *Agent) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()

	c := DefaultConfig()
	c.BindAddr = ip1.String()
	c.HTTPAddr = fmt.Sprintf("127.0.0.1:%s", port)
	c.NodeName = "test"
	c.Server = true
	c.LogLevel = logLevel
	c.BootstrapExpect = 1
	c.DevMode = true
	c.DataDir = dir

	a = NewAgent(c)
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
	dir, _ := setupAPITest(t, port)
	defer os.RemoveAll(dir)

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

	jsonStr1 := []byte(`{
		"name": "test_job",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "test"},
		"disabled": false
	}`)
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
}

func TestAPIJobCreateUpdateParentJob_SameParent(t *testing.T) {
	resp := postJob(t, "8092", []byte(`{
		"name": "test_job",
		"schedule": "@every 1m",
		"command": "date",
		"owner": "mec",
		"owner_email": "foo@bar.com",
		"disabled": true,
		"parent_job": "test_job"
	}`))

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, string(body), ErrSameParent.Error())
}

func TestAPIJobCreateUpdateParentJob_NoParent(t *testing.T) {
	resp := postJob(t, "8093", []byte(`{
		"name": "test_job",
		"schedule": "@every 1m",
		"command": "date",
		"owner": "mec",
		"owner_email": "foo@bar.com",
		"disabled": true,
		"parent_job": "parent_test_job"
	}`))

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	errJSON, _ := json.Marshal(ErrParentJobNotFound.Error())
	assert.Contains(t, string(errJSON)+"\n", string(body))
}

func TestAPIJobCreateUpdateValidationBadName(t *testing.T) {
	resp := postJob(t, "8094", []byte(`{
		"name": "BAD JOB NAME!",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "date"},
		"disabled": true
	}`))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPIJobCreateUpdateValidationValidName(t *testing.T) {
	resp := postJob(t, "8095", []byte(`{
		"name": "abcdefghijklmnopqrstuvwxyz0123456789-_ßñëäïüøüáéíóýćàèìòùâêîôûæšłç",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "date"},
		"disabled": true
	}`))

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestAPIJobCreateUpdateValidationEmptyName(t *testing.T) {
	port := "8101"
	baseURL := fmt.Sprintf("http://localhost:%s/v1", port)
	dir, a := setupAPITest(t, port)
	defer os.RemoveAll(dir)
	defer a.Stop()

	jsonStr := []byte(`{
		"name": "testjob1",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "date"},
		"disabled": true
	}`)

	resp, err := http.Post(baseURL+"/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	require.NoError(t, err, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	jsonStr = []byte(`{
		"name": "",
		"parent_job": "testjob1",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "date"},
		"disabled": true
	}`)

	resp, err = http.Post(baseURL+"/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	require.NoError(t, err, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPIJobCreateUpdateValidationBadSchedule(t *testing.T) {
	resp := postJob(t, "8097", []byte(`{
		"name": "testjob",
		"schedule": "@at badtime",
		"executor": "shell",
		"executor_config": {"command": "date"},
		"disabled": true
	}`))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPIJobCreateUpdateValidationBadConcurrency(t *testing.T) {
	resp := postJob(t, "8098", []byte(`{
		"name": "testjob",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "date"},
		"concurrency": "badvalue",
		"disabled": true
	}`))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPIJobCreateUpdateValidationBadTimezone(t *testing.T) {
	resp := postJob(t, "8099", []byte(`{
		"name": "testjob",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "date"},
		"disabled": true,
		"timezone": "notatimezone"
	}`))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPIJobCreateUpdateValidationBadShellExecutorTimeout(t *testing.T) {
	resp := postJob(t, "8099", []byte(`{
		"name": "testjob",
		"schedule": "@every 1m",
		"executor": "shell",
		"executor_config": {"command": "date", "timeout": "foreverandever"},
		"disabled": true
	}`))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPIGetNonExistentJobReturnsNotFound(t *testing.T) {
	port := "8096"
	baseURL := fmt.Sprintf("http://localhost:%s/v1", port)
	dir, a := setupAPITest(t, port)
	defer os.RemoveAll(dir)
	defer a.Stop()

	resp, _ := http.Get(baseURL + "/jobs/notajob")

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPIJobCreateUpdateJobWithInvalidParentIsNotCreated(t *testing.T) {
	port := "8100"
	baseURL := fmt.Sprintf("http://localhost:%s/v1", port)
	dir, a := setupAPITest(t, port)
	defer os.RemoveAll(dir)
	defer a.Stop()

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
	require.NoError(t, err, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, ErrParentJobNotFound.Error(), string(body))

	resp, err = http.Get(baseURL + "/jobs/test_job")
	require.NoError(t, err, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPIJobRestore(t *testing.T) {
	port := "8109"
	baseURL := fmt.Sprintf("http://localhost:%s/v1/restore", port)
	dir, a := setupAPITest(t, port)
	defer os.RemoveAll(dir)
	defer a.Stop()

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	fileWriter, err := bodyWriter.CreateFormFile("file", "testBackupJobs.json")
	if err != nil {
		t.Fatalf("CreateFormFile error: %s", err)
	}

	file, err := os.Open("../scripts/testBackupJobs.json")
	if err != nil {
		t.Fatalf("open job json file error: %s", err)
	}
	defer file.Close()

	io.Copy(fileWriter, file)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, _ := http.Post(baseURL, contentType, bodyBuffer)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	rs := string(respBody)
	t.Log("restore response: ", rs)
	if strings.Contains(rs, "fail") {
		t.Fatalf("restore json file request error: %s", rs)
	}

}

func TestAPIJobOutputTruncate(t *testing.T) {
	port := "8190"
	baseURL := fmt.Sprintf("http://localhost:%s/v1", port)
	dir, a := setupAPITest(t, port)
	defer os.RemoveAll(dir)

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
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, _ = http.Get(baseURL + "/jobs/test_job/executions")
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, string(body), "[]")

	testExecution1 := &Execution{
		JobName:    "test_job",
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(),
		Success:    true,
		Output:     "test",
		NodeName:   "testNode",
	}
	testExecution2 := &Execution{
		JobName:    "test_job",
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(),
		Success:    true,
		Output:     "test " + strings.Repeat("longer output... ", 100),
		NodeName:   "testNode2",
	}
	_, err = a.Store.SetExecution(testExecution1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = a.Store.SetExecution(testExecution2)
	if err != nil {
		t.Fatal(err)
	}

	// no truncation
	resp, _ = http.Get(baseURL + "/jobs/test_job/executions")
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var executions []apiExecution
	if err := json.Unmarshal(body, &executions); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(executions))
	assert.False(t, executions[0].OutputTruncated)
	assert.Equal(t, 1705, len(executions[0].Output))
	assert.False(t, executions[1].OutputTruncated)
	assert.Equal(t, 4, len(executions[1].Output))

	// truncate limit to 200
	resp, _ = http.Get(baseURL + "/jobs/test_job/executions?output_size_limit=200")
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if err := json.Unmarshal(body, &executions); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(executions))
	assert.True(t, executions[0].OutputTruncated)
	assert.Equal(t, 200, len(executions[0].Output))
	assert.False(t, executions[1].OutputTruncated)
	assert.Equal(t, 4, len(executions[1].Output))

	// test single execution endpoint
	resp, _ = http.Get(baseURL + "/jobs/test_job/executions/" + executions[0].Id)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var execution Execution
	if err := json.Unmarshal(body, &execution); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1705, len(execution.Output))
}

// postJob POSTs the given json to the jobs endpoint and returns the response
func postJob(t *testing.T, port string, jsonStr []byte) *http.Response {
	baseURL := fmt.Sprintf("http://localhost:%s/v1", port)
	dir, a := setupAPITest(t, port)
	defer os.RemoveAll(dir)
	defer a.Stop()

	resp, err := http.Post(baseURL+"/jobs", "encoding/json", bytes.NewBuffer(jsonStr))
	require.NoError(t, err, err)

	return resp
}
