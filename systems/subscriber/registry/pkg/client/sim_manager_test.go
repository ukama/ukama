package client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/subscriber/registry/pkg/client"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"google.golang.org/grpc"
)

type MockSimManagerServiceClient struct {
	mock.Mock
}

func (m *MockSimManagerServiceClient) SomeMethod(ctx context.Context, in *pb.GetSimRequest, opts ...grpc.CallOption) (*pb.GetSimResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.GetSimResponse), args.Error(1)
}

func TestSimManagerClientProvider(t *testing.T) {
	provider := client.NewSimManagerClientProvider("localhost:8080")

	svc, err := provider.GetSimManagerService()
	assert.NoError(t, err)
	assert.NotNil(t, svc)


	mockSvc := new(MockSimManagerServiceClient)
	mockSvc.On("GetSimResponse", mock.Anything, mock.Anything).Return(&pb.GetSimResponse{}, nil)

}