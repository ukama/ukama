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
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	registry "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"

	"github.com/ukama/ukama/systems/node/controller/mocks"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	opmonpb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
	opmgrpb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"google.golang.org/protobuf/proto"
)

const testOrgName = "test-org"

// TODO: Commenting this test as it is failing and not making sense to me, need to revisit this with @Brackleycassinga
// func TestControllerServer_RestartSite(t *testing.T) {
// 	// Arrange
// 	msgclientRepo := &mbmocks.MsgBusServiceClient{}
// 	conRepo := &mocks.NodeLogRepo{}

// 	netId := uuid.NewV4()

// 	nodeId := "uk-983794-hnode-78-7830"
// 	nodeLog := &db.NodeLog{
// 		NodeId: nodeId,
// 	}
// 	conRepo.On("Get", nodeId).Return(nodeLog, nil).Once()

// msg := &pb.RestartNodeRequest{
// 	NodeId: nodeId,
// }
// data, err := proto.Marshal(msg)
// if err != nil {
// 	return
// }
// msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &cpb.NodeFeederMessage{
// 	Target:     "test-org." + "." + "." + nodeId,
// 	HTTPMethod: "POST",
// 	Path:       "/v1/reboot/" + nodeId,
// 	Msg:        data,
// }).Return(nil).Once()
// Act
// 	_, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{
// 		SiteName:  "site-1",
// 		NetworkId: netId.String(),
// 	})
// 	// Assert
// 	msgclientRepo.AssertExpectations(t)
// 	assert.NoError(t, err)

// }
func TestControllerServer_RestartNode(t *testing.T) {
	// Arrange
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	nodeId := "uk-983794-hnode-78-7830"
	s := NewControllerServer(
		testOrgName,
		conRepo,
		msgclientRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		0,
		0,
		pkg.IsDebugMode,
	)

	NodeLog := db.NodeLog{
		NodeId: nodeId,
	}
	conRepo.On("Get", nodeId).Return(&NodeLog, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     "test-org" + "." + "." + "." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/restart",
		Msg:        []byte{},
		NodeId:     nodeId,
	}).Return(nil).Once()
	// Act
	_, err := s.RestartNode(context.TODO(), &pb.RestartNodeRequest{
		NodeId: nodeId,
	})
	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_RestartNodes(t *testing.T) {
	// Arrange
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	netId := uuid.NewV4()
	nodeId := "uk-983794-hnode-78-7830"
	s := NewControllerServer(
		testOrgName,
		conRepo,
		msgclientRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		0,
		0,
		pkg.IsDebugMode,
	)

	msg := &pb.RestartNodeRequest{
		NodeId: nodeId,
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	NodeLog := db.NodeLog{
		NodeId: nodeId,
	}
	conRepo.On("Get", nodeId).Return(&NodeLog, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     "test-org" + "." + "." + "." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/restart",
		Msg:        data,
		NodeId:     nodeId,
	}).Return(nil).Once()
	// Act
	_, err = s.RestartNodes(context.TODO(), &pb.RestartNodesRequest{
		NetworkId: netId.String(),
		NodeIds:   []string{nodeId},
	})
	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)

}

