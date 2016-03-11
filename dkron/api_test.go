package dkron

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/mitchellh/cli"
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
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a.Run(args)
	}()

	time.Sleep(10 * time.Millisecond)

	return shutdownCh, resultCh
}

func TestAPIJobCreate(t *testing.T) {
	shutdownCh, _ := setupAPITest(t)

	var jsonStr = []byte(`{"name": "test_job", "schedule": "@every 2s", "command": "date", "owner": "mec", "owner_email": "foo@bar.com", "disabled": true}`)
	resp, err := http.Post("http://localhost:8090/v1/jobs/", "encoding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if bytes.Equal(body, jsonStr) {
		t.Fatalf("error saving job: %s", string(body))
	}

	// Send a shutdown request
	shutdownCh <- struct{}{}
}
