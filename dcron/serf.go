package dcron

import (
	"fmt"
	"github.com/hashicorp/serf/client"
	"os"
	"os/exec"
	"time"
)

var serf *client.RPCClient

// spawn command that specified as proc.
func spawnProc(proc string) error {
	cs := []string{"/bin/bash", "-c", proc}
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	cmd.Env = append(os.Environ())

	fmt.Fprintf(log.Writer(), "Starting %s\n", proc)
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(log.Writer(), "Failed to start %s: %s\n", proc, err)
		return err
	}
	return cmd.Wait()
}

func InitSerfAgent() {
	go spawnProc("./bin/serf agent -config-file=config/serf.json")
	serf = initSerfClient()
	defer serf.Close()
	e := make(chan map[string]interface{}, 10)

	sh, err := serf.Stream("*", e)
	if err != nil {
		log.Error(err)
	}
	defer serf.Stop(sh)

	for {
		select {
		case <-e:
			log.Debugf("Receiving event: %v", e)
		}
	}
}

func initSerfClient() *client.RPCClient {
	serfClient, err := client.NewRPCClient("127.0.0.1:7373")
	// wait for serf
	for i := 0; err != nil && i < 5; i = i + 1 {
		log.Debug(err)
		time.Sleep(1 * time.Second)
		serfClient, err = client.NewRPCClient("127.0.0.1:7373")
		log.Debugf("Connect to serf agent retry: %d", i)
	}
	if err != nil {
		log.Error("Error connecting to serf instance", err)
		return nil
	}
	return serfClient
}

// func eventRouter() {
// 	e := make(chan<- map[string]interface{}, 1)
//
// 	sh, err := serf.Stream("*", e)
//
//
// }
