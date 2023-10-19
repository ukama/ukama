package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/health/mocks"
	pb "github.com/ukama/ukama/systems/node/health/pb/gen"
	"github.com/ukama/ukama/systems/node/health/pkg/db"
)

const testOrgName = "test-org"

var orgId = uuid.NewV4()
var testNode = ukama.NewVirtualNodeId("HomeNode")

func TestHealthServer_GetRunningApps(t *testing.T) {
	// Arrange
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	hRepo := &mocks.HealthRepo{}

	health := db.Health{
		Id:        uuid.NewV4(),
		NodeID:    testNode.StringLowercase(),
		Timestamp: "test",
		System: []db.System{
			{
				Id:    uuid.NewV4(),
				Name:  "test",
				Value: "test",
			},
		},
		Capps: []db.Capp{
			{
				Id:     uuid.NewV4(),
				Name:   "test",
				Tag:    "test",
				Status: db.Status(1),
			},
		},
	}

	hRepo.On("GetRunningApps", health.NodeID).Return(&health, nil).Once()

	s := NewHealthServer(testOrgName,hRepo,msgclientRepo, false )

	// Act
	resp, err := s.GetRunningApps(context.TODO(), &pb.GetRunningAppsRequest{
		NodeId: health.NodeID,
	})

	// Assert
	msgclientRepo.AssertExpectations(t)
	if assert.NoError(t, err) {
		assert.Equal(t, health.NodeID, resp.RunningApps.NodeId)
	}

	hRepo.AssertExpectations(t)
}
