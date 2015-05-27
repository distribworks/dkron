package dcron

import (
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
)

type serfInstance struct {
	// *client.RPCClient
	// Agent  *exec.Cmd
	config *Config
	serf   *serf.Serf
}

func NewSerfManager(config *Config, serfConfig *serf.Config) *serfManager {

	return &serfManager{config: conf}
}

// It start the local serf agent waits until it's started to connect the RPC client.
func (sm *serfInstance) Start() {
}

func (sm *serfInstance) Terminate() {
}

func (sm *serfInstance) SchedulerReloadQuery(leader string) {
}
