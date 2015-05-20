package dcron

import (
	"fmt"
	"github.com/hashicorp/serf/client"
	serfs "github.com/hashicorp/serf/serf"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type serfManager struct {
	*client.RPCClient
	Agent *exec.Cmd
}

var serf *serfManager

func (sm *serfManager) Terminate() {
	sm.Agent.Process.Signal(syscall.SIGKILL)
}

func NewSerfManager(agent *exec.Cmd) *serfManager {
	serfConfig := &client.Config{Addr: config.GetString("rpc_addr")}
	sc, err := client.ClientFromConfig(serfConfig)
	// wait for serf
	for i := 0; err != nil && i < 5; i = i + 1 {
		log.Debug(err)
		time.Sleep(1 * time.Second)
		sc, err = client.ClientFromConfig(serfConfig)
		log.Debugf("Connect to serf agent retry: %d", i)
	}
	if err != nil {
		log.Error("Error connecting to serf instance", err)
		return nil
	}
	return &serfManager{sc, agent}
}

// spawn command that specified as proc.
func spawnProc(proc string) (*exec.Cmd, error) {
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
		return nil, err
	}
	return cmd, nil
}

func initSerf() {
	discover := ""
	if config.GetString("discover") != "" {
		discover = " -discover=" + config.GetString("discover")
	}
	serfArgs := []string{discover, "-rpc-addr=" + config.GetString("rpc_addr"), "-config-file=config/dcron.json"}
	agent, err := spawnProc("./bin/serf agent" + strings.Join(serfArgs, " "))

	serf = NewSerfManager(agent)
	defer serf.Close()

	ch := make(chan map[string]interface{}, 1)

	sh, err := serf.Stream("*", ch)
	if err != nil {
		log.Error(err)
	}
	defer serf.Stop(sh)

	for {
		select {
		case event := <-ch:
			for key, val := range event {
				switch ev := val.(type) {
				case serfs.MemberEvent:
					log.Debug(ev)
				default:
					log.Debugf("Receiving event: %s => %v of type %T", key, val, val)
				}
			}
			if event["Event"] == "query" {
				log.Debug(string(event["Payload"].([]byte)))
				serf.Respond(uint64(event["ID"].(int64)), []byte("Peetttee"))
			}
		}
	}
}

type Event struct {
	Event   string
	ID      int
	LTime   uint64
	Name    string
	Payload []byte
}
