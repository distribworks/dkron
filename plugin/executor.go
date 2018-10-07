package plugin

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/victorcoder/dkron/dkron"
)

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type ExecutorPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Executor dkron.Executor
}

func (p *ExecutorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	dkron.RegisterExecutorServer(s, ExecutorServer{Impl: p.Executor})
	return nil
}

func (p *ExecutorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ExecutorClient{client: dkron.NewExecutorClient(c)}, nil
}

// Here is the gRPC client that GRPCClient talks to.
type ExecutorClient struct {
	// This is the real implementation
	client dkron.ExecutorClient
}

func (m *ExecutorClient) Execute(args *dkron.ExecuteRequest) (*dkron.ExecuteResponse, error) {
	// This is where the magic conversion to Proto happens
	return m.client.Execute(context.Background(), args)
}

// Here is the gRPC server that GRPCClient talks to.
type ExecutorServer struct {
	// This is the real implementation
	Impl dkron.Executor
}

// Execute is where the magic happens
func (m ExecutorServer) Execute(ctx context.Context, req *dkron.ExecuteRequest) (*dkron.ExecuteResponse, error) {
	return m.Impl.Execute(req)
}
