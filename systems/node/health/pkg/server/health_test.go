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
	id := uuid.NewV4()
	health := db.Health{
		Id:        id,
		NodeId:    testNode.String(),
		TimeStamp: "test",
		System: []db.System{
			{
				Id:    id,
				Name:  "test",
				Value: "test",
			},
		},
		Capps: []db.Capp{
			{
				Id:     id,
				Name:   "test",
				Tag:    "test",
				Status: db.Status(1),
			},
		},
	}

	hRepo.On("GetRunningAppsInfo", testNode).Return(&health, nil).Once()

	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)

	// Act
	resp, err := s.GetRunningApps(context.TODO(), &pb.GetRunningAppsRequest{
		NodeId: testNode.String(),
	})

	// Assert
	msgclientRepo.AssertExpectations(t)
	if assert.NoError(t, err) {
		assert.Equal(t, health.NodeId, resp.RunningApps.NodeId)
	}

	hRepo.AssertExpectations(t)

}

// func TestHealthServer_StoreRunningAppsInfo(t *testing.T) {
// 	// Arrange
// 	msgclientRepo := &mbmocks.MsgBusServiceClient{}

// 	hRepo := &mocks.HealthRepo{}

// 	health := db.Health{
// 		Id:        uuid.NewV4(),
// 		NodeID:    testNode.StringLowercase(),
// 		Timestamp: "test",
// 		System: []db.System{
// 			{
// 				Id:    uuid.NewV4(),
// 				Name:  "test",
// 				Value: "test",
// 			},
// 		},
// 		Capps: []db.Capp{
// 			{
// 				Id:     uuid.NewV4(),
// 				Name:   "test",
// 				Tag:    "test",
// 				Status: db.Status(1),
// 			},
// 		},
// 	}

// 	hRepo.On("StoreRunningAppsInfo", health).Return(nil).Once()

// 	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)

// 	// Act
// 	_, err := s.StoreRunningAppsInfo(context.TODO(), &pb.StoreRunningAppsInfoRequest{
// 		NodeId: health.NodeID,
// 	})

// 	// Assert
// 	msgclientRepo.AssertExpectations(t)
// 	if assert.NoError(t, err) {
// 		//rest is empty obejct
// 		assert.Equal(t, health.NodeID, nil)
// 	}

// 	hRepo.AssertExpectations(t)
// }
