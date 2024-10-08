package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	dktypes "github.com/distribworks/dkron/v4/types"
	grpc "google.golang.org/grpc"
)

type MockedExecutor struct{}

func (m *MockedExecutor) Execute(ctx context.Context, in *dktypes.ExecuteRequest, opts ...grpc.CallOption) (*dktypes.ExecuteResponse, error) {
	resp := &dktypes.ExecuteResponse{}
	return resp, nil
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
