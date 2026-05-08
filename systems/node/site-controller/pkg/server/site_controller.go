package server

import (
	"context"

	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/policy"
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
	_, intent, sw, _ := s.reconciler.GetState(ctx, req.SiteId)
	return &pb.SetSiteResponse{State: toPB(st, intent, sw)}, nil
}

func (s *SiteControllerServer) SetService(ctx context.Context, req *pb.SetServiceRequest) (*pb.SetServiceResponse, error) {
	st, err := s.reconciler.SetService(ctx, req.SiteId, req.State, req.Reason)
	if err != nil {
		return nil, err
	}
	_, intent, sw, _ := s.reconciler.GetState(ctx, req.SiteId)
	return &pb.SetServiceResponse{State: toPB(st, intent, sw)}, nil
}

func (s *SiteControllerServer) SetRadio(ctx context.Context, req *pb.SetRadioRequest) (*pb.SetRadioResponse, error) {
	st, err := s.reconciler.SetRadio(ctx, req.SiteId, req.State, req.Reason)
	if err != nil {
		return nil, err
	}
	_, intent, sw, _ := s.reconciler.GetState(ctx, req.SiteId)
	return &pb.SetRadioResponse{State: toPB(st, intent, sw)}, nil
}

func (s *SiteControllerServer) GetSiteState(ctx context.Context, req *pb.GetSiteStateRequest) (*pb.GetSiteStateResponse, error) {
	st, intent, sw, err := s.reconciler.GetState(ctx, req.SiteId)
	if err != nil {
		return nil, err
	}
	return &pb.GetSiteStateResponse{State: toPB(st, intent, sw)}, nil
}

func (s *SiteControllerServer) GetSwitchPolicy(ctx context.Context, req *pb.GetSwitchPolicyRequest) (*pb.GetSwitchPolicyResponse, error) {
	sw, err := s.reconciler.GetSwitchPolicy(ctx, req.SiteId)
	if err != nil {
		return nil, err
	}
	return &pb.GetSwitchPolicyResponse{Policy: switchPolicyToPB(sw)}, nil
}

func (s *SiteControllerServer) RefreshSwitchPolicy(ctx context.Context, req *pb.RefreshSwitchPolicyRequest) (*pb.RefreshSwitchPolicyResponse, error) {
	if err := s.reconciler.RefreshSwitchPolicy(ctx, req.SiteId, req.CnodeId); err != nil {
		return nil, err
	}
	return &pb.RefreshSwitchPolicyResponse{Requested: true}, nil
}

func (s *SiteControllerServer) ReportSwitchPolicy(ctx context.Context, req *pb.ReportSwitchPolicyRequest) (*pb.ReportSwitchPolicyResponse, error) {
	sp, err := policy.FromPB(req.Policy)
	if err != nil {
		return nil, err
	}
	cache, err := s.reconciler.ReportSwitchPolicy(ctx, req.SiteId, req.CnodeId, sp)
	if err != nil {
		return nil, err
	}
	return &pb.ReportSwitchPolicyResponse{Policy: switchPolicyToPB(cache)}, nil
}

func (s *SiteControllerServer) PowerCycleNode(ctx context.Context, req *pb.PowerCycleNodeRequest) (*pb.PowerCycleNodeResponse, error) {
	return &pb.PowerCycleNodeResponse{}, s.reconciler.PowerCycleNode(ctx, req.SiteId, req.Role, req.Reason)
}

func toPB(st *db.SiteState, intent *db.SiteIntent, sw *db.SiteSwitchPolicy) *pb.SiteState {
	if st == nil {
		return nil
	}
	out := &pb.SiteState{
		SiteId:       st.SiteID,
		Power:        st.PowerState,
		Service:      st.ServiceState,
		Radio:        st.RadioState,
		Access:       st.AccessState,
		Reason:       st.Reason,
		SwitchPolicy: switchPolicyToPB(sw),
	}
	if intent != nil {
		out.DesiredSite = intent.DesiredSite
		out.DesiredService = intent.DesiredService
		out.DesiredRadio = intent.DesiredRadio
	}
	return out
}

func switchPolicyToPB(sw *db.SiteSwitchPolicy) *pb.SwitchPolicyStatus {
	if sw == nil {
		return nil
	}
	sp, _ := policy.FromCache(sw)
	return &pb.SwitchPolicyStatus{
		SiteId:  sw.SiteID,
		CnodeId: sw.CNodeID,
		State:   sw.State,
		Hash:    sw.Hash,
		Source:  sw.Source,
		Error:   sw.Error,
		Valid:   sw.Valid,
		Reason:  sw.Reason,
		Policy:  policy.ToPB(sp),
	}
}
