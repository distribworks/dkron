package plugin

import (
	"context"

	"google.golang.org/grpc"

	"github.com/distribworks/dkron/v2/dkron"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
)

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type ExecutorPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Executor dkron.Executor
}

func (p *ExecutorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	dkron.RegisterExecutorServer(s, ExecutorServer{Impl: p.Executor, broker: broker})
	return nil
}

func (p *ExecutorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ExecutorClient{client: dkron.NewExecutorClient(c), broker: broker}, nil
}

// Here is the gRPC client that GRPCClient talks to.
type ExecutorClient struct {
	// This is the real implementation
	client dkron.ExecutorClient
	broker *plugin.GRPCBroker
}

func (m *ExecutorClient) Execute(args *dkron.ExecuteRequest, cb dkron.StatusHelper) (*dkron.ExecuteResponse, error) {
	// This is where the magic conversion to Proto happens
	statusHelperServer := &GRPCStatusHelperServer{Impl: cb}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		dkron.RegisterStatusHelperServer(s, statusHelperServer)

		return s
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	args.StatusServer = brokerID
	r, err := m.client.Execute(context.Background(), args)

	s.Stop()
	return r, err
}

// ExecutorServer is the gRPC server that GRPCClient talks to.
type ExecutorServer struct {
	// This is the real implementation
	Impl   dkron.Executor
	broker *plugin.GRPCBroker
}

// Execute is where the magic happens
func (m ExecutorServer) Execute(ctx context.Context, req *dkron.ExecuteRequest) (*dkron.ExecuteResponse, error) {
	conn, err := m.broker.Dial(req.StatusServer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	a := &GRPCStatusHelperClient{dkron.NewStatusHelperClient(conn)}
	return m.Impl.Execute(req, a)
}

// GRPCStatusHelperClient is an implementation of status updates over RPC.
type GRPCStatusHelperClient struct{ client dkron.StatusHelperClient }

func (m *GRPCStatusHelperClient) Update(a float32, b []byte, c bool) (int64, error) {
	resp, err := m.client.Update(context.Background(), &dkron.StatusUpdateRequest{
		Progress: a,
		Output:   b,
		Error:    c,
	})
	if err != nil {
		log.WithError(err).Info("status.Update client start error")
		return 0, err
	}
	return resp.R, err
}

// GRPCStatusHelperServer is the gRPC server that GRPCClient talks to.
type GRPCStatusHelperServer struct {
	// This is the real implementation
	Impl dkron.StatusHelper
}

func (m *GRPCStatusHelperServer) Update(ctx context.Context, req *dkron.StatusUpdateRequest) (resp *dkron.StatusUpdateResponse, err error) {
	r, err := m.Impl.Update(req.Progress, req.Output, req.Error)
	if err != nil {
		return nil, err
	}
	return &dkron.StatusUpdateResponse{R: r}, err
}
