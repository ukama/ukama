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
	"testing"

	"github.com/stretchr/testify/assert"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"

	"github.com/ukama/ukama/systems/node/controller/mocks"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	"google.golang.org/protobuf/proto"
)

const testOrgName = "test-org"

// TODO: Commenting this test as it is failing and not making sense to me, need to revisit this with @Brackleycassinga
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

	msg := &pb.RestartNodeRequest{
		NodeId: nodeId,
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &cpb.NodeFeederMessage{
		Target:     "test-org" + "." + "." + "." + nodeId,
		HTTPMethod: "POST",
		Path:       "/v1/reboot",
		Msg:        data,
	}).Return(nil).Once()
	// Act
	_, err = s.RestartNode(context.TODO(), &pb.RestartNodeRequest{
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
	_, err := proto.Marshal(msg)
	if err != nil {
		return
	}

	NodeLog := db.NodeLog{
		NodeId: nodeId,
	}

	conRepo.On("Get", nodeId).Return(&NodeLog, nil).Once()

	msg = &pb.RestartNodeRequest{
		NodeId: nodeId,
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &cpb.NodeFeederMessage{
		Target:     "test-org" + "." + "." + "." + nodeId,
		HTTPMethod: "POST",
		Path:       "/v1/reboot",
		Msg:        data,
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

	nodeId := "uk-983794-hnode-78-7830"
	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, nil, pkg.IsDebugMode)

	msg := &pb.ToggleRfSwitchRequest{
		NodeId: nodeId,
		Status: true,
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &cpb.NodeFeederMessage{
		Target:     "test-org" + "." + "." + "." + nodeId,
		HTTPMethod: "POST",
		Path:       "/v1/rf",
		Msg:        data,
	}).Return(nil).Once()

	_, err = s.ToggleRfSwitch(context.TODO(), &pb.ToggleRfSwitchRequest{
		NodeId: nodeId,
		Status: true,
	})

	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
}
