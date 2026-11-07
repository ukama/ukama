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
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/node/controller/mocks"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	opmonpb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
	opmgrpb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

const testOrgName = "test-org"

func onlineNodeClient() creg.NodeClient {
	c := &mbmocks.NodeClient{}
	c.On("Get", mock.Anything).
		Return(&creg.NodeInfo{Status: creg.NodeStatusInfo{Connectivity: "Online"}}, nil).Maybe()

	return c
}

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
	opMgr.On("MarkRunning", "op-restart", uint64(1), "controller").
		Return(&copr.OperationInfo{}, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     testOrgName + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/reboot",
		Msg:        []byte{},
		NodeId:     nodeId,
	}).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, onlineNodeClient(), opMgr, opMon, 30, 60, pkg.IsDebugMode)

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
	opMgr.On("MarkRunning", "op-rf", uint64(1), "controller").
		Return(&copr.OperationInfo{}, nil).Once()
	opMon.On("Cancel", &opmonpb.CancelIntentRequest{OperationId: "op-rf"}).
		Return(&opmonpb.CancelIntentResponse{}, nil).Once()
	opMgr.On("Complete", "op-rf", uint64(1), mock.Anything, mock.Anything).
		Return(&copr.OperationInfo{}, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     testOrgName + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/radio",
		Msg:        data,
		NodeId:     nodeId,
	}).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, onlineNodeClient(), opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err = s.ToggleRadio(context.TODO(), &pb.ToggleRadioRequest{NodeId: nodeId, State: "on"})

	msgclientRepo.AssertExpectations(t)
	opMgr.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestControllerServer_ToggleService(t *testing.T) {
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
	opMgr.On("MarkRunning", "op-svc", uint64(1), "controller").
		Return(&copr.OperationInfo{}, nil).Once()
	opMon.On("Cancel", &opmonpb.CancelIntentRequest{OperationId: "op-svc"}).
		Return(&opmonpb.CancelIntentResponse{}, nil).Once()
	opMgr.On("Complete", "op-svc", uint64(1), mock.Anything, mock.Anything).
		Return(&copr.OperationInfo{}, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     testOrgName + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/service",
		Msg:        data,
		NodeId:     nodeId,
	}).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, onlineNodeClient(), opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err = s.ToggleService(context.TODO(), &pb.ToggleServiceRequest{NodeId: nodeId, State: "on"})

	msgclientRepo.AssertExpectations(t)
	opMgr.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestNodeKey(t *testing.T) {
	nodeId := "uk-983794-tnode-78-7830"
	assert.Equal(t, "node:"+nodeId, nodeKey(nodeId))
}

