package interceptor

import (
	"context"
	"log"
	"time"

	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestSimInterceptor struct {
	testAgentAdapter adapters.AgentAdapter
}

func NewTestSimInterceptor(testAgentHost string, timeout time.Duration) *TestSimInterceptor {
	agent, err := adapters.NewTestAgentAdapter(testAgentHost, timeout)
	if err != nil {
		log.Fatalf("Failed to connect to Agent service at %s. Error: %v", testAgentHost, err)
	}

	return &TestSimInterceptor{
		testAgentAdapter: agent,
	}
}

func (t *TestSimInterceptor) UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, rpcHandler grpc.UnaryHandler) (any, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnaryServerInterceptor not implemented")
}
