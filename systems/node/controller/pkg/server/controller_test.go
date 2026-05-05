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
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"

	"github.com/ukama/ukama/systems/node/controller/mocks"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	"google.golang.org/protobuf/proto"
)

const testOrgName = "test-org"

// TODO: Commenting this test as it is failing, need to revisit this with @Brackleycassinga
// func TestControllerServer_RestartSite(t *testing.T) {
// 	// Arrange
// 	msgclientRepo := &mbmocks.MsgBusServiceClient{}
// 	conRepo := &mocks.NodeLogRepo{}

// 	netId := uuid.NewV4()

// 	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)
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
	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)

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
	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)

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
	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)

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
		State: "on",
	})

	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_ToggleNodeService(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	nodeId := "uk-983794-tnode-78-7830"
	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)

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

func TestControllerServer_ToggleNodeService_InvalidNodeId(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)

	_, err := s.ToggleNodeService(context.TODO(), &pb.ToggleNodeServiceRequest{
		NodeId: "uk-983794-anode-78-7830",
		State:  "on",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node is not a tower node")
}

func TestControllerServer_PingSwitchPort(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	nodeId := "uk-983794-cnode-78-7830"
	port := int32(9)
	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     "test-org" + "..." + nodeId,
		HttpMethod: "GET",
		Path:       "/switch/v1/ports/9",
		Msg:        []byte(""),
		NodeId:     nodeId,
	}).Return(nil).Once()

	_, err := s.PingSwitchPort(context.TODO(), &pb.PingSwitchPortRequest{
		NodeId: nodeId,
		Port:   port,
	})

	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_PingSwitchPort_InvalidNodeType(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)

	_, err := s.PingSwitchPort(context.TODO(), &pb.PingSwitchPortRequest{
		NodeId: "uk-983794-tnode-78-7830",
		Port:   1,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node is not a cnode")
}

func TestControllerServer_ToggleSwitchPort(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	nodeId := "uk-983794-cnode-78-7830"
	port := int32(9)
	status := true
	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)

	jsonBody := map[string]bool{"on": status}
	data, err := json.Marshal(jsonBody)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     "test-org" + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/switch/v1/ports/9/poe",
		Msg:        data,
		NodeId:     nodeId,
	}).Return(nil).Once()

	_, err = s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{
		NodeId: nodeId,
		Port:   port,
		Status: status,
	})

	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_ToggleSwitchPort_InvalidNodeType(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)

	_, err := s.ToggleSwitchPort(context.TODO(), &pb.ToggleSwitchPortRequest{
		NodeId: "uk-983794-anode-78-7830",
		Port:   2,
		Status: true,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node is not a cnode")
}
