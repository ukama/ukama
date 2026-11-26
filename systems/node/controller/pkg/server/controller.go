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
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"

	log "github.com/sirupsen/logrus"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	cclient "github.com/ukama/ukama/systems/node/controller/pkg/client"
	opmonpb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
	opmgrpb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

const (
	reasonNodeAccepted    = "accepted by node"
	reasonNodeRejectedFmt = "node rejected: %d %s"
	reasonNodeDroppedFmt  = "not delivered: %s"
)

var actions = map[string]struct {
	path   string
	method string
}{
	"RESTART":       {path: "/device/v1/reboot", method: "POST"},
	"PING":          {path: "/device/v1/ping", method: "GET"},
	"SWITCH":        {path: "/device/v1/switch", method: "POST"},
	"RADIO":         {path: "/device/v1/radio", method: "POST"},
	"SERVICE":       {path: "/device/v1/service", method: "POST"},
	"CONFIG":        {path: "configd/v1/config", method: "POST"},
	"DELETE_CONFIG": {path: "configd/v1/config", method: "DELETE"},
}

type ControllerServer struct {
	configurationState ConfigurationState
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
	now                  func() time.Time
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
		now:                  time.Now,
		nodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (c *ControllerServer) onlineNode(nodeID string) (*creg.NodeInfo, error) {
	node, err := c.nodeClient.Get(nodeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to resolve node %s: %v",
			nodeID, err)
	}

	if !strings.EqualFold(node.Status.Connectivity,
		ukama.NodeConnectivityOnline.String()) {
		return nil, status.Errorf(codes.NotFound,
			"node %s is not available: connectivity is %s",
			nodeID, node.Status.Connectivity)
	}

	return node, nil
}

