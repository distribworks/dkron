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

	time.Sleep(10 * time.Millisecond)

	return shutdownCh, resultCh
}

func TestAPIJobCreateUpdate(t *testing.T) {
	shutdownCh, _ := setupAPITest(t)

	var jsonStr = []byte(`{"name": "test_job", "schedule": "@every 2s", "command": "date", "owner": "mec", "owner_email": "foo@bar.com", "disabled": true}`)

	var origJob Job
	if err := json.Unmarshal(jsonStr, &origJob); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post("http://localhost:8090/v1/jobs/", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	t.Log(body)
	assert.Equal(t, jsonStr, body)

	var jsonStr1 = []byte(`{"name": "test_job", "schedule": "@every 2s", "command": "test"}`)
	resp, err = http.Post("http://localhost:8090/v1/jobs/", "encoding/json", bytes.NewBuffer(jsonStr1))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	var storedJob Job
	if err := json.Unmarshal(body, &storedJob); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, origJob.Disabled, storedJob.Disabled)
	assert.Equal(t, "test", storedJob.Command)

	// Send a shutdown request
	shutdownCh <- struct{}{}
}