func TestControllerServer_ToggleService_InvalidNodeId(t *testing.T) {
	s := NewControllerServer(
		testOrgName,
		&mocks.NodeLogRepo{},
		&mbmocks.MsgBusServiceClient{},
		nil, nil, nil, nil, nil,
		0, 0,
		pkg.IsDebugMode,
	)

	_, err := s.ToggleService(context.TODO(), &pb.ToggleServiceRequest{
		NodeId: "not-a-node-id",
		State:  "on",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format of node id")
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
	opMgr.On("MarkRunning", "op-i", uint64(1), "controller").
		Return(&copr.OperationInfo{}, nil).Once()
	opMon.On("Cancel", &opmonpb.CancelIntentRequest{OperationId: "op-i"}).
		Return(&opmonpb.CancelIntentResponse{}, nil).Once()
	opMgr.On("Complete", "op-i", uint64(1), mock.Anything, mock.Anything).
		Return(&copr.OperationInfo{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, siteClient, onlineNodeClient(), opMgr, opMon, 30, 60, pkg.IsDebugMode)

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
	opMgr.On("MarkRunning", "op-i", uint64(1), "controller").
		Return(&copr.OperationInfo{}, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(assert.AnError).Once()
	opMon.On("Cancel", &opmonpb.CancelIntentRequest{OperationId: "op-i"}).
		Return(&opmonpb.CancelIntentResponse{}, nil).Once()
	opMgr.On("Fail", "op-i", uint64(1), mock.Anything, mock.Anything).
		Return(&copr.OperationInfo{}, nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, siteClient, onlineNodeClient(), opMgr, opMon, 30, 60, pkg.IsDebugMode)

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

func TestControllerServer_ToggleRadio_InvalidNodeId(t *testing.T) {
	s := NewControllerServer(testOrgName, &mocks.NodeLogRepo{}, &mbmocks.MsgBusServiceClient{}, nil, nil, nil, nil, nil, 0, 0, pkg.IsDebugMode)

	_, err := s.ToggleRadio(context.TODO(), &pb.ToggleRadioRequest{NodeId: "not-a-node-id", State: "on"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format of node id")
}

func TestControllerServer_ToggleRadio_TowerNode(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	opMgr := &mbmocks.ManagerClient{}
	opMon := &mocks.OperationMonitor{}

	// Radio toggles are forwarded to both tower and amplifier nodes, so a
	// tower node id must be accepted here.
	nodeId := "uk-983794-tnode-78-7830"

	data, err := json.Marshal(map[string]string{"state": "on"})
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	op := &copr.OperationInfo{Id: "op-radio", FencingToken: 1, ResourceKey: "node:" + nodeId}
	opMgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == "node:"+nodeId
	})).Return(&copr.StartResponse{Operation: op}, nil).Once()
	opMon.On("Register", mock.Anything).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	opMgr.On("MarkRunning", "op-radio", uint64(1), "controller").
		Return(&copr.OperationInfo{}, nil).Once()
	opMon.On("Cancel", &opmonpb.CancelIntentRequest{OperationId: "op-radio"}).
		Return(&opmonpb.CancelIntentResponse{}, nil).Once()
	opMgr.On("Complete", "op-radio", uint64(1), mock.Anything, mock.Anything).
		Return(&copr.OperationInfo{}, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &epb.NodeFeederMessage{
		Target:     testOrgName + "..." + nodeId,
		HttpMethod: "POST",
		Path:       "/device/v1/radio",
		Msg:        data,
		NodeId:     nodeId,
	}).Return(nil).Once()

	s := NewControllerServer(testOrgName, conRepo, msgclientRepo, nil, nil, onlineNodeClient(), opMgr, opMon, 30, 60, pkg.IsDebugMode)

	_, err = s.ToggleRadio(context.TODO(), &pb.ToggleRadioRequest{NodeId: nodeId, State: "on"})

	msgclientRepo.AssertExpectations(t)
	opMgr.AssertExpectations(t)
	assert.NoError(t, err)
}

func offlineNodeClient() creg.NodeClient {
	c := &mbmocks.NodeClient{}
	c.On("Get", mock.Anything).
		Return(&creg.NodeInfo{Status: creg.NodeStatusInfo{Connectivity: "Offline"}}, nil).Maybe()

	return c
}

func TestControllerServer_OfflineNodeIsRejectedWithoutLocking(t *testing.T) {
	const nodeId = "uk-983794-hnode-78-7830"

	cases := []struct {
		name string
		call func(*ControllerServer) error
	}{
		{"RestartNode", func(s *ControllerServer) error {
			_, err := s.RestartNode(context.TODO(), &pb.RestartNodeRequest{NodeId: nodeId})
			return err
		}},
		{"ToggleRadio", func(s *ControllerServer) error {
			_, err := s.ToggleRadio(context.TODO(), &pb.ToggleRadioRequest{NodeId: nodeId, State: "on"})
			return err
		}},
		{"ToggleService", func(s *ControllerServer) error {
			_, err := s.ToggleService(context.TODO(), &pb.ToggleServiceRequest{NodeId: nodeId, State: "on"})
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conRepo := &mocks.NodeLogRepo{}
			conRepo.On("Get", nodeId).Return(&db.NodeLog{NodeId: nodeId}, nil).Maybe()
			opMgr := &mbmocks.ManagerClient{}

			s := NewControllerServer(testOrgName, conRepo, nil, nil, nil,
				offlineNodeClient(), opMgr, nil, 30, 60, pkg.IsDebugMode)

			err := tc.call(s)

			assert.Error(t, err)
			assert.Equal(t, codes.NotFound, status.Code(err))
			opMgr.AssertNotCalled(t, "Start", mock.Anything)
		})
	}
}

func TestControllerServer_operationIntentCleanup(t *testing.T) {
	op := &copr.OperationInfo{Id: "op-x", FencingToken: 3}
	cancelReq := &opmonpb.CancelIntentRequest{OperationId: "op-x"}

	t.Run("CompleteFailureKeepsIntent", func(t *testing.T) {
		opMgr := &mbmocks.ManagerClient{}
		opMon := &mocks.OperationMonitor{}
		opMgr.On("Complete", "op-x", uint64(3), pkg.ServiceName, "dispatched").
			Return(nil, assert.AnError).Once()

		c := &ControllerServer{opManager: opMgr, opMonitor: opMon}
		status := c.completeOperation(op, "ToggleRadio")

		assert.Equal(t, opmgrpb.OperationStatus_RUNNING.String(), status)
		opMon.AssertNotCalled(t, "Cancel", mock.Anything)
		opMgr.AssertExpectations(t)
	})

	t.Run("CancelFailureDoesNotChangeResult", func(t *testing.T) {
		opMgr := &mbmocks.ManagerClient{}
		opMon := &mocks.OperationMonitor{}
		opMgr.On("Complete", "op-x", uint64(3), pkg.ServiceName, "dispatched").
			Return(&copr.OperationInfo{}, nil).Once()
		opMon.On("Cancel", cancelReq).Return(nil, assert.AnError).Once()

		c := &ControllerServer{opManager: opMgr, opMonitor: opMon}
		status := c.completeOperation(op, "ToggleRadio")

		assert.Equal(t, opmgrpb.OperationStatus_SUCCESS.String(), status)
		opMon.AssertExpectations(t)
	})

	t.Run("FailCancelsIntent", func(t *testing.T) {
		opMgr := &mbmocks.ManagerClient{}
		opMon := &mocks.OperationMonitor{}
		opMgr.On("Fail", "op-x", uint64(3), pkg.ServiceName, "publish failed").
			Return(&copr.OperationInfo{}, nil).Once()
		opMon.On("Cancel", cancelReq).
			Return(&opmonpb.CancelIntentResponse{}, nil).Once()

		c := &ControllerServer{opManager: opMgr, opMonitor: opMon}
		c.failOperation(op, "ToggleRadio", "publish failed")

		opMon.AssertExpectations(t)
	})

	t.Run("FailFailureKeepsIntent", func(t *testing.T) {
		opMgr := &mbmocks.ManagerClient{}
		opMon := &mocks.OperationMonitor{}
		opMgr.On("Fail", "op-x", uint64(3), pkg.ServiceName, "publish failed").
			Return(nil, assert.AnError).Once()

		c := &ControllerServer{opManager: opMgr, opMonitor: opMon}
		c.failOperation(op, "ToggleRadio", "publish failed")

		opMon.AssertNotCalled(t, "Cancel", mock.Anything)
	})

	t.Run("RegisterFailureDoesNotCancel", func(t *testing.T) {
		opMgr := &mbmocks.ManagerClient{}
		opMon := &mocks.OperationMonitor{}
		opMgr.On("Start", mock.Anything).
			Return(&copr.StartResponse{Operation: op}, nil).Once()
		opMon.On("Register", mock.Anything).Return(nil, assert.AnError).Once()
		opMgr.On("Fail", "op-x", uint64(3), pkg.ServiceName, mock.Anything).
			Return(&copr.OperationInfo{}, nil).Once()

		c := &ControllerServer{opManager: opMgr, opMonitor: opMon}
		_, err := c.acquireAndRegister("ToggleRadio", "node:x")

		assert.Error(t, err)
		opMon.AssertNotCalled(t, "Cancel", mock.Anything)
		opMgr.AssertExpectations(t)
	})
}
