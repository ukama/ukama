/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	mocks "github.com/ukama/ukama/systems/node/state/mocks"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
)

var nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

func TestStateServer_Create(t *testing.T) {
	mockRepo := new(mocks.StateRepo)
	server := NewstateServer("testOrg", mockRepo, nil, false)

	req := &pb.CreateStateRequest{
		State: &pb.State{
			NodeId:  nodeId.String(),
			State:   pb.NodeStateEnum_STATE_CONFIGURE,
			Type:    "test",
			Version: "1.0",
		},
	}

	mockRepo.On("Create", mock.AnythingOfType("*db.State"), mock.Anything).Return(nil)

	resp, err := server.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, pb.NodeStateEnum_STATE_CONFIGURE, resp.State.State)

	mockRepo.AssertExpectations(t)
}

func TestStateServer_GetByNodeId(t *testing.T) {
	mockRepo := new(mocks.StateRepo)
	server := NewstateServer("testOrg", mockRepo, nil, false)

	req := &pb.GetByNodeIdRequest{
		NodeId: nodeId.String(),
	}

	mockState := &db.State{
		Id:              uuid.NewV4(),
		NodeId:          nodeId.String(),
		State:           ukama.StateOperational,
		LastHeartbeat:   time.Now(),
		LastStateChange: time.Now(),
		Type:            "test",
		Version:         "1.0",
	}

	mockRepo.On("GetByNodeId", mock.AnythingOfType("ukama.NodeID")).Return(mockState, nil)

	resp, err := server.GetByNodeId(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, pb.NodeStateEnum_STATE_CONFIGURE, resp.State.State)

	mockRepo.AssertExpectations(t)
}

func TestStateServer_Delete(t *testing.T) {
	mockRepo := new(mocks.StateRepo)
	server := NewstateServer("testOrg", mockRepo, nil, false)

	req := &pb.DeleteStateRequest{
		NodeId: nodeId.String(),
	}

	mockRepo.On("Delete", mock.AnythingOfType("ukama.NodeID")).Return(nil)

	resp, err := server.Delete(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockRepo.AssertExpectations(t)
}
