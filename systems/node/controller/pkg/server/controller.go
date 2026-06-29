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
	"fmt"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

	"github.com/ukama/ukama/systems/common/ukama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"

	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	cclient "github.com/ukama/ukama/systems/node/controller/pkg/client"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"
	opmonpb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
	opmgrpb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

var actions = map[string]struct {
	path   string
	method string
}{
	"RESTART": {path: "/device/v1/restart", method: "POST"},
	"PING":    {path: "/device/v1/ping", method: "GET"},
	"SWITCH":  {path: "/device/v1/switch", method: "POST"},
	"RADIO":   {path: "/device/v1/radio", method: "POST"},
	"SERVICE": {path: "/device/v1/service", method: "POST"},
}

type ControllerServer struct {
	pb.UnimplementedControllerServiceServer
	nRepo                db.NodeLogRepo
	nodeFeederRoutingKey msgbus.RoutingKeyBuilder
	msgbus               mb.MsgBusServiceClient
	networkClient        creg.NetworkClient
	siteClient           creg.SiteClient
	nodeClient           creg.NodeClient
	opManager            copr.ManagerClient
	opMonitor            cclient.OperationMonitor
	opLeaseSecs          uint32
	opDeadlineSecs       uint32
	debug                bool
	orgName              string
}

func NewControllerServer(orgName string, nRepo db.NodeLogRepo, msgBus mb.MsgBusServiceClient, cnet creg.NetworkClient, csite creg.SiteClient, cnode creg.NodeClient, opMgr copr.ManagerClient, opMon cclient.OperationMonitor, leaseSecs, deadlineSecs uint32, debug bool) *ControllerServer {
	return &ControllerServer{
		nRepo:                nRepo,
		orgName:              orgName,
		msgbus:               msgBus,
		debug:                debug,
		networkClient:        cnet,
		siteClient:           csite,
		nodeClient:           cnode,
		opManager:            opMgr,
		opMonitor:            opMon,
		opLeaseSecs:          leaseSecs,
		opDeadlineSecs:       deadlineSecs,
		nodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (c *ControllerServer) SendNodeCommand(ctx context.Context, req *pb.SendNodeCommandRequest) (*pb.SendNodeCommandResponse, error) {
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
	}
	if _, err = c.nRepo.Get(nId.String()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Node has not been registered yet: %s", err.Error())
	}

	op, err := c.acquireAndRegister("SendNodeCommand", "node:"+nId.String())
	if err != nil {
		return nil, err
	}
	if err := c.markRunning(op, "RestartNode"); err != nil {
		c.failOperation(op, "RestartNode", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(codes.Internal, "mark running: %v", err)
	}
	if err := c.publishMessage(c.orgName+"..."+nId.String(), req.Method, req.Path, nId.String(), req.Body); err != nil {
		c.failOperation(op, "SendNodeCommand", fmt.Sprintf("publish failed: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())
	}

	return &pb.SendNodeCommandResponse{OperationId: op.Id, ResourceKey: op.ResourceKey, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func (c *ControllerServer) RestartNode(ctx context.Context, req *pb.RestartNodeRequest) (*pb.RestartNodeResponse, error) {
	if req.NodeId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "node ID cannot be empty")
	}
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
	}
	if _, err = c.nRepo.Get(nId.String()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Node has not been registered yet: %s", err.Error())
	}

	op, err := c.acquireAndRegister("RestartNode", "node:"+nId.String())
	if err != nil {
		return nil, err
	}
	if err := c.markRunning(op, "RestartNode"); err != nil {
		c.failOperation(op, "RestartNode", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(codes.Internal, "mark running: %v", err)
	}
	if err := c.publishMessage(c.orgName+"..."+nId.String(), actions["RESTART"].method, actions["RESTART"].path, nId.String(), []byte("")); err != nil {
		c.failOperation(op, "RestartNode", fmt.Sprintf("publish failed: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())
	}
	return &pb.RestartNodeResponse{OperationId: op.Id, ResourceKey: op.ResourceKey, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func (c *ControllerServer) PingNode(ctx context.Context, req *pb.PingNodeRequest) (*pb.PingNodeResponse, error) {
	if req.NodeId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "node ID cannot be empty")
	}

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	err = c.publishMessage(c.orgName+"."+"."+"."+nId.String(), actions["PING"].method, actions["PING"].path, nId.String(), []byte(""))
	if err != nil {
		log.Errorf("Failed to publish message. Errors %s", err.Error())
		return nil, status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())

	}

	return &pb.PingNodeResponse{}, nil
}

func (c *ControllerServer) ToggleSwitchPort(ctx context.Context, req *pb.ToggleSwitchPortRequest) (*pb.ToggleSwitchPortResponse, error) {
	log.Infof("Toggling internet switch for node %v, port %v to %v", req.NodeId, req.Port, req.Status)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
	}

	data, err := proto.Marshal(&pb.ToggleSwitchPortRequest{
		NodeId: nId.String(),
		Status: req.Status,
		Port:   req.Port,
	})
	if err != nil {
		return nil, err
	}

	op, err := c.acquireAndRegister("ToggleInternetSwitch", nodeKey(nId.String()))
	if err != nil {
		return nil, err
	}
	if err := c.markRunning(op, "ToggleInternetSwitch"); err != nil {
		c.failOperation(op, "ToggleInternetSwitch", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(codes.Internal, "mark running: %v", err)
	}
	if err := c.publishMessage(c.orgName+"..."+nId.String(), actions["SWITCH"].method, actions["SWITCH"].path, nId.String(), data); err != nil {
		c.failOperation(op, "ToggleInternetSwitch", fmt.Sprintf("publish failed: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to publish switch port reboot message: %s", err.Error())
	}
	return &pb.ToggleSwitchPortResponse{OperationId: op.Id, ResourceKey: op.ResourceKey, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func (c *ControllerServer) ToggleRadio(ctx context.Context, req *pb.ToggleRadioRequest) (*pb.ToggleRadioResponse, error) {
	log.Infof("Toggling RADIO on/off for node %v, to %v", req.NodeId, req.State)
	// TODO: RF toggle will send command to Tnode and Anode both
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
	}
	ntype := ukama.GetNodeType(nId.String())
	if *ntype != ukama.NODE_ID_TYPE_AMPNODE {
		return nil, status.Errorf(codes.InvalidArgument, "node is not an amplifier node")
	}

	data, err := json.Marshal(map[string]string{"state": req.State})
	if err != nil {
		return nil, err
	}

	op, err := c.acquireAndRegister("ToggleRfSwitch", nodeKey(nId.String()))
	if err != nil {
		return nil, err
	}
	if err := c.markRunning(op, "ToggleRfSwitch"); err != nil {
		c.failOperation(op, "ToggleRfSwitch", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(codes.Internal, "mark running: %v", err)
	}
	if err := c.publishMessage(fmt.Sprintf("%s...%s", c.orgName, req.NodeId), actions["RADIO"].method, actions["RADIO"].path, nId.String(), data); err != nil {
		c.failOperation(op, "ToggleRfSwitch", fmt.Sprintf("publish failed: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to publish RADIO switch message: %s", err.Error())
	}
	return &pb.ToggleRadioResponse{OperationId: op.Id, ResourceKey: op.ResourceKey, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func (c *ControllerServer) ToggleNodeService(ctx context.Context, req *pb.ToggleServiceRequest) (*pb.ToggleServiceResponse, error) {
	log.Infof("Toggling Node SERVICE on/off for node %v, to %v", req.NodeId, req.State)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
	}
	ntype := ukama.GetNodeType(nId.String())
	if *ntype != ukama.NODE_ID_TYPE_TOWERNODE {
		return nil, status.Errorf(codes.InvalidArgument, "node is not a tower node")
	}

	data, err := json.Marshal(map[string]string{"state": req.State})
	if err != nil {
		return nil, err
	}

	op, err := c.acquireAndRegister("ToggleNodeService", nodeKey(nId.String()))
	if err != nil {
		return nil, err
	}
	if err := c.markRunning(op, "ToggleNodeService"); err != nil {
		c.failOperation(op, "ToggleNodeService", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(codes.Internal, "mark running: %v", err)
	}
	if err := c.publishMessage(fmt.Sprintf("%s...%s", c.orgName, req.NodeId), actions["SERVICE"].method, actions["SERVICE"].path, nId.String(), data); err != nil {
		c.failOperation(op, "ToggleNodeService", fmt.Sprintf("publish failed: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to publish Node SERVICE switch message: %s", err.Error())
	}
	return &pb.ToggleServiceResponse{OperationId: op.Id, ResourceKey: op.ResourceKey, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func nodeKey(nodeID string) string {
	return "node:" + nodeID
}

func (c *ControllerServer) acquireAndRegister(actionType, resourceKey string) (*copr.OperationInfo, error) {
	if c.opManager == nil || c.opMonitor == nil {
		log.Warnf("%s running without operation manager/monitor for %s", actionType, resourceKey)
		return nil, fmt.Errorf("operation manager/monitor is not set")
	}

	startResp, err := c.opManager.Start(copr.StartRequest{
		Type:         actionType,
		System:       "node",
		ResourceKey:  resourceKey,
		RequestedBy:  pkg.ServiceName,
		LeaseSeconds: c.opLeaseSecs,
	})
	if err != nil {
		log.Warnf("%s lock acquire for %s rejected: %v", actionType, resourceKey, err)
		return nil, err
	}
	op := startResp.Operation
	if _, err := c.opMonitor.Register(&opmonpb.RegisterIntentRequest{
		OperationId:     op.Id,
		ResourceKey:     resourceKey,
		ActionType:      actionType,
		FencingToken:    op.FencingToken,
		DeadlineSeconds: c.opDeadlineSecs,
	}); err != nil {
		log.Errorf("%s register intent for op %s failed: %v", actionType, op.Id, err)
		c.failOperation(op, actionType, fmt.Sprintf("register intent failed: %v", err))
		return nil, status.Errorf(codes.Internal, "register intent: %v", err)
	}
	log.Infof("%s acquired lock op=%s token=%d for %s", actionType, op.Id, op.FencingToken, resourceKey)
	return op, nil
}

func (c *ControllerServer) markRunning(op *copr.OperationInfo, actionType string) error {
	if c.opManager == nil || op == nil || op.Id == "" {
		return nil
	}

	if _, err := c.opManager.MarkRunning(op.Id, op.FencingToken); err != nil {
		log.Warnf("%s mark running failed for op %s: %v", actionType, op.Id, err)
		return err
	}
	return nil
}

func (c *ControllerServer) failOperation(op *copr.OperationInfo, actionType, reason string) {
	if c.opManager == nil || op == nil || op.Id == "" {
		return
	}
	if _, err := c.opManager.ForceUnlock(op.Id, pkg.ServiceName, reason); err != nil {
		log.Errorf("%s failed to force unlock operation %s: %v", actionType, op.Id, err)
	}
}

func (c *ControllerServer) publishMessage(target string, method string, path string, nodeId string, anyMsg []byte) error {
	route := "request.cloud.local" + "." + c.orgName + "." + pkg.SystemName + "." + pkg.ServiceName + "." + "nodefeeder" + "." + "publish"
	msg := &epb.NodeFeederMessage{
		Target:     target,
		HttpMethod: method,
		Path:       path,
		Msg:        anyMsg,
		NodeId:     nodeId,
	}
	log.Infof("Published controller %s on route %s on target %s ", anyMsg, route, target)
	err := c.msgbus.PublishRequest(route, msg)
	return err
}
