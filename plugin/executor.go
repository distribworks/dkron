package plugin

import (
	"context"

	"github.com/distribworks/dkron/v3/plugin/types"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type StatusHelper interface {
	Update([]byte, bool) (int64, error)
}

// Executor is the interface that we're exposing as a plugin.
type Executor interface {
	Execute(args *types.ExecuteRequest, cb StatusHelper) (*types.ExecuteResponse, error)
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
	types.RegisterExecutorServer(s, ExecutorServer{Impl: p.Executor, broker: broker})
	return nil
}

func (p *ExecutorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ExecutorClient{client: types.NewExecutorClient(c), broker: broker}, nil
}

// Here is the gRPC client that GRPCClient talks to.
type ExecutorClient struct {
	// This is the real implementation
	client types.ExecutorClient
	broker *plugin.GRPCBroker
}

func (m *ExecutorClient) Execute(args *types.ExecuteRequest, cb StatusHelper) (*types.ExecuteResponse, error) {
	// This is where the magic conversion to Proto happens
	statusHelperServer := &GRPCStatusHelperServer{Impl: cb}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		types.RegisterStatusHelperServer(s, statusHelperServer)

		return s
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	args.StatusServer = brokerID
	r, err := m.client.Execute(context.Background(), args)

	s.Stop()
	return r, err
}

// Here is the gRPC server that GRPCClient talks to.
type ExecutorServer struct {
	// This is the real implementation
	Impl   Executor
	broker *plugin.GRPCBroker
}

// Execute is where the magic happens
func (m ExecutorServer) Execute(ctx context.Context, req *types.ExecuteRequest) (*types.ExecuteResponse, error) {
	conn, err := m.broker.Dial(req.StatusServer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	a := &GRPCStatusHelperClient{types.NewStatusHelperClient(conn)}
	return m.Impl.Execute(req, a)
}

// GRPCStatusHelperClient is an implementation of status updates over RPC.
type GRPCStatusHelperClient struct{ client types.StatusHelperClient }

func (m *GRPCStatusHelperClient) Update(b []byte, c bool) (int64, error) {
	resp, err := m.client.Update(context.Background(), &types.StatusUpdateRequest{
		Output: b,
		Error:  c,
	})
	if err != nil {
		return 0, err
	}
	return resp.R, err
}

// GRPCStatusHelperServer is the gRPC server that GRPCClient talks to.
type GRPCStatusHelperServer struct {
	// This is the real implementation
	Impl StatusHelper
}

func (m *GRPCStatusHelperServer) Update(ctx context.Context, req *types.StatusUpdateRequest) (resp *types.StatusUpdateResponse, err error) {
	r, err := m.Impl.Update(req.Output, req.Error)
	if err != nil {
		return nil, err
	}
	return &types.StatusUpdateResponse{R: r}, err
}
