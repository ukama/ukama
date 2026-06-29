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
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	registry "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"

	"github.com/ukama/ukama/systems/node/controller/mocks"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	opmonpb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
)

const testOrgName = "test-org"

func TestControllerServer_RestartNode(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-hnode-78-7830"

	NodeLog := db.NodeLog{NodeId: nodeId}
	conRepo.On("Get", nodeId).Return(&NodeLog, nil).Once()

	op := &copr.OperationInfo{Id: "op-restart", FencingToken: 1, ResourceKey: "node:" + nodeId}
	opMgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == "node:"+nodeId
	})).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-restart", uint64(1)).Return(&copr.OperationInfo{}, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     testOrgName + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/restart",
		Msg:        []byte{},
		NodeId:     nodeId,
	}).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	resp, err := s.RestartNode(context.TODO(), &pb.RestartNodeRequest{NodeId: nodeId})

	msgclientRepo.AssertExpectations(t)
	opMgr.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, "op-restart", resp.OperationId)
	assert.Equal(t, "node:"+nodeId, resp.ResourceKey)
}

func TestControllerServer_ToggleRadio(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-anode-78-7830"

	jsonBody := map[string]string{"state": "on"}
	data, err := json.Marshal(jsonBody)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	op := &copr.OperationInfo{Id: "op-rf", FencingToken: 1, ResourceKey: "node:" + nodeId}
	opMgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == "node:"+nodeId
	})).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-rf", uint64(1)).Return(&copr.OperationInfo{}, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     testOrgName + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/radio",
		Msg:        data,
		NodeId:     nodeId,
	}).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err = s.ToggleRadio(context.TODO(), &pb.ToggleRadioRequest{NodeId: nodeId, State: "on"})

	msgclientRepo.AssertExpectations(t)
	opMgr.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_ToggleNodeService(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-tnode-78-7830"

	jsonBody := map[string]string{"state": "on"}
	data, err := json.Marshal(jsonBody)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	op := &copr.OperationInfo{Id: "op-svc", FencingToken: 1, ResourceKey: "node:" + nodeId}
	opMgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == "node:"+nodeId
	})).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-svc", uint64(1)).Return(&copr.OperationInfo{}, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     testOrgName + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/service",
		Msg:        data,
		NodeId:     nodeId,
	}).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err = s.ToggleNodeService(context.TODO(), &pb.ToggleServiceRequest{NodeId: nodeId, State: "on"})

	msgclientRepo.AssertExpectations(t)
	opMgr.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_ToggleNodeService_SiteLevelLock(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	nodeClient := &mbmocks.NodeClient{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-tnode-78-7830"
	siteId := uuid.NewV4().String()

	nodeClient.On("Get", nodeId).Return(&registry.NodeInfo{
		Id:   nodeId,
		Site: registry.NodeSiteInfo{SiteId: siteId},
	}, nil).Once()

	op := &copr.OperationInfo{Id: "op-1", FencingToken: 1, ResourceKey: "site:" + siteId}
	opMgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == "site:"+siteId
	})).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-1", uint64(1)).Return(&copr.OperationInfo{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nodeClient, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	resp, err := s.ToggleNodeService(context.TODO(), &pb.ToggleServiceRequest{NodeId: nodeId, State: "on"})

	assert.NoError(t, err)
	assert.Equal(t, "site:"+siteId, resp.ResourceKey)
	opMgr.AssertExpectations(t)
	nodeClient.AssertExpectations(t)
}

func TestControllerServer_ToggleNodeService_LockFallsBackToNode(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	nodeClient := &mbmocks.NodeClient{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-tnode-78-7830"

	nodeClient.On("Get", nodeId).Return(nil, assert.AnError).Once()

	op := &copr.OperationInfo{Id: "op-2", FencingToken: 1, ResourceKey: "node:" + nodeId}
	opMgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == "node:"+nodeId
	})).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-2", uint64(1)).Return(&copr.OperationInfo{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nodeClient, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err := s.ToggleNodeService(context.TODO(), &pb.ToggleServiceRequest{NodeId: nodeId, State: "on"})

	assert.NoError(t, err)
	opMgr.AssertExpectations(t)
	nodeClient.AssertExpectations(t)
}

