/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	contpb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/reconciler"
	"github.com/ukama/ukama/systems/node/site-controller/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SiteControllerServer struct {
	pb.UnimplementedSiteControllerServiceServer
	orgName        string
	reconciler     *reconciler.Reconciler
	dbStructs      *db.DBStruct
	msgBus         msgBusServiceClient.MsgBusServiceClient
	siteRegistry   creg.SiteClient
	nodeClient     creg.NodeClient
	healthClient   providers.HealthClientProvider
	controller     providers.ControllerClientProvider
	baseRoutingKey msgbus.RoutingKeyBuilder
}

func NewSiteControllerServer(orgName string, r *reconciler.Reconciler, mb msgBusServiceClient.MsgBusServiceClient, nodeClient creg.NodeClient, siteClient creg.SiteClient, healthClient providers.HealthClientProvider, controller providers.ControllerClientProvider, dbStructs *db.DBStruct) *SiteControllerServer {
	return &SiteControllerServer{reconciler: r, msgBus: mb, baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName), nodeClient: nodeClient, siteRegistry: siteClient, healthClient: healthClient, controller: controller, orgName: orgName, dbStructs: dbStructs}
}

func (s *SiteControllerServer) SetSite(ctx context.Context, req *pb.SetSiteRequest) (*pb.SetSiteResponse, error) {
	st, err := s.reconciler.SetSite(ctx, req.SiteId, req.State, req.Reason, req.RequestedBy)
	if err != nil {
		return nil, mapErr(err)
	}
	intent, err := s.getIntent(ctx, req.SiteId)
	if err != nil {
		return nil, mapErr(err)
	}

	return &pb.SetSiteResponse{State: siteStateToPB(st, intent)}, nil
}

func (s *SiteControllerServer) SetService(ctx context.Context, req *pb.SetServiceRequest) (*pb.SetServiceResponse, error) {
	nodeID, err := s.resolveNode(req.SiteId, ukama.NODE_ID_TYPE_TOWERNODE)
	if err != nil {
		return nil, err
	}
	client, err := s.controller.GetClient()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "node-controller unavailable: %v", err)
	}
	resp, err := client.ToggleService(ctx, &contpb.ToggleServiceRequest{NodeId: nodeID, State: req.State})
	if err != nil {
		return nil, mapErr(err)
	}
	log.Infof("site-controller: forwarded SERVICE %s for site %s to node %s, op=%s", req.State, req.SiteId, nodeID, resp.GetOperationId())
	return &pb.SetServiceResponse{}, nil
}

func (s *SiteControllerServer) SetRadio(ctx context.Context, req *pb.SetRadioRequest) (*pb.SetRadioResponse, error) {
	nodeID, err := s.resolveNode(req.SiteId, ukama.NODE_ID_TYPE_AMPNODE)
	if err != nil {
		return nil, err
	}
	client, err := s.controller.GetClient()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "node-controller unavailable: %v", err)
	}
	resp, err := client.ToggleRadio(ctx, &contpb.ToggleRadioRequest{NodeId: nodeID, State: req.State})
	if err != nil {
		return nil, mapErr(err)
	}
	log.Infof("site-controller: forwarded RADIO %s for site %s to node %s, op=%s", req.State, req.SiteId, nodeID, resp.GetOperationId())
	return &pb.SetRadioResponse{}, nil
}

func (s *SiteControllerServer) RestartSite(ctx context.Context, req *pb.RestartSiteRequest) (*pb.RestartSiteResponse, error) {
	client, err := s.controller.GetClient()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "node-controller unavailable: %v", err)
	}
	nodes, err := s.nodeClient.GetNodesBySite(req.SiteId)
	if err != nil {
		return nil, mapErr(err)
	}
	operationIds := make([]string, 0)
	for _, node := range nodes.Nodes {
		if node.Type != ukama.NODE_ID_TYPE_CNODE {
			resp, err := client.RestartNode(ctx, &contpb.RestartNodeRequest{NodeId: node.Id})
			if err != nil {
				return nil, mapErr(err)
			}
			log.Infof("site-controller: restarted node %s for site %s, op=%s", node.Id, req.SiteId, resp.GetOperationId())
			operationIds = append(operationIds, resp.GetOperationId())
		}
	}

	log.Infof("site-controller: forwarded RESTART for site %s", req.SiteId)
	return &pb.RestartSiteResponse{OperationIds: operationIds, Status: "success"}, nil
}

func (s *SiteControllerServer) ToggleInternetSwitch(ctx context.Context, req *pb.ToggleInternetSwitchRequest) (*pb.ToggleInternetSwitchResponse, error) {
	client, err := s.controller.GetClient()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "node-controller unavailable: %v", err)
	}

	nodes, err := s.nodeClient.GetNodesBySite(req.SiteId)
	if err != nil {
		return nil, mapErr(err)
	}

	var cnodeId string
	for _, node := range nodes.Nodes {
		if node.Type == ukama.NODE_ID_TYPE_CNODE {
			cnodeId = node.Id
			break
		}
	}
	if cnodeId == "" {
		return nil, status.Errorf(codes.NotFound, "no CNODE found for site %s", req.SiteId)
	}

	resp, err := client.ToggleSwitchPort(ctx, &contpb.ToggleSwitchPortRequest{NodeId: cnodeId, Status: req.Status, Port: req.Port})
	if err != nil {
		return nil, mapErr(err)
	}
	log.Infof("site-controller: forwarded internet switch for site %s port %d to %v, op=%s", req.SiteId, req.Port, req.Status, resp.GetOperationId())
	return &pb.ToggleInternetSwitchResponse{OperationId: resp.GetOperationId(), ResourceKey: resp.GetResourceKey(), Status: resp.GetStatus()}, nil
}

