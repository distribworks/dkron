package plugin

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/victorcoder/dkron/dkron"
	"github.com/victorcoder/dkron/proto"
)

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type ExecutorPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Executor dkron.Executor
}

func (p *ExecutorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterExecutorServer(s, ExecutorServer{Impl: p.Executor})
	return nil
}

func (p *ExecutorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ExecutorClient{client: proto.NewExecutorClient(c)}, nil
}

// Here is the gRPC client that GRPCClient talks to.
type ExecutorClient struct {
	// This is the real implementation
	client proto.ExecutorClient
}

func (m *ExecutorClient) Execute(args *proto.ExecutorArgs) error {
	// This is where the magic conversion to Proto happens
	a := &proto.ExecutorArgs{
		Execution: &proto.Execution{},
	}
	_, err := m.client.Execute(context.Background(), &proto.ExecuteRequest{
		ExecutorArgs: a,
	})
	return err
}

// Here is the gRPC server that GRPCClient talks to.
type ExecutorServer struct {
	// This is the real implementation
	Impl dkron.Executor
}

func (m ExecutorServer) Execute(ctx context.Context, req *proto.ExecuteRequest) (*proto.ExecuteResponse, error) {
	// This is where the magic conversion to native dkron happens
	args := &dkron.ExecutorArgs{
		Execution: dkron.Execution{},
	}
	err := m.Impl.Execute(args)
	return &proto.ExecuteResponse{err.Error()}, err
}
