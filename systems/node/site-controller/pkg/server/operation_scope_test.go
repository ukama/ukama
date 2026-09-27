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
	"github.com/ukama/ukama/systems/common/ukama"
	contpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	contmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const testControllerID = "uk-983794-cnode-78-7830"

func siteLockNodes() *creg.NodesBySite {
	return &creg.NodesBySite{Nodes: []creg.NodeInfo{
		{Id: testTowerID, Type: ukama.NODE_ID_TYPE_TOWERNODE},
		{Id: testControllerID, Type: ukama.NODE_ID_TYPE_CNODE},
		{Id: testAmpID, Type: ukama.NODE_ID_TYPE_AMPNODE},
	}}
}

func TestSiteActionsRejectBusyControllerBeforeDispatch(t *testing.T) {
	for _, action := range []string{"restart", "radio"} {
		t.Run(action, func(t *testing.T) {
			nodes := &mbmocks.NodeClient{}
			nodes.On("GetNodesBySite", "site").Return(siteLockNodes(), nil).Once()
			controller := &contmocks.ControllerServiceClient{}
			manager := &mbmocks.ManagerClient{}
			manager.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
				assert.ElementsMatch(t, []string{"node:"+testTowerID, "node:"+testAmpID, "node:"+testControllerID}, req.ConflictResourceKeys)
				assert.ElementsMatch(t, []string{"node:"+testTowerID, "node:"+testAmpID}, append([]string{req.ResourceKey}, req.AdditionalResourceKeys...))
				return true
			})).Return(nil, status.Error(codes.AlreadyExists, "controller software update active")).Once()
			s := newTestServer(nodes, controller)
			s.opManager = manager
			var err error
			if action == "restart" {
				_, err = s.RestartSite(context.Background(), &pb.RestartSiteRequest{SiteId: "site"})
			} else {
				_, err = s.SetRadio(context.Background(), &pb.SetRadioRequest{SiteId: "site", State: "off"})
			}
			require.Equal(t, codes.AlreadyExists, status.Code(err))
			assert.Empty(t, controller.Calls)
			manager.AssertExpectations(t)
		})
	}
}

func TestRestartSiteDispatchesOnlyReservedTowerAndAmplifier(t *testing.T) {
	nodes := &mbmocks.NodeClient{}
	nodes.On("GetNodesBySite", "site").Return(siteLockNodes(), nil).Once()
	controller := &contmocks.ControllerServiceClient{}
	for _, id := range []string{testTowerID, testAmpID} {
		id := id
		controller.On("RestartNode", mock.MatchedBy(func(ctx context.Context) bool {
			md, _ := metadata.FromOutgoingContext(ctx)
			opID, token, err := copr.ReservationFromContext(metadata.NewIncomingContext(context.Background(), md))
			return err == nil && opID == "node:"+id && token > 0
		}), mock.MatchedBy(func(req *contpb.RestartNodeRequest) bool {
			return req.NodeId == id
		})).Return(&contpb.RestartNodeResponse{OperationId: "node:"+id}, nil).Once()
	}
	s := newTestServer(nodes, controller)
	resp, err := s.RestartSite(context.Background(), &pb.RestartSiteRequest{SiteId: "site"})
	require.NoError(t, err)
	assert.Len(t, resp.OperationIds, 2)
	controller.AssertExpectations(t)
	assert.Len(t, controller.Calls, 2)
}

func TestSiteRestartFailureReleasesOnlyUnsentReservations(t *testing.T) {
	nodes := &mbmocks.NodeClient{}
	nodes.On("GetNodesBySite", "site").Return(siteLockNodes(), nil).Once()
	controller := &contmocks.ControllerServiceClient{}
	controller.On("RestartNode", mock.Anything, mock.MatchedBy(func(req *contpb.RestartNodeRequest) bool {
		return req.NodeId == testTowerID
	})).Return(nil, status.Error(codes.DeadlineExceeded, "unknown dispatch outcome")).Once()
	s := newTestServer(nodes, controller)
	manager := s.opManager.(*mbmocks.ManagerClient)
	manager.On("ForceUnlock", "node:"+testAmpID, "site_controller", mock.Anything).Return(&copr.OperationInfo{}, nil).Once()
	_, err := s.RestartSite(context.Background(), &pb.RestartSiteRequest{SiteId: "site"})
	require.Error(t, err)
	manager.AssertNotCalled(t, "ForceUnlock", "node:"+testTowerID, mock.Anything, mock.Anything)
	manager.AssertExpectations(t)
	controller.AssertExpectations(t)
}