func (s *SiteControllerServer) resolveNode(siteID, nodeType string) (string, error) {
	if siteID == "" {
		return "", status.Errorf(codes.InvalidArgument, "site id cannot be empty")
	}
	resp, err := s.nodeClient.GetNodesBySite(siteID)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to resolve nodes for site %s: %v", siteID, err)
	}
	for _, n := range resp.Nodes {
		nId, err := ukama.ValidateNodeId(n.Id)
		if err != nil {
			continue
		}
		if t := ukama.GetNodeType(nId.String()); t != nil && *t == nodeType {
			return nId.String(), nil
		}
	}
	return "", status.Errorf(codes.NotFound, "no node of type %s found for site %s", nodeType, siteID)
}

func (s *SiteControllerServer) GetSiteState(ctx context.Context, req *pb.GetSiteStateRequest) (*pb.GetSiteStateResponse, error) {
	snap, err := s.reconciler.GetSnapshot(ctx, req.SiteId)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.GetSiteStateResponse{Snapshot: snapshotToPB(snap)}, nil
}

func (s *SiteControllerServer) UpsertPortMap(ctx context.Context, req *pb.UpsertPortMapRequest) (*pb.UpsertPortMapResponse, error) {
	ports := make([]db.SitePortMap, 0, len(req.Ports))
	for _, p := range req.Ports {
		// cn := p.CnodeId
		// if cn == "" {
		// 	cn = req.CnodeId
		// }
		ports = append(ports, db.SitePortMap{
			Port: int(p.Port), Role: p.Role, NodeID: p.NodeId, Class: p.Class, Policy: p.Policy,
		})
	}
	if err := s.reconciler.UpsertPortMap(ctx, req.SiteId, req.CnodeId, ports); err != nil {
		return nil, mapErr(err)
	}
	return &pb.UpsertPortMapResponse{}, nil
}

func (s *SiteControllerServer) GetPortMap(ctx context.Context, req *pb.GetPortMapRequest) (*pb.GetPortMapResponse, error) {
	ports, err := s.reconciler.GetPortMap(ctx, req.SiteId)
	if err != nil {
		return nil, mapErr(err)
	}
	out := make([]*pb.PortMapEntry, 0, len(ports))
	for _, p := range ports {
		out = append(out, &pb.PortMapEntry{
			Port: int32(p.Port), Role: p.Role, NodeId: p.NodeID, Class: p.Class, Policy: p.Policy,
		})
	}
	return &pb.GetPortMapResponse{Ports: out}, nil
}

func (s *SiteControllerServer) ApplySwitchPolicy(ctx context.Context, req *pb.ApplySwitchPolicyRequest) (*pb.ApplySwitchPolicyResponse, error) {
	if err := s.reconciler.ApplySwitchPolicy(ctx, req.SiteId); err != nil {
		return nil, mapErr(err)
	}
	return &pb.ApplySwitchPolicyResponse{Applied: true}, nil
}

func (s *SiteControllerServer) PowerCycleNode(ctx context.Context, req *pb.PowerCycleNodeRequest) (*pb.PowerCycleNodeResponse, error) {
	if err := s.reconciler.PowerCycleNode(ctx, req.SiteId, req.Role, req.Reason); err != nil {
		return nil, mapErr(err)
	}
	return &pb.PowerCycleNodeResponse{}, nil
}

func (s *SiteControllerServer) getIntent(ctx context.Context, siteID string) (*db.SiteIntent, error) {
	_, intent, err := s.reconciler.GetState(ctx, siteID)
	return intent, err
}

func intentToPB(in *db.SiteIntent) *pb.SiteIntentMsg {
	if in == nil {
		return nil
	}
	return &pb.SiteIntentMsg{
		SiteId: in.SiteID, DesiredService: in.DesiredService,
		DesiredRadio: in.DesiredRadio, Reason: in.Reason, RequestedBy: in.RequestedBy,
	}
}

func siteStateToPB(st *db.SiteState, intent *db.SiteIntent) *pb.DerivedStateMsg {
	if st == nil {
		return nil
	}
	out := &pb.DerivedStateMsg{
		SiteId: st.SiteID, Power: st.PowerState, Service: st.ServiceState, Radio: st.RadioState, Access: st.AccessState, Reason: st.Reason,
	}
	if intent != nil {
		out.DesiredService = intent.DesiredService
		out.DesiredRadio = intent.DesiredRadio
	}
	return out
}

func snapshotToPB(s *reconciler.SiteSnapshot) *pb.SiteSnapshot {
	if s == nil {
		return nil
	}
	ports := make([]*pb.PortMapEntry, 0, len(s.Ports))
	for _, p := range s.Ports {
		ports = append(ports, &pb.PortMapEntry{
			Port: int32(p.Port), Role: p.Role, NodeId: p.NodeID, Class: p.Class, Policy: p.Policy,
		})
	}
	return &pb.SiteSnapshot{
		Intent:         intentToPB(s.Intent),
		Derived:        siteStateToPB(s.ObservedState, s.Intent),
		ComponentsJson: s.ComponentsJSON,
		Ports:          ports,
	}
}

func mapErr(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := status.FromError(err); ok {
		return err
	}
	msg := err.Error()
	if strings.Contains(msg, "missing") || strings.Contains(msg, "invalid") {
		return status.Errorf(codes.InvalidArgument, "%s", msg)
	}
	if strings.Contains(msg, "not found") {
		return status.Errorf(codes.NotFound, "%s", msg)
	}
	return status.Errorf(codes.Internal, "%s", msg)
}
