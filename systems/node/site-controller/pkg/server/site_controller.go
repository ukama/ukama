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

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/reconciler"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SiteControllerServer struct {
	pb.UnimplementedSiteControllerServiceServer
	orgName        string
	reconciler     *reconciler.Reconciler
	msgbus         mb.MsgBusServiceClient
	nodeClient     creg.NodeClient
	baseRoutingKey msgbus.RoutingKeyBuilder
}

func NewSiteControllerServer(orgName string, r *reconciler.Reconciler, mb mb.MsgBusServiceClient, nodeClient creg.NodeClient) *SiteControllerServer {
	return &SiteControllerServer{reconciler: r, baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName), nodeClient: nodeClient}
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
	return &pb.SetSiteResponse{State: derivedStateToPB(st, intent)}, nil
}

func (s *SiteControllerServer) SetService(ctx context.Context, req *pb.SetServiceRequest) (*pb.SetServiceResponse, error) {
	st, err := s.reconciler.SetService(ctx, req.SiteId, req.State, req.Reason, req.RequestedBy)
	if err != nil {
		return nil, mapErr(err)
	}
	intent, err := s.getIntent(ctx, req.SiteId)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.SetServiceResponse{State: derivedStateToPB(st, intent)}, nil
}

func (s *SiteControllerServer) SetRadio(ctx context.Context, req *pb.SetRadioRequest) (*pb.SetRadioResponse, error) {
	st, err := s.reconciler.SetRadio(ctx, req.SiteId, req.State, req.Reason, req.RequestedBy)
	if err != nil {
		return nil, mapErr(err)
	}
	intent, err := s.getIntent(ctx, req.SiteId)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.SetRadioResponse{State: derivedStateToPB(st, intent)}, nil
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
		cn := p.CnodeId
		if cn == "" {
			cn = req.CnodeId
		}
		ports = append(ports, db.SitePortMap{
			Port: int(p.Port), Role: p.Role, NodeID: p.NodeId, Class: p.Class, Policy: p.Policy, CNodeID: cn,
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
			Port: int32(p.Port), Role: p.Role, NodeId: p.NodeID, Class: p.Class, Policy: p.Policy, CnodeId: p.CNodeID,
		})
	}
	return &pb.GetPortMapResponse{Ports: out}, nil
}

func (s *SiteControllerServer) ApplySwitchPolicy(ctx context.Context, req *pb.ApplySwitchPolicyRequest) (*pb.ApplySwitchPolicyResponse, error) {
	err := s.reconciler.ApplySwitchPolicy(ctx, req.SiteId)
	if err != nil {
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
		SiteId: in.SiteID, DesiredSite: in.DesiredSite, DesiredService: in.DesiredService,
		DesiredRadio: in.DesiredRadio, Reason: in.Reason, RequestedBy: in.RequestedBy,
	}
}

func derivedStateToPB(st *db.SiteState, intent *db.SiteIntent) *pb.DerivedStateMsg {
	if st == nil {
		return nil
	}
	out := &pb.DerivedStateMsg{
		SiteId: st.SiteID, Power: st.PowerState, Service: st.ServiceState, Radio: st.RadioState, Access: st.AccessState, Reason: st.Reason,
	}
	if intent != nil {
		out.DesiredSite = intent.DesiredSite
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
			Port: int32(p.Port), Role: p.Role, NodeId: p.NodeID, Class: p.Class, Policy: p.Policy, CnodeId: p.CNodeID,
		})
	}
	return &pb.SiteSnapshot{
		Intent:         intentToPB(s.Intent),
		Derived:        derivedStateToPB(s.DerivedState, s.Intent),
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
