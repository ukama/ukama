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
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"

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
	opManager            cclient.OperationManager
	opMonitor            cclient.OperationMonitor
	opLeaseSecs          uint32
	opDeadlineSecs       uint32
	debug                bool
	orgName              string
}

func NewControllerServer(orgName string, nRepo db.NodeLogRepo, msgBus mb.MsgBusServiceClient, cnet creg.NetworkClient, csite creg.SiteClient, cnode creg.NodeClient, opMgr cclient.OperationManager, opMon cclient.OperationMonitor, leaseSecs, deadlineSecs uint32, debug bool) *ControllerServer {
	return &ControllerServer{
		nRepo:                nRepo,
		orgName:              orgName,
		nodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:               msgBus,
		debug:                debug,
		networkClient:        cnet,
		siteClient:           csite,
		nodeClient:           cnode,
		opManager:            opMgr,
		opMonitor:            opMon,
		opLeaseSecs:          leaseSecs,
		opDeadlineSecs:       deadlineSecs,
	}
}

func (c *ControllerServer) RestartSite(ctx context.Context, req *pb.RestartSiteRequest) (*pb.RestartSiteResponse, error) {
	log.Infof("Restarting site %v", req)

	if req.SiteId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "site name cannot be empty")
	}
	if req.NetworkId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "network cannot be empty")
	}
	netId, err := uuid.FromString(req.GetNetworkId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid network ID format: %s", err.Error())
	}
	if _, err = c.siteClient.Get(req.GetSiteId()); err != nil {
		return nil, fmt.Errorf("failed to validate site %s and network %s. Error %s", req.GetSiteId(), netId.String(), err.Error())
	}
	if _, err = c.networkClient.Get(netId.String()); err != nil {
		return nil, fmt.Errorf("failed to validate network with network %s. Error %s", netId.String(), err.Error())
	}
	nodes, err := c.nodeClient.GetNodesBySite(req.SiteId)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes with site %s and network %s. Error %s", req.GetSiteId(), netId.String(), err.Error())
	}

	validatedNodeIds := make([]string, 0, len(nodes.Nodes))
	for _, node := range nodes.Nodes {
		nId, err := ukama.ValidateNodeId(node.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
		}
		if _, err = c.nRepo.Get(nId.String()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Node has not been registered yet: %s", err.Error())
		}
		validatedNodeIds = append(validatedNodeIds, nId.String())
	}

	op, err := c.acquireAndRegister("RestartSite", "site:"+req.GetSiteId())
	if err != nil {
		return nil, err
	}
	if err := c.markRunning(op, "RestartSite"); err != nil {
		c.failOperation(op, "RestartSite", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(codes.Internal, "mark running: %v", err)
	}
	for _, nodeId := range validatedNodeIds {
		if err := c.publishMessage(c.orgName+"..."+nodeId, actions["RESTART"].method, actions["RESTART"].path, nodeId, []byte("")); err != nil {
			c.failOperation(op, "RestartSite", fmt.Sprintf("publish failed: %v", err))
			return nil, status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())
		}
	}
	return &pb.RestartSiteResponse{OperationIds: []string{op.Id}, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
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

func (c *ControllerServer) RestartNodes(ctx context.Context, req *pb.RestartNodesRequest) (*pb.RestartNodesResponse, error) {
	if len(req.NodeIds) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "node IDs cannot be empty")
	}

	validatedNodeIds := make([]string, 0, len(req.NodeIds))
	for _, nodeId := range req.NodeIds {
		nId, err := ukama.ValidateNodeId(nodeId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
		}
		if _, err = c.nRepo.Get(nId.String()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Node has not been registered yet: %s", err.Error())
		}
		validatedNodeIds = append(validatedNodeIds, nId.String())
	}

	ops := make([]*opmgrpb.Operation, 0, len(validatedNodeIds))
	for _, nodeId := range validatedNodeIds {
		op, err := c.acquireAndRegister("RestartNodes", "node:"+nodeId)
		if err != nil {
			c.failOperations(ops, "RestartNodes", "RestartNodes aborted before dispatch because a node lock was unavailable")
			return nil, err
		}
		ops = append(ops, op)
	}

	operationIds := make([]string, 0, len(ops))
	for i, nodeId := range validatedNodeIds {
		data, err := proto.Marshal(&pb.RestartNodeRequest{NodeId: nodeId})
		if err != nil {
			c.failOperations(ops[i:], "RestartNodes", fmt.Sprintf("marshal failed: %v", err))
			return nil, err
		}

		op := ops[i]
		if err := c.markRunning(op, "RestartNodes"); err != nil {
			c.failOperation(op, "RestartNodes", fmt.Sprintf("mark running failed: %v", err))
			return nil, status.Errorf(codes.Internal, "mark running: %v", err)
		}
		if err := c.publishMessage(c.orgName+"..."+nodeId, actions["RESTART"].method, actions["RESTART"].path, nodeId, data); err != nil {
			c.failOperation(op, "RestartNodes", fmt.Sprintf("publish failed: %v", err))
			return nil, status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())
		}
		operationIds = append(operationIds, op.Id)
	}
	return &pb.RestartNodesResponse{OperationIds: operationIds, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func (c *ControllerServer) ToggleInternetSwitch(ctx context.Context, req *pb.ToggleInternetSwitchRequest) (*pb.ToggleInternetSwitchResponse, error) {
	log.Infof("Toggling internet switch for site %v, port %v to %v", req.SiteId, req.Port, req.Status)

	if req.SiteId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "site ID cannot be empty")
	}
	siteId, err := uuid.FromString(req.GetSiteId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid site ID format: %s", err.Error())
	}
	if _, err = c.siteClient.Get(req.SiteId); err != nil {
		return nil, fmt.Errorf("failed to validate site %s. Error %s", req.SiteId, err.Error())
	}

	data, err := proto.Marshal(&pb.ToggleInternetSwitchRequest{
		SiteId: siteId.String(),
		Status: req.Status,
		Port:   req.Port,
	})
	if err != nil {
		return nil, err
	}

	op, err := c.acquireAndRegister("ToggleInternetSwitch", "site:"+siteId.String())
	if err != nil {
		return nil, err
	}
	if err := c.markRunning(op, "ToggleInternetSwitch"); err != nil {
		c.failOperation(op, "ToggleInternetSwitch", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(codes.Internal, "mark running: %v", err)
	}
	if err := c.publishMessage(c.orgName+"..."+siteId.String(), actions["SWITCH"].method, actions["SWITCH"].path, siteId.String(), data); err != nil {
		c.failOperation(op, "ToggleInternetSwitch", fmt.Sprintf("publish failed: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to publish switch port reboot message: %s", err.Error())
	}
	return &pb.ToggleInternetSwitchResponse{OperationId: op.Id, ResourceKey: op.ResourceKey, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func (c *ControllerServer) ToggleRfSwitch(ctx context.Context, req *pb.ToggleRfSwitchRequest) (*pb.ToggleRfSwitchResponse, error) {
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

	op, err := c.acquireAndRegister("ToggleRfSwitch", c.siteResourceKey(nId.String()))
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
	return &pb.ToggleRfSwitchResponse{OperationId: op.Id, ResourceKey: op.ResourceKey, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func (c *ControllerServer) ToggleNodeService(ctx context.Context, req *pb.ToggleNodeServiceRequest) (*pb.ToggleNodeServiceResponse, error) {
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

	op, err := c.acquireAndRegister("ToggleNodeService", c.siteResourceKey(nId.String()))
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
	return &pb.ToggleNodeServiceResponse{OperationId: op.Id, ResourceKey: op.ResourceKey, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func (c *ControllerServer) siteResourceKey(nodeID string) string {
	if c.nodeClient == nil {
		return "node:" + nodeID
	}
	n, err := c.nodeClient.Get(nodeID)
	if err != nil || n == nil || n.Site.SiteId == "" {
		log.Warnf("could not resolve site for node %s, using node-level lock: %v", nodeID, err)
		return "node:" + nodeID
	}
	return "site:" + n.Site.SiteId
}

func (c *ControllerServer) acquireAndRegister(actionType, resourceKey string) (*opmgrpb.Operation, error) {
	if c.opManager == nil || c.opMonitor == nil {
		log.Warnf("%s running without operation manager/monitor for %s", actionType, resourceKey)
		return &opmgrpb.Operation{
			Id:          "",
			ResourceKey: resourceKey,
		}, nil
	}

	startResp, err := c.opManager.Start(&opmgrpb.StartOperationRequest{
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

func (c *ControllerServer) markRunning(op *opmgrpb.Operation, actionType string) error {
	if c.opManager == nil || op == nil || op.Id == "" {
		return nil
	}

	if _, err := c.opManager.MarkRunning(op.Id, op.FencingToken); err != nil {
		log.Warnf("%s mark running failed for op %s: %v", actionType, op.Id, err)
		return err
	}
	return nil
}

func (c *ControllerServer) failOperation(op *opmgrpb.Operation, actionType, reason string) {
	if c.opManager == nil || op == nil || op.Id == "" {
		return
	}
	if _, err := c.opManager.FailOperation(op.Id, pkg.ServiceName, reason); err != nil {
		log.Errorf("%s failed to mark operation %s failed: %v", actionType, op.Id, err)
	}
}

func (c *ControllerServer) failOperations(ops []*opmgrpb.Operation, actionType, reason string) {
	for _, op := range ops {
		c.failOperation(op, actionType, reason)
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
