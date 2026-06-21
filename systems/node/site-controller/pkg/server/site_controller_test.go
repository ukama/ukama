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
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	registry "github.com/ukama/ukama/systems/common/rest/client/registry"
	contpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	contmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	testOrgName  = "test-org"
	testTowerID  = "uk-983794-tnode-78-7830"
	testAmpID    = "uk-983794-anode-78-7830"
)

type fakeControllerProvider struct {
	client contpb.ControllerServiceClient
}

func (f *fakeControllerProvider) GetClient() (contpb.ControllerServiceClient, error) {
	return f.client, nil
}

func newTestServer(nodeClient *mbmocks.NodeClient, controllerClient contpb.ControllerServiceClient) *SiteControllerServer {
	return NewSiteControllerServer(testOrgName, nil, nil, nodeClient, nil, nil, &fakeControllerProvider{client: controllerClient}, nil)
}

func TestSetService_ForwardsToTowerNode(t *testing.T) {
	siteID := "11111111-1111-1111-1111-111111111111"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}

	nodeClient.On("GetNodesBySite", siteID).Return(&registry.NodesBySite{
		Nodes: []registry.NodeInfo{{Id: testAmpID}, {Id: testTowerID}},
	}, nil).Once()

	controllerClient.On("ToggleNodeService", mock.Anything, mock.MatchedBy(func(req *contpb.ToggleNodeServiceRequest) bool {
		return req.NodeId == testTowerID && req.State == "on"
	})).Return(&contpb.ToggleNodeServiceResponse{OperationId: "op-1"}, nil).Once()

	s := newTestServer(nodeClient, controllerClient)

	_, err := s.SetService(context.TODO(), &pb.SetServiceRequest{SiteId: siteID, State: "on"})

	assert.NoError(t, err)
	nodeClient.AssertExpectations(t)
	controllerClient.AssertExpectations(t)
}

func TestSetRadio_ForwardsToAmplifierNode(t *testing.T) {
	siteID := "22222222-2222-2222-2222-222222222222"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}

	nodeClient.On("GetNodesBySite", siteID).Return(&registry.NodesBySite{
		Nodes: []registry.NodeInfo{{Id: testTowerID}, {Id: testAmpID}},
	}, nil).Once()

	controllerClient.On("ToggleRfSwitch", mock.Anything, mock.MatchedBy(func(req *contpb.ToggleRfSwitchRequest) bool {
		return req.NodeId == testAmpID && req.State == "off"
	})).Return(&contpb.ToggleRfSwitchResponse{OperationId: "op-2"}, nil).Once()

	s := newTestServer(nodeClient, controllerClient)

	_, err := s.SetRadio(context.TODO(), &pb.SetRadioRequest{SiteId: siteID, State: "off"})

	assert.NoError(t, err)
	nodeClient.AssertExpectations(t)
	controllerClient.AssertExpectations(t)
}

func TestSetService_NoTowerNode_NotFound(t *testing.T) {
	siteID := "33333333-3333-3333-3333-333333333333"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}

	nodeClient.On("GetNodesBySite", siteID).Return(&registry.NodesBySite{
		Nodes: []registry.NodeInfo{{Id: testAmpID}},
	}, nil).Once()

	s := newTestServer(nodeClient, controllerClient)

	_, err := s.SetService(context.TODO(), &pb.SetServiceRequest{SiteId: siteID, State: "on"})

	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
	nodeClient.AssertExpectations(t)
	controllerClient.AssertNotCalled(t, "ToggleNodeService", mock.Anything, mock.Anything)
}

func TestSetService_EmptySite_InvalidArgument(t *testing.T) {
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}

	s := newTestServer(nodeClient, controllerClient)

	_, err := s.SetService(context.TODO(), &pb.SetServiceRequest{SiteId: "", State: "on"})

	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}