func TestControllerServer_ToggleNodeService_InvalidNodeId(t *testing.T) {
	s := NewControllerServer(
		testOrgName,
		&mocks.NodeLogRepo{},
		&mbmocks.MsgBusServiceClient{},
		nil, nil, nil, nil, nil,
		0, 0,
		pkg.IsDebugMode,
	)

	_, err := s.ToggleNodeService(context.TODO(), &pb.ToggleServiceRequest{
		NodeId: "uk-983794-anode-78-7830",
		State:  "on",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node is not a tower node")
}

func TestSiteKey_NormalizesSoSiteActionsAreMutuallyExclusive(t *testing.T) {
	id := uuid.NewV4()

	// A site action (ToggleSwitchPort) keys on the site id directly, while a node
	// action (ToggleNodeService/ToggleRadio) keys on the site resolved from the
	// registry. Both must produce the same key for the same site, regardless of
	// casing, so only one operation can hold the lock at a time.
	direct := siteKey(id.String())
	resolved := siteKey(strings.ToUpper(id.String()))

	assert.Equal(t, "site:"+id.String(), direct)
	assert.Equal(t, direct, resolved)

	// Non-uuid input falls back to the raw value rather than dropping the prefix.
	assert.Equal(t, "site:not-a-uuid", siteKey("not-a-uuid"))
}

func TestControllerServer_ToggleSwitchPort_SiteLevelLock(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	siteClient := &mbmocks.SiteClient{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	siteId := uuid.NewV4().String()

	siteClient.On("Get", siteId).Return(&registry.SiteInfo{}, nil).Once()

	op := &copr.OperationInfo{Id: "op-i", FencingToken: 1, ResourceKey: "site:" + siteId}
	opMgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == "site:"+siteId
	})).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-i", uint64(1)).Return(&copr.OperationInfo{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, siteClient, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	resp, err := s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{SiteId: siteId, Status: true, Port: 2})

	assert.NoError(t, err)
	assert.Equal(t, "site:"+siteId, resp.ResourceKey)
	opMgr.AssertExpectations(t)
	siteClient.AssertExpectations(t)
}

func TestControllerServer_ToggleSwitchPort_PublishFailureFailsOperation(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	siteClient := &mbmocks.SiteClient{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	siteId := uuid.NewV4().String()

	siteClient.On("Get", siteId).Return(&registry.SiteInfo{}, nil).Once()

	op := &copr.OperationInfo{Id: "op-i", FencingToken: 1, ResourceKey: "site:" + siteId}
	opMgr.On("Start", mock.Anything).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-i", uint64(1)).Return(&copr.OperationInfo{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(assert.AnError).Once()
	opMgr.On("ForceUnlock", "op-i", mock.Anything, mock.Anything).Return(&copr.OperationInfo{}, nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, siteClient, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err := s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{SiteId: siteId, Status: true, Port: 2})

	assert.Error(t, err)
	opMgr.AssertExpectations(t)
}

func TestControllerServer_ToggleSwitchPort_Validation(t *testing.T) {
	s := NewControllerServer(testOrgName, &mocks.NodeLogRepo{}, &mbmocks.MsgBusServiceClient{}, nil, nil, nil, nil, nil, 0, 0, pkg.IsDebugMode)

	_, err := s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{SiteId: ""})
	assert.Error(t, err)

	_, err = s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{SiteId: "not-a-uuid"})
	assert.Error(t, err)
}

func TestControllerServer_ToggleRadio_InvalidNodeType(t *testing.T) {
	s := NewControllerServer(testOrgName, &mocks.NodeLogRepo{}, &mbmocks.MsgBusServiceClient{}, nil, nil, nil, nil, nil, 0, 0, pkg.IsDebugMode)

	_, err := s.ToggleRadio(context.TODO(), &pb.ToggleRadioRequest{NodeId: "uk-983794-tnode-78-7830", State: "on"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node is not an amplifier node")
}

func TestControllerServer_ToggleRadio_SiteLevelLock(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	nodeClient := &mbmocks.NodeClient{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-anode-78-7830"
	siteId := uuid.NewV4().String()

	nodeClient.On("Get", nodeId).Return(&registry.NodeInfo{
		Id:   nodeId,
		Site: registry.NodeSiteInfo{SiteId: siteId},
	}, nil).Once()

	op := &copr.OperationInfo{Id: "op-rf", FencingToken: 1, ResourceKey: "site:" + siteId}
	opMgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == "site:"+siteId
	})).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-rf", uint64(1)).Return(&copr.OperationInfo{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nodeClient, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	resp, err := s.ToggleRadio(context.TODO(), &pb.ToggleRadioRequest{NodeId: nodeId, State: "on"})

	assert.NoError(t, err)
	assert.Equal(t, "site:"+siteId, resp.ResourceKey)
	opMgr.AssertExpectations(t)
	nodeClient.AssertExpectations(t)
}
