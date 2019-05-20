package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/victorcoder/dkron/plugintypes"
)

type ExecutionProcessorPlugin struct {
	Processor plugintypes.ExecutionProcessor
}

func (p *ExecutionProcessorPlugin) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &ExecutionProcessorServer{Broker: b, Processor: p.Processor}, nil
}

func (p *ExecutionProcessorPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ExecutionProcessor{Broker: b, Client: c}, nil
}

// Here is an implementation that talks over RPC
type ExecutionProcessor struct {
	Broker *plugin.MuxBroker
	Client *rpc.Client
}

// The Process method that actually call the plugin Process method.
func (e *ExecutionProcessor) Process(args *plugintypes.ExecutionProcessorArgs) plugintypes.Execution {
	var resp plugintypes.Execution
	err := e.Client.Call("Plugin.Process", args, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that client talks to, conforming to
// the requirements of net/rpc
type ExecutionProcessorServer struct {
	// This is the real implementation
	Broker    *plugin.MuxBroker
	Processor plugintypes.ExecutionProcessor
}

func (e *ExecutionProcessorServer) Process(args *plugintypes.ExecutionProcessorArgs, resp *plugintypes.Execution) error {
	*resp = e.Processor.Process(args)
	return nil
}
