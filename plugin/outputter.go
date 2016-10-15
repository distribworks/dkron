package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/victorcoder/dkron/dkron"
)

type OutputPlugin struct {
	Outputter dkron.Outputter
}

func (p *OutputPlugin) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &OutputterServer{Broker: b, Impl: p.Outputter}, nil
}

func (p *OutputPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &Outputter{Broker: b, Client: c}, nil
}

// Here is an implementation that talks over RPC
type Outputter struct {
	Broker *plugin.MuxBroker
	Client *rpc.Client
}

func (o *Outputter) Output(execution *dkron.Execution) []byte {
	var resp []byte
	err := o.Client.Call("Plugin.Output", execution, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that Outputter talks to, conforming to
// the requirements of net/rpc
type OutputterServer struct {
	// This is the real implementation
	Broker *plugin.MuxBroker
	Impl   dkron.Outputter
}

func (s *OutputterServer) Output(execution *dkron.Execution, resp *[]byte) error {
	*resp = s.Impl.Output(execution)
	return nil
}
