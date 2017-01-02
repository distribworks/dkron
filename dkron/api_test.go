package dkron

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func setupAPITest(t *testing.T) (chan<- struct{}, <-chan int) {
	shutdownCh := make(chan struct{})
	// defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	args := []string{
		"-bind", testutil.GetBindAddr().String(),
		"-http-addr", "127.0.0.1:8090",
		"-node", "test",
		"-server",
		"-log-level", logLevel,
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a.Run(args)
	}()

	for {
		if a.ready {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(1 * time.Second)

	return shutdownCh, resultCh
}

func TestAPIJobCreateUpdate(t *testing.T) {
	shutdownCh, _ := setupAPITest(t)

	jsonStr := []byte(`{"name": "test_job", "schedule": "@every 2s", "command": "date", "owner": "mec", "owner_email": "foo@bar.com", "disabled": true}`)

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

	jsonStr1 := []byte(`{"name": "test_job", "schedule": "@every 2s", "command": "test", "disabled": false}`)
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
	assert.NotEqual(t, origJob.Command, overwriteJob.Command)
	assert.Equal(t, "test", overwriteJob.Command)

	// Send a shutdown request
	shutdownCh <- struct{}{}
}

func TestAPIJobCreateUpdateParentJob_SameParent(t *testing.T) {
	shutdownCh, _ := setupAPITest(t)

	jsonStr := []byte(`{
		"name": "test_job",
		"schedule": "@every 2s",
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
	errJson, err := json.Marshal(ErrSameParent.Error())
	assert.Equal(t, string(errJson)+"\n", string(body))

	// Send a shutdown request
	shutdownCh <- struct{}{}
}

func TestAPIJobCreateUpdateParentJob_NoParent(t *testing.T) {
	shutdownCh, _ := setupAPITest(t)

	jsonStr := []byte(`{
		"name": "test_job",
		"schedule": "@every 2s",
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
	errJson, err := json.Marshal(ErrParentJobNotFound.Error())
	assert.Equal(t, string(errJson)+"\n", string(body))

	// Send a shutdown request
	shutdownCh <- struct{}{}
}
