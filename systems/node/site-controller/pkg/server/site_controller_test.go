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
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/operation"
	registry "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	contpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	contmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
	scmocks "github.com/ukama/ukama/systems/node/site-controller/mocks"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	testOrgName = "test-org"
	testTowerID = "uk-983794-tnode-78-7830"
	testAmpID   = "uk-983794-anode-78-7830"
	testCnodeID = "uk-983794-cnode-78-7830"
)

type fakeControllerProvider struct {
	client contpb.ControllerServiceClient
}

func (f *fakeControllerProvider) GetClient() (contpb.ControllerServiceClient, error) {
	return f.client, nil
}

var online = registry.NodeStatusInfo{Connectivity: ukama.NodeConnectivityOnline.String()}

func newTestServer(nodeClient *mbmocks.NodeClient, controllerClient contpb.ControllerServiceClient) *SiteControllerServer {
	opManager := &mbmocks.ManagerClient{}
	opManager.On("GetByResource", mock.Anything).Return(nil, nil)
	return newTestServerWithOps(nodeClient, controllerClient, opManager)
}

func newTestServerWithOps(nodeClient *mbmocks.NodeClient, controllerClient contpb.ControllerServiceClient, opManager *mbmocks.ManagerClient) *SiteControllerServer {
	return NewSiteControllerServer(testOrgName, nil, nil, nodeClient, nil, nil, &fakeControllerProvider{client: controllerClient}, opManager, nil)
}

func TestSetService_ForwardsToTowerNode(t *testing.T) {
	siteID := "11111111-1111-1111-1111-111111111111"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}

	nodeClient.On("GetNodesBySite", siteID).Return(&registry.NodesBySite{
		Nodes: []registry.NodeInfo{{Id: testAmpID, Status: online}, {Id: testTowerID, Status: online}},
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
			{Id: testTowerID, Type: ukama.NODE_ID_TYPE_TOWERNODE, Status: online},
			{Id: testAmpID, Type: ukama.NODE_ID_TYPE_AMPNODE, Status: online},
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
		Nodes: []registry.NodeInfo{{Id: testAmpID, Status: online}},
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

func siteWithAllNodes() *registry.NodesBySite {
	return &registry.NodesBySite{
		Nodes: []registry.NodeInfo{
			{Id: testTowerID, Type: ukama.NODE_ID_TYPE_TOWERNODE, Status: online},
			{Id: testAmpID, Type: ukama.NODE_ID_TYPE_AMPNODE, Status: online},
			{Id: testCnodeID, Type: ukama.NODE_ID_TYPE_CNODE, Status: online},
		},
	}
}

func cnodeUpdating() *mbmocks.ManagerClient {
	opManager := &mbmocks.ManagerClient{}
	opManager.On("GetByResource", "node:"+testCnodeID).Return(&operation.OperationInfo{
		Id: "op-update", Type: "UpdateSoftware", Status: operation.StatusRunning, ResourceKey: "node:" + testCnodeID,
	}, nil)
	opManager.On("GetByResource", mock.Anything).Return(nil, nil)
	return opManager
}

func TestRestartSite_RejectedWhileCnodeUpdating(t *testing.T) {
	siteID := "44444444-4444-4444-4444-444444444444"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	nodeClient.On("GetNodesBySite", siteID).Return(siteWithAllNodes(), nil).Once()

	s := newTestServerWithOps(nodeClient, controllerClient, cnodeUpdating())

	_, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: siteID})

	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	assert.Contains(t, err.Error(), testCnodeID)
	controllerClient.AssertNotCalled(t, "RestartNode", mock.Anything, mock.Anything)
}

