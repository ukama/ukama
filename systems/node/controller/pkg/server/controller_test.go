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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
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

func TestNodeKey(t *testing.T) {
	nodeId := "uk-983794-tnode-78-7830"
	assert.Equal(t, "node:"+nodeId, nodeKey(nodeId))
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

func TestControllerServer_ToggleSwitchPort_NodeLevelLock(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	siteClient := &mbmocks.SiteClient{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-tnode-78-7830"
	resourceKey := nodeKey(nodeId)

	op := &copr.OperationInfo{Id: "op-i", FencingToken: 1, ResourceKey: resourceKey}
	opMgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == resourceKey
	})).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-i", uint64(1)).Return(&copr.OperationInfo{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, siteClient, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	resp, err := s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{NodeId: nodeId, Status: true, Port: 2})

	assert.NoError(t, err)
	assert.Equal(t, resourceKey, resp.ResourceKey)
	opMgr.AssertExpectations(t)
	siteClient.AssertExpectations(t)
}

func TestControllerServer_ToggleSwitchPort_PublishFailureFailsOperation(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	siteClient := &mbmocks.SiteClient{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-tnode-78-7830"
	resourceKey := nodeKey(nodeId)

	op := &copr.OperationInfo{Id: "op-i", FencingToken: 1, ResourceKey: resourceKey}
	opMgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == resourceKey
	})).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-i", uint64(1)).Return(&copr.OperationInfo{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(assert.AnError).Once()
	opMgr.On("ForceUnlock", "op-i", mock.Anything, mock.Anything).Return(&copr.OperationInfo{}, nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, siteClient, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err := s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{NodeId: nodeId, Status: true, Port: 2})

	assert.Error(t, err)
	opMgr.AssertExpectations(t)
}

func TestControllerServer_ToggleSwitchPort_Validation(t *testing.T) {
	s := NewControllerServer(testOrgName, &mocks.NodeLogRepo{}, &mbmocks.MsgBusServiceClient{}, nil, nil, nil, nil, nil, 0, 0, pkg.IsDebugMode)

	_, err := s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{NodeId: ""})
	assert.Error(t, err)

	_, err = s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{NodeId: "not-a-uuid"})
	assert.Error(t, err)
}

func TestControllerServer_ToggleRadio_InvalidNodeType(t *testing.T) {
	s := NewControllerServer(testOrgName, &mocks.NodeLogRepo{}, &mbmocks.MsgBusServiceClient{}, nil, nil, nil, nil, nil, 0, 0, pkg.IsDebugMode)

	_, err := s.ToggleRadio(context.TODO(), &pb.ToggleRadioRequest{NodeId: "uk-983794-tnode-78-7830", State: "on"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node is not an amplifier node")
}
