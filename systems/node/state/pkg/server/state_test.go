package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			CurrentState: pb.NodeStateEnum_STATE_ACTIVE,
			Connectivity: pb.Connectivity_CONNECTIVITY_ONLINE,
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
		CurrentState:    db.StateActive,
		Connectivity:    db.Online,
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

func TestStateServer_Update(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	server := NewstateServer("testOrg", mockRepo, nil, false)

	nodeId := "uk-sa2433-hnode-v0-000c"
	req := &pb.UpdateStateRequest{
		State: &pb.State{
			NodeId:       nodeId,
			CurrentState: pb.NodeStateEnum_STATE_MAINTENANCE,
			Connectivity: pb.Connectivity_CONNECTIVITY_ONLINE,
			Type:         "updatedType",
			Version:      "2.0",
		},
	}

	existingState := &db.State{
		Id:              uuid.NewV4(),
		NodeId:          nodeId,
		CurrentState:    db.StateActive,
		Connectivity:    db.Online,
		LastHeartbeat:   time.Now().Add(-1 * time.Hour),
		LastStateChange: time.Now().Add(-1 * time.Hour),
		Type:            "testType",
		Version:         "1.0",
	}

	mockRepo.On("GetByNodeId", mock.AnythingOfType("ukama.NodeID")).
		Return(existingState, nil)

	mockRepo.On("Update", mock.AnythingOfType("*db.State")).
		Return(nil)

	resp, err := server.Update(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, nodeId, resp.State.NodeId)
	assert.Equal(t, pb.NodeStateEnum_STATE_MAINTENANCE, resp.State.CurrentState)
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

func TestStateServer_ListAll(t *testing.T) {
	mockRepo := mocks.NewStateRepo(t)
	server := NewstateServer("testOrg", mockRepo, nil, false)

	req := &pb.ListAllRequest{}

	expectedStates := []db.State{
		{
			Id:              uuid.NewV4(),
			NodeId:          "UK-SA2433-hnode-V0-000C",
			CurrentState:    db.StateActive,
			Connectivity:    db.Online,
			LastHeartbeat:   time.Now(),
			LastStateChange: time.Now(),
			Type:            "testType1",
			Version:         "1.0",
		},
		{
			Id:              uuid.NewV4(),
			NodeId:          "UK-SA2433-hnode-V0-000D",
			CurrentState:    db.StateMaintenance,
			Connectivity:    db.Offline,
			LastHeartbeat:   time.Now(),
			LastStateChange: time.Now(),
			Type:            "testType2",
			Version:         "2.0",
		},
	}

	mockRepo.On("ListAll").
		Return(expectedStates, nil)

	resp, err := server.ListAll(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.States, 2)
	assert.Equal(t, expectedStates[0].NodeId, resp.States[0].NodeId)
	assert.Equal(t, expectedStates[1].NodeId, resp.States[1].NodeId)
	mockRepo.AssertExpectations(t)
}