func TestRestartSite_IdleSiteRestartsTowerAndAmplifierOnly(t *testing.T) {
	siteID := "55555555-5555-5555-5555-555555555555"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	opManager := &mbmocks.ManagerClient{}
	nodeClient.On("GetNodesBySite", siteID).Return(siteWithAllNodes(), nil).Once()
	opManager.On("GetByResource", "node:"+testTowerID).Return(nil, nil).Once()
	opManager.On("GetByResource", "node:"+testAmpID).Return(nil, nil).Once()
	opManager.On("GetByResource", "node:"+testCnodeID).Return(&operation.OperationInfo{
		Id: "op-done", Type: "UpdateSoftware", Status: operation.StatusSuccess,
	}, nil).Once()
	controllerClient.On("RestartNode", mock.Anything, mock.MatchedBy(func(req *contpb.RestartNodeRequest) bool {
		return req.NodeId == testTowerID || req.NodeId == testAmpID
	})).Return(&contpb.RestartNodeResponse{OperationId: "op-r"}, nil).Twice()

	s := newTestServerWithOps(nodeClient, controllerClient, opManager)

	resp, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: siteID})

	assert.NoError(t, err)
	assert.Len(t, resp.OperationIds, 2)
	opManager.AssertExpectations(t)
	controllerClient.AssertExpectations(t)
}

func TestSetRadio_RejectedWhileCnodeUpdating(t *testing.T) {
	siteID := "66666666-6666-6666-6666-666666666666"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	nodeClient.On("GetNodesBySite", siteID).Return(siteWithAllNodes(), nil).Once()

	s := newTestServerWithOps(nodeClient, controllerClient, cnodeUpdating())

	_, err := s.SetRadio(context.TODO(), &pb.SetRadioRequest{SiteId: siteID, State: "off"})

	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	controllerClient.AssertNotCalled(t, "ToggleRadio", mock.Anything, mock.Anything)
}

func TestSetService_RejectedWhileCnodeUpdating(t *testing.T) {
	siteID := "77777777-7777-7777-7777-777777777777"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	nodeClient.On("GetNodesBySite", siteID).Return(siteWithAllNodes(), nil).Once()

	s := newTestServerWithOps(nodeClient, controllerClient, cnodeUpdating())

	_, err := s.SetService(context.TODO(), &pb.SetServiceRequest{SiteId: siteID, State: "off"})

	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	controllerClient.AssertNotCalled(t, "ToggleService", mock.Anything, mock.Anything)
}

func TestSetService_IgnoresBusyAmplifier(t *testing.T) {
	siteID := "88888888-8888-8888-8888-888888888888"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	opManager := &mbmocks.ManagerClient{}
	nodeClient.On("GetNodesBySite", siteID).Return(siteWithAllNodes(), nil).Once()
	opManager.On("GetByResource", "node:"+testTowerID).Return(nil, nil).Once()
	opManager.On("GetByResource", "node:"+testCnodeID).Return(nil, nil).Once()
	controllerClient.On("ToggleService", mock.Anything, mock.Anything).Return(&contpb.ToggleServiceResponse{OperationId: "op-s"}, nil).Once()

	s := newTestServerWithOps(nodeClient, controllerClient, opManager)

	_, err := s.SetService(context.TODO(), &pb.SetServiceRequest{SiteId: siteID, State: "on"})

	assert.NoError(t, err)
	opManager.AssertNotCalled(t, "GetByResource", "node:"+testAmpID)
	controllerClient.AssertExpectations(t)
}

func TestRestartSite_LockReadFailureIsUnavailable(t *testing.T) {
	siteID := "99999999-9999-9999-9999-999999999999"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	opManager := &mbmocks.ManagerClient{}
	nodeClient.On("GetNodesBySite", siteID).Return(siteWithAllNodes(), nil).Once()
	opManager.On("GetByResource", mock.Anything).Return(nil, errors.New("connection refused"))

	s := newTestServerWithOps(nodeClient, controllerClient, opManager)

	_, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: siteID})

	assert.Equal(t, codes.Unavailable, status.Code(err))
	controllerClient.AssertNotCalled(t, "RestartNode", mock.Anything, mock.Anything)
}

