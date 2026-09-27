/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/node/controller/mocks"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestToggleServiceRejectsBusySiblingBeforeDispatch(t *testing.T) {
	id := "uk-983794-tnode-78-7830"
	nodes := &mbmocks.NodeClient{}
	nodes.On("Get", id).Return(&creg.NodeInfo{
		Id: id, Site: creg.NodeSiteInfo{SiteId: "site"},
		Status: creg.NodeStatusInfo{Connectivity: "Online"},
	}, nil).Twice()
	nodes.On("GetNodesBySite", "site").Return(&creg.NodesBySite{Nodes: []creg.NodeInfo{
		{Id: id}, {Id: "cnode"}, {Id: "anode"},
	}}, nil).Once()
	manager := &mbmocks.ManagerClient{}
	manager.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		assert.ElementsMatch(t, []string{"node:"+id, "node:cnode", "node:anode"}, req.ConflictResourceKeys)
		return req.Type == "ToggleService" && req.ResourceKey == "node:"+id
	})).Return(nil, status.Error(codes.AlreadyExists, "cnode update active")).Once()
	monitor := &mocks.OperationMonitor{}
	bus := &mbmocks.MsgBusServiceClient{}
	s := NewControllerServer(testOrgName, nil, bus, nil, nil, nodes, manager, monitor, 120, 120, false)
	_, err := s.ToggleService(context.Background(), &pb.ToggleServiceRequest{NodeId: id, State: "off"})
	require.Equal(t, codes.AlreadyExists, status.Code(err))
	assert.Empty(t, monitor.Calls)
	assert.Empty(t, bus.Calls)
	nodes.AssertExpectations(t)
	manager.AssertExpectations(t)
}

func TestReservedOperationIsValidatedAndClaimedOnce(t *testing.T) {
	for _, scenario := range []string{"valid", "wrong node", "wrong action", "already running", "claim lost"} {
		t.Run(scenario, func(t *testing.T) {
			op := &copr.OperationInfo{Id: "reservation", Type: "RestartNode", System: "node",
				ResourceKey: "node:t", RequestedBy: "site_controller", Status: copr.StatusPending, FencingToken: 42}
			ctx := copr.WithReservation(context.Background(), op)
			md, _ := metadata.FromOutgoingContext(ctx)
			ctx = metadata.NewIncomingContext(context.Background(), md)
			switch scenario {
			case "wrong node": op.ResourceKey = "node:c"
			case "wrong action": op.Type = "ToggleRadio"
			case "already running": op.Status = copr.StatusRunning
			}
			manager := &mbmocks.ManagerClient{}
			manager.On("Get", op.Id).Return(op, nil).Once()
			if scenario == "valid" {
				running := *op
				running.Status = copr.StatusRunning
				manager.On("MarkRunning", op.Id, op.FencingToken).Return(&running, nil).Once()
			} else if scenario == "claim lost" {
				manager.On("MarkRunning", op.Id, op.FencingToken).Return(nil, status.Error(codes.FailedPrecondition, "already claimed")).Once()
			}
			s := &ControllerServer{opManager: manager}
			got, err := s.acquireOperation(ctx, "RestartNode", "node:t")
			if scenario == "valid" {
				require.NoError(t, err)
				assert.Equal(t, copr.StatusRunning, got.Status)
			} else {
				require.Error(t, err)
			}
			manager.AssertNotCalled(t, "Start", mock.Anything)
			manager.AssertNotCalled(t, "ForceUnlock", mock.Anything, mock.Anything, mock.Anything)
			manager.AssertExpectations(t)
		})
	}
}
