package plugin

import (
	"context"

	"github.com/distribworks/dkron/v2/plugin/types"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Executor is the interface that we're exposing as a plugin.
type Executor interface {
	Execute(args *types.ExecuteRequest) (*types.ExecuteResponse, error)
}

// ExecutorPluginConfig is the plugin config
type ExecutorPluginConfig map[string]string

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type ExecutorPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Executor Executor
}

func (p *ExecutorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	types.RegisterExecutorServer(s, ExecutorServer{Impl: p.Executor})
	return nil
}

func (p *ExecutorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ExecutorClient{client: types.NewExecutorClient(c)}, nil
}

// Here is the gRPC client that GRPCClient talks to.
type ExecutorClient struct {
	// This is the real implementation
	client types.ExecutorClient
}

func (m *ExecutorClient) Execute(args *types.ExecuteRequest) (*types.ExecuteResponse, error) {
	// This is where the magic conversion to Proto happens
	r, err := m.client.Execute(context.Background(), args)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
	return r, nil
}

// Here is the gRPC server that GRPCClient talks to.
type ExecutorServer struct {
	// This is the real implementation
	Impl Executor
}

// Execute is where the magic happens
func (m ExecutorServer) Execute(ctx context.Context, req *types.ExecuteRequest) (*types.ExecuteResponse, error) {
	return m.Impl.Execute(req)
}