func TestRestartSite_NoOperationManagerIsUnavailable(t *testing.T) {
	siteID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	nodeClient.On("GetNodesBySite", siteID).Return(siteWithAllNodes(), nil).Once()

	s := NewSiteControllerServer(testOrgName, nil, nil, nodeClient, nil, nil, &fakeControllerProvider{client: controllerClient}, nil, nil)

	_, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: siteID})

	assert.Equal(t, codes.Unavailable, status.Code(err))
	controllerClient.AssertNotCalled(t, "RestartNode", mock.Anything, mock.Anything)
}

func TestSetRadio_OfflineAmplifierRejectedBeforeAnyToggle(t *testing.T) {
	siteID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	site := siteWithAllNodes()
	site.Nodes[1].Status.Connectivity = ukama.NodeConnectivityOffline.String()
	nodeClient.On("GetNodesBySite", siteID).Return(site, nil).Once()

	s := newTestServer(nodeClient, controllerClient)

	_, err := s.SetRadio(context.TODO(), &pb.SetRadioRequest{SiteId: siteID, State: "off"})

	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	assert.Contains(t, err.Error(), testAmpID)
	controllerClient.AssertNotCalled(t, "ToggleRadio", mock.Anything, mock.Anything)
}

func TestRestartSite_OfflineCnodeRejected(t *testing.T) {
	siteID := "cccccccc-cccc-cccc-cccc-cccccccccccc"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	site := siteWithAllNodes()
	site.Nodes[2].Status.Connectivity = ukama.NodeConnectivityOffline.String()
	nodeClient.On("GetNodesBySite", siteID).Return(site, nil).Once()

	s := newTestServer(nodeClient, controllerClient)

	_, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: siteID})

	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	assert.Contains(t, err.Error(), testCnodeID)
	controllerClient.AssertNotCalled(t, "RestartNode", mock.Anything, mock.Anything)
}

func TestSetRadio_MissingAmplifierNotFound(t *testing.T) {
	siteID := "dddddddd-dddd-dddd-dddd-dddddddddddd"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	nodeClient.On("GetNodesBySite", siteID).Return(&registry.NodesBySite{
		Nodes: []registry.NodeInfo{
			{Id: testTowerID, Type: ukama.NODE_ID_TYPE_TOWERNODE, Status: online},
			{Id: testCnodeID, Type: ukama.NODE_ID_TYPE_CNODE, Status: online},
		},
	}, nil).Once()

	s := newTestServer(nodeClient, controllerClient)

	_, err := s.SetRadio(context.TODO(), &pb.SetRadioRequest{SiteId: siteID, State: "on"})

	assert.Equal(t, codes.NotFound, status.Code(err))
	controllerClient.AssertNotCalled(t, "ToggleRadio", mock.Anything, mock.Anything)
}

