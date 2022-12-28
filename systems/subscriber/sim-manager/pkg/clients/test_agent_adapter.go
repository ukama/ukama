package clients

import (
	"context"
	"log"

	testagentpb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
)

type TestAgentAdapter struct {
	testAgentService TestAgentClientProvider
	testAgentClient  testagentpb.TestAgentServiceClient
}

func NewTestAgentAdapter(testAgentService TestAgentClientProvider) *TestAgentAdapter {
	testAgentClient, err := testAgentService.GetClient()
	if err != nil {
		log.Fatal()
	}

	return &TestAgentAdapter{
		testAgentClient: testAgentClient,
	}
}

func (s *TestAgentAdapter) ActivateSim(ctx context.Context, simID string) error {
	_, err := s.testAgentClient.ActivateSim(ctx, &testagentpb.ActivateSimRequest{SimID: simID})

	return err
}

func (s *TestAgentAdapter) DeactivateSim(ctx context.Context, simID string) error {
	_, err := s.testAgentClient.DeactivateSim(ctx, &testagentpb.DeactivateSimRequest{SimID: simID})

	return err
}
