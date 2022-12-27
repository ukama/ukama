package providers

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestAgentClientProvider creates a local client to interact with
// a remote instance of Test Agent service.
type TestAgentClientProvider interface {
	GetClient() (pb.TestAgentServiceClient, error)
}

type testAgentClientProvider struct {
	testAgentService pb.TestAgentServiceClient
	testAgentHost    string
}

func NewTestAgentClientProvider(orgHost string) TestAgentClientProvider {
	return &testAgentClientProvider{testAgentHost: orgHost}
}

func (t *testAgentClientProvider) GetClient() (pb.TestAgentServiceClient, error) {
	if t.testAgentService == nil {
		var conn *grpc.ClientConn

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Infoln("Connecting to Org service ", t.testAgentHost)

		conn, err := grpc.DialContext(ctx, t.testAgentHost, grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Errorf("Failed to connect to Org service %s. Error: %v", t.testAgentHost, err)

			return nil, fmt.Errorf("failed to connect to remote org service: %w", err)
		}

		t.testAgentService = pb.NewTestAgentServiceClient(conn)
	}

	return t.testAgentService, nil
}