func TestRestartSite_EmptyRegistryTypeNeverRestartsCnode(t *testing.T) {
	siteID := "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	nodeClient := &mbmocks.NodeClient{}
	controllerClient := &contmocks.ControllerServiceClient{}
	nodeClient.On("GetNodesBySite", siteID).Return(&registry.NodesBySite{
		Nodes: []registry.NodeInfo{
			{Id: testCnodeID, Status: online},
			{Id: testTowerID, Status: online},
			{Id: testAmpID, Status: online},
		},
	}, nil).Once()
	controllerClient.On("RestartNode", mock.Anything, mock.MatchedBy(func(req *contpb.RestartNodeRequest) bool {
		return req.NodeId == testTowerID || req.NodeId == testAmpID
	})).Return(&contpb.RestartNodeResponse{OperationId: "op-r"}, nil).Twice()

	s := newTestServer(nodeClient, controllerClient)

	resp, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: siteID})

	assert.NoError(t, err)
	assert.Len(t, resp.OperationIds, 2)
	controllerClient.AssertExpectations(t)
	controllerClient.AssertNotCalled(t, "RestartNode", mock.Anything, mock.MatchedBy(func(req *contpb.RestartNodeRequest) bool {
		return req.NodeId == testCnodeID
	}))
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

	err := c.recordObservedState(context.Background(), siteID, &db.SiteState{
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

	err := c.recordObservedState(context.Background(), siteID, &db.SiteState{RadioState: "off"})

	assert.NoError(t, err)
	states.AssertExpectations(t)
}

func TestRecordObservedState_SurvivesMissingReconciler(t *testing.T) {
	siteID := "55555555-5555-5555-5555-555555555555"

	states := &scmocks.StateRepo{}
	states.On("Upsert", mock.Anything).Return(nil).Once()

	c := &SiteControllerEventServer{states: states}

	err := c.recordObservedState(context.Background(), siteID, &db.SiteState{RadioState: "on"})

	assert.NoError(t, err, "recording an observation must not depend on the reconciler")
	states.AssertExpectations(t)
}

func TestRecordObservedState_ReportsStoreFailure(t *testing.T) {
	siteID := "66666666-6666-6666-6666-666666666666"

	states := &scmocks.StateRepo{}
	states.On("Upsert", mock.Anything).Return(assert.AnError).Once()

	c := &SiteControllerEventServer{states: states}

	err := c.recordObservedState(context.Background(), siteID, &db.SiteState{RadioState: "on"})

	assert.Error(t, err, "a failed observation write must surface so the event is retried")
}

func TestNormalizeNodeType(t *testing.T) {
	cases := []struct {
		reported string
		want     string
	}{
		{"tower", ukama.NODE_ID_TYPE_TOWERNODE},
		{"tnode", ukama.NODE_ID_TYPE_TOWERNODE},
		{"Tower", ukama.NODE_ID_TYPE_TOWERNODE},
		{"amplifier", ukama.NODE_ID_TYPE_AMPNODE},
		{"anode", ukama.NODE_ID_TYPE_AMPNODE},
		{"control", ukama.NODE_ID_TYPE_CNODE},
		{"controller", ukama.NODE_ID_TYPE_CNODE},
		{"cnode", ukama.NODE_ID_TYPE_CNODE},
		{"something-else", "something-else"},
	}

	for _, tc := range cases {
		t.Run(tc.reported, func(t *testing.T) {
			assert.Equal(t, tc.want, normalizeNodeType(tc.reported),
				"nodes report long names while the constants are short codes")
		})
	}
}

func TestHandleNodeOnline_IgnoresNonTowerNodes(t *testing.T) {
	states := &scmocks.StateRepo{}

	c := &SiteControllerEventServer{states: states}

	err := c.handleNodeOnline(context.Background(), &epb.NodeOnlineEvent{NodeId: testAmpID})

	assert.NoError(t, err)
	states.AssertNotCalled(t, "Upsert", mock.Anything)
}

func TestHandleNodeOnline_RejectsUnparseableNodeId(t *testing.T) {
	states := &scmocks.StateRepo{}

	c := &SiteControllerEventServer{states: states}

	err := c.handleNodeOnline(context.Background(), &epb.NodeOnlineEvent{NodeId: "not-a-node-id"})

	assert.Error(t, err)
	states.AssertNotCalled(t, "Upsert", mock.Anything)
}

func registryError(code int) error {
	return fmt.Errorf("GetNode failure: %w", fmt.Errorf("rest api GET failure with error: %w",
		&client.ErrorStatus{StatusCode: code, RawError: errors.New("registry error")}))
}

// The event server has no health client here, so reaching ListInterfaces would panic.
func TestHandleHealthReport_IgnoresNodeWithoutSite(t *testing.T) {
	nodeClient := &mbmocks.NodeClient{}
	nodeClient.On("Get", testAmpID).Return(&registry.NodeInfo{Id: testAmpID}, nil).Once()
	nodeClient.On("List", registry.ListNodesRequest{NodeId: testAmpID}).
		Return(&registry.ListNodesResponse{Nodes: []*registry.NodeInfo{{Id: testAmpID}}}, nil).Once()
	states := &scmocks.StateRepo{}
	c := &SiteControllerEventServer{s: newTestServer(nodeClient, nil), states: states}

	err := c.handleHealthReport(context.Background(), &epb.HealthReportEvent{NodeId: testAmpID, NodeType: ukama.NODE_ID_TYPE_AMPNODE})

	assert.NoError(t, err)
	nodeClient.AssertExpectations(t)
	states.AssertNotCalled(t, "Upsert", mock.Anything)
}

func TestHandleHealthReport_IgnoresNodeMissingFromRegistry(t *testing.T) {
	nodeClient := &mbmocks.NodeClient{}
	nodeClient.On("Get", testAmpID).Return(nil, registryError(http.StatusNotFound)).Once()
	c := &SiteControllerEventServer{s: newTestServer(nodeClient, nil)}

	err := c.handleHealthReport(context.Background(), &epb.HealthReportEvent{NodeId: testAmpID, NodeType: ukama.NODE_ID_TYPE_AMPNODE})

	assert.NoError(t, err)
	nodeClient.AssertNotCalled(t, "List", mock.Anything)
}

func TestHandleHealthReport_ReturnsRegistryOutage(t *testing.T) {
	nodeClient := &mbmocks.NodeClient{}
	nodeClient.On("Get", testAmpID).Return(nil, registryError(http.StatusInternalServerError)).Once()
	c := &SiteControllerEventServer{s: newTestServer(nodeClient, nil)}

	err := c.handleHealthReport(context.Background(), &epb.HealthReportEvent{NodeId: testAmpID, NodeType: ukama.NODE_ID_TYPE_AMPNODE})

	assert.Error(t, err)
}

func TestHandleNodeOnline_IgnoresTowerWithoutSite(t *testing.T) {
	for name, resp := range map[string]*registry.ListNodesResponse{
		"no nodes":  {},
		"no site":   {Nodes: []*registry.NodeInfo{{Id: testTowerID}}},
		"nil entry": {Nodes: []*registry.NodeInfo{nil}},
	} {
		t.Run(name, func(t *testing.T) {
			nodeClient := &mbmocks.NodeClient{}
			nodeClient.On("List", registry.ListNodesRequest{NodeId: testTowerID}).Return(resp, nil).Once()
			states := &scmocks.StateRepo{}
			c := &SiteControllerEventServer{s: newTestServer(nodeClient, nil), states: states}

			err := c.handleNodeOnline(context.Background(), &epb.NodeOnlineEvent{NodeId: testTowerID})

			assert.NoError(t, err)
			states.AssertNotCalled(t, "Upsert", mock.Anything)
		})
	}
}

func TestHandleNodeOnline_ReturnsRegistryOutage(t *testing.T) {
	nodeClient := &mbmocks.NodeClient{}
	nodeClient.On("List", registry.ListNodesRequest{NodeId: testTowerID}).Return(nil, registryError(http.StatusServiceUnavailable)).Once()
	states := &scmocks.StateRepo{}
	c := &SiteControllerEventServer{s: newTestServer(nodeClient, nil), states: states}

	err := c.handleNodeOnline(context.Background(), &epb.NodeOnlineEvent{NodeId: testTowerID})

	assert.Error(t, err)
	states.AssertNotCalled(t, "Upsert", mock.Anything)
}

func TestEventNotification_SiteDeleteRemovesSite(t *testing.T) {
	siteID := "77777777-7777-7777-7777-777777777777"

	sites := &scmocks.SiteRepo{}
	sites.On("Delete", siteID).Return(nil).Once()

	c := &SiteControllerEventServer{sites: sites, orgName: testOrgName}

	anyMsg, err := anypb.New(&epb.EventDeleteSite{SiteId: siteID, NetworkId: "net"})
	assert.NoError(t, err)

	_, err = c.EventNotification(context.Background(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, "event.cloud.local.{{ .Org}}.registry.site.site.delete"),
		Msg:        anyMsg,
	})

	assert.NoError(t, err)
	sites.AssertExpectations(t)
}