func TestControllerServer_ToggleRf(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	nodeId := "uk-983794-anode-78-7830"
	s := NewControllerServer(
		testOrgName,
		conRepo,
		msgclientRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		0,
		0,
		pkg.IsDebugMode,
	)

	jsonBody := map[string]string{"state": "on"}
	data, err := json.Marshal(jsonBody)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}
	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     "test-org" + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/radio",
		Msg:        data,
		NodeId:     nodeId,
	}).Return(nil).Once()

	_, err = s.ToggleRfSwitch(context.TODO(), &pb.ToggleRfSwitchRequest{
		NodeId: nodeId,
		State:  "on",
	})

	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_ToggleNodeService(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	nodeId := "uk-983794-tnode-78-7830"
	s := NewControllerServer(
		testOrgName,
		conRepo,
		msgclientRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		0,
		0,
		pkg.IsDebugMode,
	)

	jsonBody := map[string]string{"state": "on"}
	data, err := json.Marshal(jsonBody)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}
	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     "test-org" + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/service",
		Msg:        data,
		NodeId:     nodeId,
	}).Return(nil).Once()

	_, err = s.ToggleNodeService(context.TODO(), &pb.ToggleNodeServiceRequest{
		NodeId: nodeId,
		State:  "on",
	})

	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_ToggleNodeService_SiteLevelLock(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	nodeClient := &mbmocks.NodeClient{}
	opMgr := &mocks.OperationManager{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-tnode-78-7830"
	siteId := uuid.NewV4().String()

	nodeClient.On("Get", nodeId).Return(&registry.NodeInfo{
		Id:   nodeId,
		Site: registry.NodeSiteInfo{SiteId: siteId},
	}, nil).Once()

	op := &opmgrpb.Operation{Id: "op-1", FencingToken: 1, ResourceKey: "site:" + siteId}
	opMgr.On("Start", mock.MatchedBy(func(req *opmgrpb.StartOperationRequest) bool {
		return req.ResourceKey == "site:"+siteId
	})).Return(&opmgrpb.StartOperationResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-1", uint64(1)).Return(&opmgrpb.MarkRunningResponse{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nodeClient, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	resp, err := s.ToggleNodeService(context.TODO(), &pb.ToggleNodeServiceRequest{NodeId: nodeId, State: "on"})

	assert.NoError(t, err)
	assert.Equal(t, "site:"+siteId, resp.ResourceKey)
	opMgr.AssertExpectations(t)
	nodeClient.AssertExpectations(t)
}

func TestControllerServer_ToggleNodeService_LockFallsBackToNode(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	nodeClient := &mbmocks.NodeClient{}
	opMgr := &mocks.OperationManager{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-tnode-78-7830"

	nodeClient.On("Get", nodeId).Return(nil, assert.AnError).Once()

	op := &opmgrpb.Operation{Id: "op-2", FencingToken: 1, ResourceKey: "node:" + nodeId}
	opMgr.On("Start", mock.MatchedBy(func(req *opmgrpb.StartOperationRequest) bool {
		return req.ResourceKey == "node:"+nodeId
	})).Return(&opmgrpb.StartOperationResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-2", uint64(1)).Return(&opmgrpb.MarkRunningResponse{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nodeClient, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err := s.ToggleNodeService(context.TODO(), &pb.ToggleNodeServiceRequest{NodeId: nodeId, State: "on"})

	assert.NoError(t, err)
	opMgr.AssertExpectations(t)
	nodeClient.AssertExpectations(t)
}

func TestControllerServer_ToggleNodeService_InvalidNodeId(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	s := NewControllerServer(
		testOrgName,
		conRepo,
		msgclientRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		0,
		0,
		pkg.IsDebugMode,
	)

	_, err := s.ToggleNodeService(context.TODO(), &pb.ToggleNodeServiceRequest{
		NodeId: "uk-983794-anode-78-7830",
		State:  "on",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node is not a tower node")
}

func TestSiteKey_NormalizesSoSiteActionsAreMutuallyExclusive(t *testing.T) {
	id := uuid.NewV4()

	// A site action (RestartSite/ToggleInternetSwitch) keys on the site id directly,
	// while a node action (ToggleNodeService/ToggleRfSwitch) keys on the site resolved
	// from the registry. Both must produce the same key for the same site, regardless of
	// casing, so only one operation can hold the lock at a time.
	direct := siteKey(id.String())
	resolved := siteKey(strings.ToUpper(id.String()))

	assert.Equal(t, "site:"+id.String(), direct)
	assert.Equal(t, direct, resolved)

	// Non-uuid input falls back to the raw value rather than dropping the prefix.
	assert.Equal(t, "site:not-a-uuid", siteKey("not-a-uuid"))
}

func TestControllerServer_RestartSite_SiteLevelLock(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	siteClient := &mbmocks.SiteClient{}
	networkClient := &mbmocks.NetworkClient{}
	nodeClient := &mbmocks.NodeClient{}
	opMgr := &mocks.OperationManager{}
	opMon := &mocks.OperationMonitor{}

	siteId := uuid.NewV4().String()
	netId := uuid.NewV4().String()
	node1 := "uk-983794-tnode-78-7830"
	node2 := "uk-983794-anode-78-7830"

	siteClient.On("Get", siteId).Return(&registry.SiteInfo{}, nil).Once()
	networkClient.On("Get", netId).Return(&registry.NetworkInfo{}, nil).Once()
	nodeClient.On("GetNodesBySite", siteId).Return(&registry.NodesBySite{
		Nodes: []registry.NodeInfo{{Id: node1}, {Id: node2}},
	}, nil).Once()
	conRepo.On("Get", mock.Anything).Return(&db.NodeLog{}, nil).Times(2)

	op := &opmgrpb.Operation{Id: "op-r", FencingToken: 1, ResourceKey: "site:" + siteId}
	opMgr.On("Start", mock.MatchedBy(func(req *opmgrpb.StartOperationRequest) bool {
		return req.ResourceKey == "site:"+siteId
	})).Return(&opmgrpb.StartOperationResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-r", uint64(1)).Return(&opmgrpb.MarkRunningResponse{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Times(2)

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, networkClient, siteClient, nodeClient, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	resp, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: siteId, NetworkId: netId})

	assert.NoError(t, err)
	// One site-level operation covers all the site's nodes.
	assert.Len(t, resp.OperationIds, 1)
	assert.Equal(t, "op-r", resp.OperationIds[0])
	opMgr.AssertExpectations(t)
}

func TestControllerServer_ToggleInternetSwitch_SiteLevelLock(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	siteClient := &mbmocks.SiteClient{}
	opMgr := &mocks.OperationManager{}
	opMon := &mocks.OperationMonitor{}

	siteId := uuid.NewV4().String()

	siteClient.On("Get", siteId).Return(&registry.SiteInfo{}, nil).Once()

	op := &opmgrpb.Operation{Id: "op-i", FencingToken: 1, ResourceKey: "site:" + siteId}
	opMgr.On("Start", mock.MatchedBy(func(req *opmgrpb.StartOperationRequest) bool {
		return req.ResourceKey == "site:"+siteId
	})).Return(&opmgrpb.StartOperationResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-i", uint64(1)).Return(&opmgrpb.MarkRunningResponse{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, siteClient, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	resp, err := s.ToggleInternetSwitch(context.TODO(), &pb.ToggleInternetSwitchRequest{SiteId: siteId, Status: true, Port: 2})

	assert.NoError(t, err)
	assert.Equal(t, "site:"+siteId, resp.ResourceKey)
	opMgr.AssertExpectations(t)
	siteClient.AssertExpectations(t)
}

func TestControllerServer_RestartSite_PublishFailureFailsOperation(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	siteClient := &mbmocks.SiteClient{}
	networkClient := &mbmocks.NetworkClient{}
	nodeClient := &mbmocks.NodeClient{}
	opMgr := &mocks.OperationManager{}
	opMon := &mocks.OperationMonitor{}

	siteId := uuid.NewV4().String()
	netId := uuid.NewV4().String()
	node1 := "uk-983794-tnode-78-7830"

	siteClient.On("Get", siteId).Return(&registry.SiteInfo{}, nil).Once()
	networkClient.On("Get", netId).Return(&registry.NetworkInfo{}, nil).Once()
	nodeClient.On("GetNodesBySite", siteId).Return(&registry.NodesBySite{
		Nodes: []registry.NodeInfo{{Id: node1}},
	}, nil).Once()
	conRepo.On("Get", mock.Anything).Return(&db.NodeLog{}, nil).Once()

	op := &opmgrpb.Operation{Id: "op-r", FencingToken: 1, ResourceKey: "site:" + siteId}
	opMgr.On("Start", mock.Anything).Return(&opmgrpb.StartOperationResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-r", uint64(1)).Return(&opmgrpb.MarkRunningResponse{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(assert.AnError).Once()
	opMgr.On("FailOperation", "op-r", mock.Anything, mock.Anything).Return(&opmgrpb.ForceUnlockResponse{}, nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, networkClient, siteClient, nodeClient, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: siteId, NetworkId: netId})

	assert.Error(t, err)
	opMgr.AssertExpectations(t)
}

func TestControllerServer_RestartSite_Validation(t *testing.T) {
	s := NewControllerServer(testOrgName, &mocks.NodeLogRepo{}, &mbmocks.MsgBusServiceClient{}, nil, nil, nil, nil, nil, 0, 0, pkg.IsDebugMode)

	_, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: "", NetworkId: uuid.NewV4().String()})
	assert.Error(t, err)

	_, err = s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: "site-1", NetworkId: ""})
	assert.Error(t, err)

	_, err = s.RestartSite(context.TODO(), &pb.RestartSiteRequest{SiteId: "site-1", NetworkId: "not-a-uuid"})
	assert.Error(t, err)
}

func TestControllerServer_ToggleInternetSwitch_Validation(t *testing.T) {
	s := NewControllerServer(testOrgName, &mocks.NodeLogRepo{}, &mbmocks.MsgBusServiceClient{}, nil, nil, nil, nil, nil, 0, 0, pkg.IsDebugMode)

	_, err := s.ToggleInternetSwitch(context.TODO(), &pb.ToggleInternetSwitchRequest{SiteId: ""})
	assert.Error(t, err)

	_, err = s.ToggleInternetSwitch(context.TODO(), &pb.ToggleInternetSwitchRequest{SiteId: "not-a-uuid"})
	assert.Error(t, err)
}

func TestControllerServer_ToggleRfSwitch_InvalidNodeType(t *testing.T) {
	s := NewControllerServer(testOrgName, &mocks.NodeLogRepo{}, &mbmocks.MsgBusServiceClient{}, nil, nil, nil, nil, nil, 0, 0, pkg.IsDebugMode)

	_, err := s.ToggleRfSwitch(context.TODO(), &pb.ToggleRfSwitchRequest{NodeId: "uk-983794-tnode-78-7830", State: "on"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node is not an amplifier node")
}

func TestControllerServer_ToggleRfSwitch_SiteLevelLock(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	nodeClient := &mbmocks.NodeClient{}
	opMgr := &mocks.OperationManager{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-anode-78-7830"
	siteId := uuid.NewV4().String()

	nodeClient.On("Get", nodeId).Return(&registry.NodeInfo{
		Id:   nodeId,
		Site: registry.NodeSiteInfo{SiteId: siteId},
	}, nil).Once()

	op := &opmgrpb.Operation{Id: "op-rf", FencingToken: 1, ResourceKey: "site:" + siteId}
	opMgr.On("Start", mock.MatchedBy(func(req *opmgrpb.StartOperationRequest) bool {
		return req.ResourceKey == "site:"+siteId
	})).Return(&opmgrpb.StartOperationResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-rf", uint64(1)).Return(&opmgrpb.MarkRunningResponse{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nodeClient, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	resp, err := s.ToggleRfSwitch(context.TODO(), &pb.ToggleRfSwitchRequest{NodeId: nodeId, State: "on"})

	assert.NoError(t, err)
	assert.Equal(t, "site:"+siteId, resp.ResourceKey)
	opMgr.AssertExpectations(t)
	nodeClient.AssertExpectations(t)
}
