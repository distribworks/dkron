package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/victorcoder/dkron/dkron"
)

type ExecutionProcessorPlugin struct {
	Processor dkron.ExecutionProcessor
}

func (p *ExecutionProcessorPlugin) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &ExecutionProcessorServer{Broker: b, Impl: p.Processor}, nil
}

func (p *ExecutionProcessorPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &Outputter{Broker: b, Client: c}, nil
}

// Here is an implementation that talks over RPC
type ExecutionProcessor struct {
	Broker *plugin.MuxBroker
	Client *rpc.Client
}

func (e *ExecutionProcessor) Process(execution *dkron.Execution) *dkron.Execution {
	var resp dkron.Execution
	err := e.Client.Call("Plugin.Process", execution, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that Outputter talks to, conforming to
// the requirements of net/rpc
type ExecutionProcessorServer struct {
	// This is the real implementation
	Broker    *plugin.MuxBroker
	Processor dkron.ExecutionProcessor
}

func (e *ExecutionProcessorServer) Process(execution *dkron.Execution, resp *dkron.ExecutionProcessor) error {
	*resp = s.Processor.Process(execution)
	return nil
}
