package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dktypes "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	grpc "google.golang.org/grpc"
)

type MockedExecutor struct{}

func (m *MockedExecutor) Execute(ctx context.Context, in *dktypes.ExecuteRequest, opts ...grpc.CallOption) (*dktypes.ExecuteResponse, error) {
	resp := &dktypes.ExecuteResponse{}
	return resp, nil
}

func (m *MockedExecutor) ConfigSchema(ctx context.Context, in *dktypes.ConfigSchemaRequest, opts ...grpc.CallOption) (*dktypes.ConfigSchemaResponse, error) {
	return &dktypes.ConfigSchemaResponse{}, nil
}

type MockedStatusHelper struct{}

func (m MockedStatusHelper) Update([]byte, bool) (int64, error) {
	return 0, nil
}

type MockedBroker struct{}

func (m *MockedBroker) AcceptAndServe(id uint32, s func([]grpc.ServerOption) *grpc.Server) {
}

func (m *MockedBroker) NextId() uint32 {
	return 0
}

func TestExecuteDoesNotPanicIfGRPCIsNotInitializedOnTime(t *testing.T) {
	var brokerMock MockedBroker
	var execMock MockedExecutor
	execClient := ExecutorClient{
		client: &execMock,
		broker: &brokerMock,
	}

	var requestStub dktypes.ExecuteRequest
	var statusHelperMock MockedStatusHelper
	assert.NotPanics(t, func() { execClient.Execute(&requestStub, statusHelperMock) })
}

type BlockingSchemaExecutor struct {
	MockedExecutor
}

func (m *BlockingSchemaExecutor) ConfigSchema(ctx context.Context, in *dktypes.ConfigSchemaRequest, opts ...grpc.CallOption) (*dktypes.ConfigSchemaResponse, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

func TestConfigSchemaTimesOut(t *testing.T) {
	var brokerMock MockedBroker
	execClient := ExecutorClient{
		client: &BlockingSchemaExecutor{},
		broker: &brokerMock,
	}

	start := time.Now()
	schema, err := execClient.ConfigSchema()

	require.Error(t, err)
	assert.Empty(t, schema)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Less(t, time.Since(start), 2*time.Second)
}
