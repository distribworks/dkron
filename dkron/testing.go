package dkron

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T, cb func(*Config)) (*Agent, func()) {
	s, c, err := TestServerErr(t, cb)
	require.NoError(t, err, "failed to start test server")
	return s, c
}

func TestServerErr(t *testing.T, cb func(*Config)) (*Agent, func(), error) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)

	aName := "test1"
	ip, returnFn := testutil.TakeIP()
	aAddr := ip.String()
	defer returnFn()

	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	c := DefaultConfig()
	c.BindAddr = aAddr
	//c.StartJoin = []string{a2Addr}
	c.NodeName = aName
	c.Server = true
	//c.LogLevel = logLevel
	c.BootstrapExpect = 3
	c.DevMode = true
	c.DataDir = dir

	agent := NewAgent(c)
	if err := agent.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	return agent, func() {
		ch := make(chan error)
		go func() {
			defer close(ch)

			// Shutdown server
			err := agent.Stop()
			if err != nil {
				ch <- fmt.Errorf("failed to shutdown server: %w", err)
			}
			os.RemoveAll(dir)
		}()

		select {
		case e := <-ch:
			if e != nil {
				t.Fatal(e.Error())
			}
		case <-time.After(1 * time.Minute):
			t.Fatal("timed out while shutting down server")
		}
	}, nil
}
