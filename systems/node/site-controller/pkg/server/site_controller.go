package server

import (
	"context"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/reconciler"
)

type SiteControllerServer struct {
	pb.UnimplementedSiteControllerServiceServer
	reconciler *reconciler.Reconciler
}

func NewSiteControllerServer(r *reconciler.Reconciler) *SiteControllerServer {
	return &SiteControllerServer{reconciler: r}
}
func (s *SiteControllerServer) SetSite(ctx context.Context, req *pb.SetSiteRequest) (*pb.SetSiteResponse, error) {
	st, err := s.reconciler.SetSite(ctx, req.SiteId, req.State, req.Reason)
	if err != nil {
		return nil, err
	}
	intent, _ := s.getIntent(ctx, req.SiteId)
	return &pb.SetSiteResponse{State: toPB(st, intent)}, nil
}
func (s *SiteControllerServer) SetService(ctx context.Context, req *pb.SetServiceRequest) (*pb.SetServiceResponse, error) {
	st, err := s.reconciler.SetService(ctx, req.SiteId, req.State, req.Reason)
	if err != nil {
		return nil, err
	}
	intent, _ := s.getIntent(ctx, req.SiteId)
	return &pb.SetServiceResponse{State: toPB(st, intent)}, nil
}
func (s *SiteControllerServer) SetRadio(ctx context.Context, req *pb.SetRadioRequest) (*pb.SetRadioResponse, error) {
	st, err := s.reconciler.SetRadio(ctx, req.SiteId, req.State, req.Reason)
	if err != nil {
		return nil, err
	}
	intent, _ := s.getIntent(ctx, req.SiteId)
	return &pb.SetRadioResponse{State: toPB(st, intent)}, nil
}
func (s *SiteControllerServer) GetSiteState(ctx context.Context, req *pb.GetSiteStateRequest) (*pb.GetSiteStateResponse, error) {
	st, intent, err := s.reconciler.GetState(ctx, req.SiteId)
	if err != nil {
		return nil, err
	}
	return &pb.GetSiteStateResponse{State: toPB(st, intent)}, nil
}
func (s *SiteControllerServer) UpsertPortMap(ctx context.Context, req *pb.UpsertPortMapRequest) (*pb.UpsertPortMapResponse, error) {
	ports := make([]db.SitePortMap, 0, len(req.Ports))
	for _, p := range req.Ports {
		ports = append(ports, db.SitePortMap{Port: int(p.Port), Role: p.Role, NodeID: p.NodeId, Class: p.Class, Policy: p.Policy, CNodeID: req.CnodeId})
	}
	return &pb.UpsertPortMapResponse{}, s.reconciler.UpsertPortMap(ctx, req.SiteId, req.CnodeId, ports)
}
func (s *SiteControllerServer) GetPortMap(ctx context.Context, req *pb.GetPortMapRequest) (*pb.GetPortMapResponse, error) {
	ports, err := s.reconciler.GetPortMap(ctx, req.SiteId)
	if err != nil {
		return nil, err
	}
	out := make([]*pb.PortMapEntry, 0, len(ports))
	for _, p := range ports {
		out = append(out, &pb.PortMapEntry{Port: int32(p.Port), Role: p.Role, NodeId: p.NodeID, Class: p.Class, Policy: p.Policy})
	}
	return &pb.GetPortMapResponse{Ports: out}, nil
}
func (s *SiteControllerServer) ApplySwitchPolicy(ctx context.Context, req *pb.ApplySwitchPolicyRequest) (*pb.ApplySwitchPolicyResponse, error) {
	err := s.reconciler.ApplySwitchPolicy(ctx, req.SiteId)
	return &pb.ApplySwitchPolicyResponse{Applied: err == nil}, err
}
func (s *SiteControllerServer) PowerCycleNode(ctx context.Context, req *pb.PowerCycleNodeRequest) (*pb.PowerCycleNodeResponse, error) {
	return &pb.PowerCycleNodeResponse{}, s.reconciler.PowerCycleNode(ctx, req.SiteId, req.Role, req.Reason)
}
func (s *SiteControllerServer) getIntent(ctx context.Context, siteID string) (*db.SiteIntent, error) {
	_, intent, err := s.reconciler.GetState(ctx, siteID)
	return intent, err
}
func toPB(st *db.SiteState, intent *db.SiteIntent) *pb.SiteState {
	if st == nil {
		return nil
	}
	out := &pb.SiteState{SiteId: st.SiteID, Power: st.PowerState, Service: st.ServiceState, Radio: st.RadioState, Access: st.AccessState, Reason: st.Reason}
	if intent != nil {
		out.DesiredSite = intent.DesiredSite
		out.DesiredService = intent.DesiredService
		out.DesiredRadio = intent.DesiredRadio
	}
	return out
}
