package dcron

import (
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/serf/client"
)

type serfManager struct {
	*client.RPCClient
	Agent  *exec.Cmd
	config *Config
}

// var serf = &serfManager{}

func NewSerfManager(conf *Config) *serfManager {
	return &serfManager{config: conf}
}

// It start the local serf agent waits until it's started to connect the RPC client.
func (sm *serfManager) Start() {
	discover := ""
	if sm.config.Discover != "" {
		discover = " -discover=" + sm.config.Discover
	}
	bind := "-bind=" + sm.config.BindAddr
	rpc_addr := "-rpc-addr=" + sm.config.RPCAddr
	node := "-node=" + sm.config.NodeName

	serfArgs := []string{discover, node, rpc_addr, bind, "-config-file=config/dcron.json"}

	log.Debug("./bin/serf agent " + strings.Join(serfArgs, " "))
	agent, err := spawnProc("./bin/serf agent " + strings.Join(serfArgs, " "))
	if err != nil {
		log.Error(err)
	}

	sm.Agent = agent

	serfConfig := &client.Config{Addr: sm.config.RPCAddr}
	sc, err := client.ClientFromConfig(serfConfig)
	// wait for serf
	for i := 0; err != nil && i < 5; i = i + 1 {
		log.Debug(err)
		time.Sleep(1 * time.Second)
		sc, err = client.ClientFromConfig(serfConfig)
		log.Debugf("Connect to serf agent retry: %d", i)
	}
	if err != nil {
		log.Fatal("Error connecting to serf instance", err)
	}

	sm.RPCClient = sc
}

func (sm *serfManager) Terminate() {
	sm.Close()
	sm.Agent.Process.Signal(syscall.SIGKILL)
}

type Event struct {
	Event   string
	ID      int
	LTime   uint64
	Name    string
	Payload []byte
}
