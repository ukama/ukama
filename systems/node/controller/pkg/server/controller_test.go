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
	"google.golang.org/protobuf/proto"
)

const testOrgName = "test-org"

func TestControllerServer_RestartNode(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	opMgr := &mbmocks.OperatorClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-hnode-78-7830"
	setupOperationMocks(opMgr, opMon, "op-rn", "node:"+nodeId)

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	conRepo.On("Get", nodeId).Return(&db.NodeLog{NodeId: nodeId}, nil).Once()
	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     "test-org" + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/restart",
		Msg:        []byte{},
		NodeId:     nodeId,
	}).Return(nil).Once()

	resp, err := s.RestartNode(context.TODO(), &pb.RestartNodeRequest{NodeId: nodeId})

	assert.NoError(t, err)
	assert.Equal(t, "op-rn", resp.OperationId)
	assert.Equal(t, "node:"+nodeId, resp.ResourceKey)
	msgclientRepo.AssertExpectations(t)
	opMgr.AssertExpectations(t)
}

func TestControllerServer_PingNode(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	nodeId := "uk-983794-hnode-78-7830"
	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, nil, nil, 0, 0, pkg.IsDebugMode)

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     "test-org" + "." + "." + "." + nodeId,
		HttpMethod: "GET",
		Path:       "/device/v1/ping",
		Msg:        []byte{},
		NodeId:     nodeId,
	}).Return(nil).Once()

	_, err := s.PingNode(context.TODO(), &pb.PingNodeRequest{NodeId: nodeId})

	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_ToggleRadio(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	opMgr := &mocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-anode-78-7830"
	setupOperationMocks(opMgr, opMon, "op-rf", "node:"+nodeId)

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

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

	_, err = s.ToggleRadio(context.TODO(), &pb.ToggleRadioRequest{NodeId: nodeId, State: "on"})

	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_ToggleNodeService(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	opMgr := &mocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-tnode-78-7830"
	setupOperationMocks(opMgr, opMon, "op-svc", "node:"+nodeId)

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

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

	_, err = s.ToggleNodeService(context.TODO(), &pb.ToggleServiceRequest{NodeId: nodeId, State: "on"})

	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_ToggleNodeService_SiteLevelLock(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	nodeClient := &mbmocks.NodeClient{}
	opMgr := &mocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-tnode-78-7830"
	siteId := uuid.NewV4().String()

	nodeClient.On("Get", nodeId).Return(&registry.NodeInfo{
		Id:   nodeId,
		Site: registry.NodeSiteInfo{SiteId: siteId},
	}, nil).Once()
	setupOperationMocks(opMgr, opMon, "op-1", "site:"+siteId)
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
	opMgr := &mocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-tnode-78-7830"

	nodeClient.On("Get", nodeId).Return(nil, assert.AnError).Once()
	setupOperationMocks(opMgr, opMon, "op-2", "node:"+nodeId)
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nodeClient, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err := s.ToggleNodeService(context.TODO(), &pb.ToggleServiceRequest{NodeId: nodeId, State: "on"})

	assert.NoError(t, err)
	opMgr.AssertExpectations(t)
	nodeClient.AssertExpectations(t)
}

func TestControllerServer_ToggleNodeService_InvalidNodeId(t *testing.T) {
	s := NewControllerServer(testOrgName, &mocks.NodeLogRepo{}, &mbmocks.MsgBusServiceClient{}, nil, nil, nil, nil, nil, 0, 0, pkg.IsDebugMode)

	_, err := s.ToggleNodeService(context.TODO(), &pb.ToggleServiceRequest{
		NodeId: "uk-983794-anode-78-7830",
		State:  "on",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node is not a tower node")
}

func TestSiteKey_NormalizesSoSiteActionsAreMutuallyExclusive(t *testing.T) {
	id := uuid.NewV4()

	// A site action (ToggleSwitchPort) keys on the site id directly, while a node action
	// (ToggleNodeService/ToggleRadio) keys on the site resolved from the registry. Both must
	// produce the same key for the same site, regardless of casing, so only one operation
	// can hold the lock at a time.
	direct := siteKey(id.String())
	resolved := siteKey(strings.ToUpper(id.String()))

	assert.Equal(t, "site:"+id.String(), direct)
	assert.Equal(t, direct, resolved)
	assert.Equal(t, "site:not-a-uuid", siteKey("not-a-uuid"))
}

func TestControllerServer_ToggleSwitchPort_SiteLevelLock(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	siteClient := &mbmocks.SiteClient{}
	opMgr := &mocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	siteId := uuid.NewV4().String()

	siteClient.On("Get", siteId).Return(&registry.SiteInfo{}, nil).Once()
	setupOperationMocks(opMgr, opMon, "op-i", "site:"+siteId)
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
	opMgr := &mocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	siteId := uuid.NewV4().String()

	siteClient.On("Get", siteId).Return(&registry.SiteInfo{}, nil).Once()
	setupOperationMocks(opMgr, opMon, "op-i", "site:"+siteId)
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

func TestControllerServer_ToggleSwitchPort_PublishesMarshaledRequest(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	siteClient := &mbmocks.SiteClient{}
	opMgr := &mocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	siteId := uuid.NewV4().String()
	normalizedSiteId, err := uuid.FromString(siteId)
	assert.NoError(t, err)

	data, err := proto.Marshal(&pb.ToggleSwitchPortRequest{
		SiteId: normalizedSiteId.String(),
		Status: true,
		Port:   2,
	})
	assert.NoError(t, err)

	siteClient.On("Get", siteId).Return(&registry.SiteInfo{}, nil).Once()
	setupOperationMocks(opMgr, opMon, "op-i", "site:"+normalizedSiteId.String())
	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     "test-org" + "..." + normalizedSiteId.String(),
		HttpMethod: "POST",
		Path:       "/device/v1/switch",
		Msg:        data,
		NodeId:     normalizedSiteId.String(),
	}).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, siteClient, nil, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err = s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{SiteId: siteId, Status: true, Port: 2})

	assert.NoError(t, err)
	msgclientRepo.AssertExpectations(t)
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
	opMgr := &mocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	nodeId := "uk-983794-anode-78-7830"
	siteId := uuid.NewV4().String()

	nodeClient.On("Get", nodeId).Return(&registry.NodeInfo{
		Id:   nodeId,
		Site: registry.NodeSiteInfo{SiteId: siteId},
	}, nil).Once()
	setupOperationMocks(opMgr, opMon, "op-rf", "site:"+siteId)
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nodeClient, opMgr, opMon, 30, 60, pkg.IsDebugMode)

	resp, err := s.ToggleRadio(context.TODO(), &pb.ToggleRadioRequest{NodeId: nodeId, State: "on"})

	assert.NoError(t, err)
	assert.Equal(t, "site:"+siteId, resp.ResourceKey)
	opMgr.AssertExpectations(t)
	nodeClient.AssertExpectations(t)
}
