package rest

import (
	"github.com/gin-gonic/gin"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
)

func (r *Router) postSiteOnHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetSiteResponse, error) {
	return r.clients.SiteController.SetSite(req.SiteId, "on", req.Reason)
}

func (r *Router) postSiteOffHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetSiteResponse, error) {
	return r.clients.SiteController.SetSite(req.SiteId, "off", req.Reason)
}

func (r *Router) postServiceOnHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetServiceResponse, error) {
	return r.clients.SiteController.SetService(req.SiteId, "on", req.Reason)
}

func (r *Router) postServiceOffHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetServiceResponse, error) {
	return r.clients.SiteController.SetService(req.SiteId, "off", req.Reason)
}

func (r *Router) postRadioOnHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetRadioResponse, error) {
	return r.clients.SiteController.SetRadio(req.SiteId, "on", req.Reason)
}

func (r *Router) postRadioOffHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetRadioResponse, error) {
	return r.clients.SiteController.SetRadio(req.SiteId, "off", req.Reason)
}

func (r *Router) getSiteStateHandler(c *gin.Context, req *SiteStateRequest) (*pb.GetSiteStateResponse, error) {
	return r.clients.SiteController.GetSiteState(req.SiteId)
}

func (r *Router) getSiteSwitchPolicyHandler(c *gin.Context, req *SiteStateRequest) (*pb.GetSwitchPolicyResponse, error) {
	return r.clients.SiteController.GetSwitchPolicy(req.SiteId)
}

func (r *Router) postRefreshSwitchPolicyHandler(c *gin.Context, req *RefreshSwitchPolicyRequest) (*pb.RefreshSwitchPolicyResponse, error) {
	return r.clients.SiteController.RefreshSwitchPolicy(req.SiteId, req.CNodeId, req.Reason)
}

func (r *Router) putReportSwitchPolicyHandler(c *gin.Context, req *ReportSwitchPolicyRequest) (*pb.ReportSwitchPolicyResponse, error) {
	policy := &pb.SwitchPolicy{
		SiteId:    req.Policy.SiteID,
		Source:    req.Policy.Source,
		UpdatedAt: req.Policy.UpdatedAt,
		State:     req.Policy.State,
		Hash:      req.Policy.Hash,
		Error:     req.Policy.Error,
		Ports:     make([]*pb.SwitchPolicyPort, 0, len(req.Policy.Ports)),
	}
	for _, p := range req.Policy.Ports {
		policy.Ports = append(policy.Ports, &pb.SwitchPolicyPort{
			Port:   p.Port,
			Role:   p.Role,
			NodeId: p.NodeId,
			Class:  p.Class,
			Policy: p.Policy,
		})
	}
	return r.clients.SiteController.ReportSwitchPolicy(req.SiteId, req.CNodeId, policy)
}

func (r *Router) postPowerCycleNodeHandler(c *gin.Context, req *PowerCycleNodeRequest) (*pb.PowerCycleNodeResponse, error) {
	return r.clients.SiteController.PowerCycleNode(req.SiteId, req.Role, req.Reason)
}