func (c *ControllerServer) SendNodeCommand(ctx context.Context, req *pb.SendNodeCommandRequest) (*pb.SendNodeCommandResponse, error) {
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
	}
	if _, err = c.nRepo.Get(nId.String()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Node has not been registered yet: %s", err.Error())
	}

	ops, err := c.lockNodes(ctx, "SendNodeCommand", nId.String())
	if err != nil {
		return nil, err
	}
	op := ops[0]
	if err := c.markRunning(op, "SendNodeCommand"); err != nil {
		c.failOperation(op, "SendNodeCommand", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(markRunningCode(err), "mark running: %v", err)
	}
	if err := c.publishMessage(c.orgName+"..."+nId.String(), req.Method, req.Path,
		nId.String(), op, req.Body); err != nil {
		c.failOperation(op, "SendNodeCommand", fmt.Sprintf("publish failed: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())
	}

	return &pb.SendNodeCommandResponse{OperationId: op.Id, ResourceKey: op.ResourceKey,
		Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func (c *ControllerServer) RestartNodes(ctx context.Context,
	req *pb.RestartNodesRequest) (*pb.RestartNodesResponse, error) {
	ids, err := validNodeIds(req.NodeIds)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		if _, err := c.nRepo.Get(id); err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"Node has not been registered yet: %s", err.Error())
		}
	}

	ops, err := c.lockNodes(ctx, "RestartNode", ids...)
	if err != nil {
		return nil, err
	}
	if err := c.markAllRunning(ops, "RestartNode"); err != nil {
		return nil, err
	}

	resp := &pb.RestartNodesResponse{}
	for i, op := range ops {
		if err := c.publishMessage(c.orgName+"..."+ids[i], actions["RESTART"].method,
			actions["RESTART"].path, ids[i], op, []byte("")); err != nil {
			c.failAll(ops[i:], "RestartNode", fmt.Sprintf("publish failed: %v", err))
			return nil, status.Errorf(codes.Internal, "Failed to publish message: %s",
				err.Error())
		}
		resp.Operations = append(resp.Operations, &pb.RestartNodeResponse{
			OperationId: op.Id,
			ResourceKey: op.ResourceKey,
			Status:      opmgrpb.OperationStatus_RUNNING.String(),
		})
	}
	return resp, nil
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

	err = c.publishMessage(c.orgName+"."+"."+"."+nId.String(), actions["PING"].method,
		actions["PING"].path, nId.String(), nil, []byte(""))
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

	ops, err := c.lockNodes(ctx, "ToggleInternetSwitch", nId.String())
	if err != nil {
		return nil, err
	}
	op := ops[0]
	if err := c.markRunning(op, "ToggleInternetSwitch"); err != nil {
		c.failOperation(op, "ToggleInternetSwitch", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(markRunningCode(err), "mark running: %v", err)
	}
	if err := c.publishMessage(c.orgName+"..."+nId.String(), actions["SWITCH"].method,
		actions["SWITCH"].path, nId.String(), op, data); err != nil {
		c.failOperation(op, "ToggleInternetSwitch", fmt.Sprintf("publish failed: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to publish switch port reboot message: %s", err.Error())
	}
	return &pb.ToggleSwitchPortResponse{OperationId: op.Id, ResourceKey: op.ResourceKey,
		Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func (c *ControllerServer) ToggleRadios(ctx context.Context,
	req *pb.ToggleRadiosRequest) (*pb.ToggleRadiosResponse, error) {
	log.Infof("Toggling RADIO on/off for nodes %v, to %v", req.NodeIds, req.State)
	ids, err := validNodeIds(req.NodeIds)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(map[string]string{"state": req.State})
	if err != nil {
		return nil, err
	}

	ops, err := c.lockNodes(ctx, "ToggleRadio", ids...)
	if err != nil {
		return nil, err
	}
	if err := c.markAllRunning(ops, "ToggleRadio"); err != nil {
		return nil, err
	}

	resp := &pb.ToggleRadiosResponse{}
	for i, op := range ops {
		if err := c.publishMessage(c.orgName+"..."+ids[i], actions["RADIO"].method,
			actions["RADIO"].path, ids[i], op, data); err != nil {
			c.failAll(ops[i:], "ToggleRadio", fmt.Sprintf("publish failed: %v", err))
			return nil, status.Errorf(codes.Internal,
				"Failed to publish RADIO switch message: %s", err.Error())
		}
		resp.Operations = append(resp.Operations, &pb.ToggleRadioResponse{
			OperationId: op.Id,
			ResourceKey: op.ResourceKey,
			Status:      opmgrpb.OperationStatus_RUNNING.String(),
		})
	}
	return resp, nil
}

func (c *ControllerServer) ToggleService(ctx context.Context, req *pb.ToggleServiceRequest) (*pb.ToggleServiceResponse, error) {
	log.Infof("Toggling Node SERVICE on/off for node %v, to %v", req.NodeId, req.State)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of node id. Error %s", err.Error())
	}

	data, err := json.Marshal(map[string]string{"state": req.State})
	if err != nil {
		return nil, err
	}

	ops, err := c.lockNodes(ctx, "ToggleService", nId.String())
	if err != nil {
		return nil, err
	}
	op := ops[0]
	if err := c.markRunning(op, "ToggleService"); err != nil {
		c.failOperation(op, "ToggleService", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(markRunningCode(err), "mark running: %v", err)
	}
	if err := c.publishMessage(fmt.Sprintf("%s...%s", c.orgName, req.NodeId),
		actions["SERVICE"].method, actions["SERVICE"].path, nId.String(), op,
		data); err != nil {
		c.failOperation(op, "ToggleService", fmt.Sprintf("publish failed: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to publish Node SERVICE switch message: %s", err.Error())
	}
	return &pb.ToggleServiceResponse{OperationId: op.Id, ResourceKey: op.ResourceKey,
		Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
}

func nodeKey(nodeID string) string {
	return "node:" + nodeID
}

func (c *ControllerServer) lockNodes(ctx context.Context, actionType string,
	nodeIDs ...string) ([]*copr.OperationInfo, error) {
	targets := make([]string, 0, len(nodeIDs))
	siteID := ""
	for i, id := range nodeIDs {
		node, err := c.onlineNode(id)
		if err != nil {
			return nil, err
		}
		if i > 0 && node.Site.SiteId != siteID {
			return nil, status.Errorf(codes.InvalidArgument,
				"%s: nodes %v do not belong to the same site", actionType, nodeIDs)
		}
		siteID = node.Site.SiteId
		targets = append(targets, nodeKey(id))
	}

	if c.opManager == nil || c.opMonitor == nil {
		return nil, status.Errorf(codes.Unavailable,
			"%s: operation manager/monitor is not set", actionType)
	}

	siblings, err := c.siteNodeKeys(siteID)
	if err != nil {
		return nil, err
	}

	resp, err := c.opManager.Start(copr.StartRequest{
		Type:         actionType,
		System:       "node",
		RequestedBy:  ugrpc.RequestedBy(ctx, pkg.ServiceName),
		LeaseSeconds: c.opLeaseSecs,
		ResourceKeys: targets,
		MustBeFree:   siblings,
	})
	if err == nil && (resp == nil || len(resp.Operations) == 0) {
		err = fmt.Errorf("start returned no operations")
	}
	if errors.Is(err, copr.ErrBusy) {
		log.Warnf("%s rejected for %v: %v", actionType, nodeIDs, err)
		return nil, status.Errorf(codes.AlreadyExists, "%s rejected: %v",
			actionType, err)
	}
	if err != nil {
		log.Errorf("%s lock for %v failed: %v", actionType, nodeIDs, err)
		return nil, status.Errorf(codes.Unavailable, "%s: lock nodes: %v",
			actionType, err)
	}

	ops := resp.Operations
	for i, op := range ops {
		if _, err := c.opMonitor.Register(&opmonpb.RegisterIntentRequest{
			OperationId:     op.Id,
			ResourceKey:     op.ResourceKey,
			ActionType:      actionType,
			FencingToken:    op.FencingToken,
			DeadlineSeconds: c.opDeadlineSecs,
		}); err != nil {
			log.Errorf("%s register intent for op %s failed: %v", actionType, op.Id, err)
			reason := fmt.Sprintf("register intent failed: %v", err)
			for _, registered := range ops[:i] {
				c.failOperation(registered, actionType, reason)
			}
			for _, pending := range ops[i:] {
				c.markFailed(pending, actionType, reason)
			}
			return nil, status.Errorf(codes.Internal, "register intent: %v", err)
		}
		log.Infof("%s acquired lock op=%s token=%d for %s", actionType, op.Id,
			op.FencingToken, op.ResourceKey)
	}
	return ops, nil
}

func (c *ControllerServer) siteNodeKeys(siteID string) ([]string, error) {
	if siteID == "" {
		return nil, nil
	}
	resp, err := c.nodeClient.GetNodesBySite(siteID)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "read nodes of site %s: %v",
			siteID, err)
	}
	keys := make([]string, 0, len(resp.Nodes))
	for _, n := range resp.Nodes {
		keys = append(keys, nodeKey(canonicalNodeID(n.Id)))
	}
	return keys, nil
}

func canonicalNodeID(id string) string {
	if nId, err := ukama.ValidateNodeId(id); err == nil {
		return nId.String()
	}
	return id
}

func validNodeIds(raw []string) ([]string, error) {
	if len(raw) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "node ID cannot be empty")
	}
	seen := make(map[string]bool, len(raw))
	ids := make([]string, 0, len(raw))
	for _, r := range raw {
		nId, err := ukama.ValidateNodeId(r)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of node id. Error %s", err.Error())
		}
		if !seen[nId.String()] {
			seen[nId.String()] = true
			ids = append(ids, nId.String())
		}
	}
	return ids, nil
}

func (c *ControllerServer) markAllRunning(ops []*copr.OperationInfo,
	actionType string) error {
	for _, op := range ops {
		if err := c.markRunning(op, actionType); err != nil {
			c.failAll(ops, actionType, fmt.Sprintf("mark running failed: %v", err))
			return status.Errorf(markRunningCode(err), "mark running: %v", err)
		}
	}
	return nil
}

func (c *ControllerServer) failAll(ops []*copr.OperationInfo, actionType,
	reason string) {
	for _, op := range ops {
		c.failOperation(op, actionType, reason)
	}
}

func (c *ControllerServer) markRunning(op *copr.OperationInfo, actionType string) error {
	if c.opManager == nil || op == nil || op.Id == "" {
		return nil
	}

	if _, err := c.opManager.MarkRunning(op.Id, op.FencingToken,
		pkg.ServiceName); err != nil {
		log.Warnf("%s mark running failed for op %s: %v", actionType, op.Id, err)
		return err
	}
	return nil
}

func markRunningCode(err error) codes.Code {
	switch {
	case errors.Is(err, copr.ErrNotRunnable):
		return codes.Aborted
	case errors.Is(err, copr.ErrUnavailable):
		return codes.Unavailable
	default:
		return codes.Internal
	}
}

func (c *ControllerServer) failOperation(op *copr.OperationInfo, actionType, reason string) {
	if c.markFailed(op, actionType, reason) {
		c.cancelIntent(op, actionType)
	}
}

func (c *ControllerServer) markFailed(op *copr.OperationInfo, actionType,
	reason string) bool {
	if c.opManager == nil || op == nil || op.Id == "" {
		return false
	}
	if _, err := c.opManager.Fail(op.Id, op.FencingToken, pkg.ServiceName,
		reason); err != nil {
		log.Errorf("%s failed to fail operation %s: %v", actionType, op.Id, err)
		return false
	}
	return true
}

func (c *ControllerServer) dispatchIntent(op *copr.OperationInfo, actionType string) {
	if c.opMonitor == nil {
		return
	}
	if _, err := c.opMonitor.Dispatch(&opmonpb.DispatchIntentRequest{
		OperationId: op.Id,
	}); err != nil {
		log.Warnf("%s failed to mark intent of operation %s dispatched: %v",
			actionType, op.Id, err)
	}
}

func (c *ControllerServer) cancelIntent(op *copr.OperationInfo, actionType string) {
	if c.opMonitor == nil {
		return
	}
	if _, err := c.opMonitor.Cancel(&opmonpb.CancelIntentRequest{
		OperationId: op.Id,
	}); err != nil {
		log.Warnf("%s failed to cancel intent for operation %s: %v",
			actionType, op.Id, err)
	}
}

func (c *ControllerServer) completeOperation(op *copr.OperationInfo, actionType string) string {
	if c.opManager == nil || op == nil || op.Id == "" {
		return opmgrpb.OperationStatus_RUNNING.String()
	}
	if _, err := c.opManager.Complete(op.Id, op.FencingToken, pkg.ServiceName,
		reasonNodeAccepted); err != nil {
		log.Errorf("%s failed to complete operation %s (lock frees at deadline): %v", actionType, op.Id, err)
		return opmgrpb.OperationStatus_RUNNING.String()
	}
	c.cancelIntent(op, actionType)
	return opmgrpb.OperationStatus_SUCCESS.String()
}

func (c *ControllerServer) publishMessage(target string, method string, path string,
	nodeId string, op *copr.OperationInfo, anyMsg []byte) error {
	route := "request.cloud.local" + "." + c.orgName + "." + pkg.SystemName + "." + pkg.ServiceName + "." + "nodefeeder" + "." + "publish"
	msg := &epb.NodeFeederMessage{
		Target:     target,
		HttpMethod: method,
		Path:       path,
		Msg:        anyMsg,
		NodeId:     nodeId,
	}
	if op != nil {
		msg.OperationId = op.Id
		msg.FencingToken = op.FencingToken
		msg.ExpiresAt = timestamppb.New(c.now().
			Add(time.Duration(c.opDeadlineSecs) * time.Second))
	}
	log.Infof("Published controller %s on route %s on target %s ", anyMsg, route, target)
	err := c.msgbus.PublishRequest(route, msg)
	return err
}
