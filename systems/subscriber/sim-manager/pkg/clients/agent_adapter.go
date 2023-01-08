package clients

import (
	"context"
	"log"
	"time"

	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AgentAdapter interface {
	ActivateSim(context.Context, string) error
	DeactivateSim(context.Context, string) error
	// TerminateSim(context.Context, string) error
	Close()
}

type AgentFactory struct {
	timeout time.Duration
	factory map[sims.SimType]AgentAdapter
}

func NewAgentFactory(testAgentHost string, timeout time.Duration) *AgentFactory {
	var factory = make(map[sims.SimType]AgentAdapter)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	testAgentConn, err := grpc.DialContext(ctx, testAgentHost, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to Test Agent service %s. Error: %v", testAgentHost, err)

	}

	factory[sims.SimTypeInterNone] = &TestAgentAdapter{
		conn:   testAgentConn,
		host:   testAgentHost,
		client: pb.NewTestAgentServiceClient(testAgentConn)}

	return &AgentFactory{
		timeout: timeout,
		factory: factory,
	}
}

func (a *AgentFactory) GetAgentAdapter(simType sims.SimType) (AgentAdapter, bool) {
	agent, ok := a.factory[simType]

	return agent, ok
}

func (a *AgentFactory) Close() {
	for _, adapter := range a.factory {
		adapter.Close()
	}
}