func TestEventNotification_RejectsMismatchedPayload(t *testing.T) {
	sites := &scmocks.SiteRepo{}
	c := &SiteControllerEventServer{sites: sites, orgName: testOrgName}

	anyMsg, err := anypb.New(&epb.EventAddSite{SiteId: "99999999-9999-9999-9999-999999999999"})
	assert.NoError(t, err)

	_, err = c.EventNotification(context.Background(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, "event.cloud.local.{{ .Org}}.registry.site.site.delete"),
		Msg:        anyMsg,
	})

	assert.Error(t, err)
	sites.AssertNotCalled(t, "Delete", mock.Anything)
}

func TestHandleDeleteSite_ReportsStoreFailure(t *testing.T) {
	sites := &scmocks.SiteRepo{}
	sites.On("Delete", mock.Anything).Return(assert.AnError).Once()

	c := &SiteControllerEventServer{sites: sites}

	err := c.handleDeleteSite(&epb.EventDeleteSite{SiteId: "88888888-8888-8888-8888-888888888888"})

	assert.Error(t, err, "a failed delete must surface so the event is retried")
}

func TestEventNotification_NodeDeleteClearsPortMaps(t *testing.T) {
	nodeID := "uk-sa2644-tnode-v0-9a1e"

	ports := &scmocks.PortMapRepo{}
	ports.On("ClearNode", nodeID).Return(nil).Once()

	c := &SiteControllerEventServer{ports: ports, orgName: testOrgName}

	anyMsg, err := anypb.New(&epb.EventRegistryNodeDelete{NodeId: nodeID})
	assert.NoError(t, err)

	_, err = c.EventNotification(context.Background(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, "event.cloud.local.{{ .Org}}.registry.node.node.delete"),
		Msg:        anyMsg,
	})

	assert.NoError(t, err)
	ports.AssertExpectations(t)
}

