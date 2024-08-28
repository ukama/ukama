package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/state/mocks"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
)

func TestStateServer_Create(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	server := NewstateServer("testOrg", mockRepo, nil, false)

	req := &pb.CreateStateRequest{
		State: &pb.State{
			NodeId:       "uk-sa2433-hnode-v0-000c",
			State: pb.NodeStateEnum_STATE_CONFIGURE,
			Type:         "testType",
			Version:      "1.0",
		},
	}

	mockRepo.On("Create", mock.AnythingOfType("*db.State"), mock.AnythingOfType("func(*db.State, *gorm.DB) error")).
		Return(nil)

	resp, err := server.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.State.NodeId, resp.State.NodeId)
	mockRepo.AssertExpectations(t)
}

func TestStateServer_GetByNodeId(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	server := NewstateServer("testOrg", mockRepo, nil, false)

	nodeId := "uk-sa2433-hnode-v0-000c"
	req := &pb.GetByNodeIdRequest{
		NodeId: nodeId,
	}

	expectedState := &db.State{
		Id:              uuid.NewV4(),
		NodeId:          nodeId,
		State:    ukama.StateConfigure,
		LastHeartbeat:   time.Now(),
		LastStateChange: time.Now(),
		Type:            "testType",
		Version:         "1.0",
	}

	mockRepo.On("GetByNodeId", mock.AnythingOfType("ukama.NodeID")).
		Return(expectedState, nil)

	resp, err := server.GetByNodeId(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, nodeId, resp.State.NodeId)
	mockRepo.AssertExpectations(t)
}


func TestStateServer_Delete(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	server := NewstateServer("testOrg", mockRepo, nil, false)

	nodeId := "UK-SA2433-hnode-V0-000C"
	req := &pb.DeleteStateRequest{
		NodeId: nodeId,
	}

	mockRepo.On("Delete", mock.AnythingOfType("ukama.NodeID")).
		Return(nil)

	resp, err := server.Delete(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockRepo.AssertExpectations(t)
}
