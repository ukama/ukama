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
	"github.com/ukama/ukama/systems/common/ukama"
	contpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	contmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
	hpb "github.com/ukama/ukama/systems/node/health/pb/gen"
	scmocks "github.com/ukama/ukama/systems/node/site-controller/mocks"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/reconciler"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	testOrgName = "test-org"
	testTowerID = "uk-983794-tnode-78-7830"
	testAmpID   = "uk-983794-anode-78-7830"
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

	controllerClient.On("ToggleService", mock.Anything, mock.MatchedBy(func(req *contpb.ToggleServiceRequest) bool {
		return req.NodeId == testTowerID && req.State == "on"
	})).Return(&contpb.ToggleServiceResponse{OperationId: "op-1"}, nil).Once()

	s := newTestServer(nodeClient, controllerClient)

	_, err := s.SetService(context.TODO(), &pb.SetServiceRequest{SiteId: siteID, State: "on"})

	assert.NoError(t, err)
	nodeClient.AssertExpectations(t)
	controllerClient.AssertExpectations(t)
}

func TestSetRadio_ForwardsToTowerAndAmplifierNodes(t *testing.T) {
	siteID := "22222222-2222-2222-2222-222222222222"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}

	nodeClient.On("GetNodesBySite", siteID).Return(&registry.NodesBySite{
		Nodes: []registry.NodeInfo{
			{Id: testTowerID, Type: ukama.NODE_ID_TYPE_TOWERNODE},
			{Id: testAmpID, Type: ukama.NODE_ID_TYPE_AMPNODE},
		},
	}, nil).Once()

	controllerClient.On("ToggleRadio", mock.Anything, mock.MatchedBy(func(req *contpb.ToggleRadioRequest) bool {
		return req.NodeId == testTowerID && req.State == "off"
	})).Return(&contpb.ToggleRadioResponse{OperationId: "op-t"}, nil).Once()

	controllerClient.On("ToggleRadio", mock.Anything, mock.MatchedBy(func(req *contpb.ToggleRadioRequest) bool {
		return req.NodeId == testAmpID && req.State == "off"
	})).Return(&contpb.ToggleRadioResponse{OperationId: "op-a"}, nil).Once()

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

func TestRecordObservedState_WritesStateNotIntent(t *testing.T) {
	siteID := "22222222-2222-2222-2222-222222222222"

	states := &scmocks.StateRepo{}
	intents := &scmocks.IntentRepo{}

	states.On("Upsert", mock.MatchedBy(func(s *db.SiteState) bool {
		return s.SiteID == siteID &&
			s.RadioState == "on" &&
			s.ServiceState == "" &&
			s.Reason == reasonHealthReport
	})).Return(nil).Once()

	c := &SiteControllerEventServer{states: states, intents: intents}

	err := c.recordObservedState(siteID, &db.SiteState{
		RadioState: "on",
		Reason:     reasonHealthReport,
	})

	assert.NoError(t, err)
	states.AssertExpectations(t)
	intents.AssertNotCalled(t, "Upsert", mock.Anything)
}

func TestRecordObservedState_LeavesOtherAxisUntouched(t *testing.T) {
	siteID := "33333333-3333-3333-3333-333333333333"

	states := &scmocks.StateRepo{}
	states.On("Upsert", mock.MatchedBy(func(s *db.SiteState) bool {
		return s.ServiceState == "" && s.PowerState == "" && s.AccessState == ""
	})).Return(nil).Once()

	c := &SiteControllerEventServer{states: states}

	err := c.recordObservedState(siteID, &db.SiteState{RadioState: "off"})

	assert.NoError(t, err)
	states.AssertExpectations(t)
}

func TestObservedServiceState(t *testing.T) {
	cases := []struct {
		name string
		apps []*hpb.App
		want string
	}{
		{
			name: "pcrf active means service is running",
			apps: []*hpb.App{{Name: "pcrf", Status: "Active"}},
			want: reconciler.StateRunning,
		},
		{
			name: "pcrf inactive means service is off",
			apps: []*hpb.App{{Name: "pcrf", Status: "Inactive"}},
			want: reconciler.StateOff,
		},
		{
			name: "pcrf missing means service is off",
			apps: []*hpb.App{{Name: "lookout", Status: "Active"}},
			want: reconciler.StateOff,
		},
		{
			name: "no apps reported means service is off",
			apps: nil,
			want: reconciler.StateOff,
		},
		{
			name: "another active app does not imply service",
			apps: []*hpb.App{{Name: "meshd", Status: "Active"}, {Name: "pcrf", Status: "Stopped"}},
			want: reconciler.StateOff,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, observedServiceState(tc.apps))
		})
	}
}

func TestObservedServiceState_UsesReconcilerVocabulary(t *testing.T) {
	running := observedServiceState([]*hpb.App{{Name: "pcrf", Status: "Active"}})

	assert.Equal(t, reconciler.StateRunning, running,
		"observed service must use the value the reconciler expects for a desired service of on")
}