func TestHandleNodeDelete_ReportsStoreFailure(t *testing.T) {
	ports := &scmocks.PortMapRepo{}
	ports.On("ClearNode", mock.Anything).Return(assert.AnError).Once()

	c := &SiteControllerEventServer{ports: ports}

	err := c.handleNodeDelete(&epb.EventRegistryNodeDelete{NodeId: "uk-sa2644-tnode-v0-9a1e"})

	assert.Error(t, err, "a failed clear must surface so the event is retried")
}

func registryErr(code int) error {
	return fmt.Errorf("GetSite failure: %w", fmt.Errorf("rest api GET failure with error: %w",
		&client.ErrorStatus{StatusCode: code, RawError: errors.New("registry error")}))
}

func TestPruneDeletedSites_DeletesOnlySitesRegistryReportsMissing(t *testing.T) {
	live := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	gone := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	unknown := "cccccccc-cccc-cccc-cccc-cccccccccccc"

	sites := &scmocks.SiteRepo{}
	sites.On("List").Return([]db.Site{{SiteID: live}, {SiteID: gone}, {SiteID: unknown}}, nil).Once()
	sites.On("Delete", gone).Return(nil).Once()

	siteClient := &mbmocks.SiteClient{}
	siteClient.On("Get", live).Return(&registry.SiteInfo{Id: live}, nil).Once()
	siteClient.On("Get", gone).Return(nil, registryErr(http.StatusNotFound)).Once()
	siteClient.On("Get", unknown).Return(nil, registryErr(http.StatusInternalServerError)).Once()

	srv := NewSiteControllerServer(testOrgName, nil, nil, nil, siteClient, nil, nil, nil, nil)
	c := &SiteControllerEventServer{s: srv, sites: sites}

	c.PruneDeletedSites()

	sites.AssertExpectations(t)
	siteClient.AssertExpectations(t)
	sites.AssertNotCalled(t, "Delete", live)
	sites.AssertNotCalled(t, "Delete", unknown)
}
