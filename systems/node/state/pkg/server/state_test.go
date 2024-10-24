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
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/state/mocks"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStateServer_AddNodeState(t *testing.T) {
	mockStateRepo := &mocks.StateRepo{}
	mockMsgBusClient := &mbmocks.MsgBusServiceClient{}

	stateServer := NewStateServer("test-org", "test-org-id", mockStateRepo, mockMsgBusClient)
    nodeId:=ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String()
	testCases := []struct {
		name    string
		req     *pb.AddStateRequest
		mockFn  func()
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Valid request",
			req: &pb.AddStateRequest{
				NodeId:       nodeId,
				CurrentState: cpb.NodeState_Unknown,
				SubState:     []string{"sub-state-1"},
				Events:       []string{"event-1"},
				NodeType:     "test-type",
				NodeIp:       "1.2.3.4",
				NodePort:     80,
				MeshIp:       "5.6.7.8",
				MeshPort:     8080,
				MeshHostName: "mesh-host",
			},
			mockFn: func() {
				mockStateRepo.On("GetLatestState", nodeId).Return(nil, nil)
				mockStateRepo.On("AddState", mock.Anything, mock.Anything).Return(nil)
			
			},
			wantErr: false,	},
		
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFn()
			resp, err := stateServer.AddNodeState(context.Background(), tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tc.errCode, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Id)
			}
			mockStateRepo.AssertExpectations(t)
		})
	}
}

func TestStateServer_GetLatestState(t *testing.T) {
	mockStateRepo := &mocks.StateRepo{}
	mockMsgBusClient := &mbmocks.MsgBusServiceClient{}

	stateServer := NewStateServer("test-org", "test-org-id", mockStateRepo, mockMsgBusClient)
	nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String()

	testCases := []struct {
		name    string
		req     *pb.GetLatestStateRequest
		mockFn  func()
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Valid request",
			req: &pb.GetLatestStateRequest{
				NodeId: nodeId,
			},
			mockFn: func() {
				mockStateRepo.On("GetLatestState", nodeId).Return(&db.State{
					Id:           uuid.NewV4(),
					NodeId:       nodeId,
					CurrentState: cpb.NodeState_Unknown,
					SubState:     []string{"sub-state-1"},
					Events:       []string{"event-1"},
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}, nil)

			},
			wantErr: false,
		},
		
	
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFn()
			resp, err := stateServer.GetLatestState(context.Background(), tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tc.errCode, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			mockStateRepo.AssertExpectations(t)
		})
	}
}


func TestStateServer_GetStates(t *testing.T) {

	mockStateRepo := &mocks.StateRepo{}
	mockMsgBusClient := &mbmocks.MsgBusServiceClient{}

	stateServer := NewStateServer("test-org", "test-org-id", mockStateRepo, mockMsgBusClient)
	nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String()

	testCases := []struct {
		name    string
		req     *pb.GetStatesRequest
		mockFn  func()
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Valid request",
			req: &pb.GetStatesRequest{
				NodeId: nodeId,
			},
			mockFn: func() {
				mockStateRepo.On("GetStateHistory", nodeId).Return([]db.State{
					{
						Id:           uuid.NewV4(),
						NodeId:       nodeId,
						CurrentState: cpb.NodeState_Unknown,
						SubState:     []string{"sub-state-1"},
						Events:       []string{"event-1"},
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					},
				}, nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFn()
			resp, err := stateServer.GetStates(context.Background(), tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tc.errCode, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			mockStateRepo.AssertExpectations(t)
		})
	}
}
