package rest

import (
	"github.com/gin-gonic/gin"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
)

func (r *Router) postSiteOnHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetSiteResponse, error) {
	return r.clients.SiteController.SetSite(req.SiteId, "on", req.Reason, req.RequestedBy)
}

func (r *Router) postSiteOffHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetSiteResponse, error) {
	return r.clients.SiteController.SetSite(req.SiteId, "off", req.Reason, req.RequestedBy)
}

func (r *Router) postServiceOnHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetServiceResponse, error) {
	return r.clients.SiteController.SetService(req.SiteId, "on", req.Reason, req.RequestedBy)
}

func (r *Router) postServiceOffHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetServiceResponse, error) {
	return r.clients.SiteController.SetService(req.SiteId, "off", req.Reason, req.RequestedBy)
}

func (r *Router) postRadioOnHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetRadioResponse, error) {
	return r.clients.SiteController.SetRadio(req.SiteId, "on", req.Reason, req.RequestedBy)
}

func (r *Router) postRadioOffHandler(c *gin.Context, req *SiteActionRequest) (*pb.SetRadioResponse, error) {
	return r.clients.SiteController.SetRadio(req.SiteId, "off", req.Reason, req.RequestedBy)
}

func (r *Router) getSiteStateHandler(c *gin.Context, req *SiteStateRequest) (*pb.GetSiteStateResponse, error) {
	return r.clients.SiteController.GetSiteState(req.SiteId)
}

func (r *Router) getSitePortMapHandler(c *gin.Context, req *SiteStateRequest) (*pb.GetPortMapResponse, error) {
	return r.clients.SiteController.GetPortMap(req.SiteId)
}

func (r *Router) putSitePortMapHandler(c *gin.Context, req *SitePortMapRequest) (*pb.UpsertPortMapResponse, error) {
	ports := make([]*pb.PortMapEntry, 0, len(req.Ports))
	for _, p := range req.Ports {
		ports = append(ports, &pb.PortMapEntry{Port: p.Port, Role: p.Role, NodeId: p.NodeId, Class: p.Class, Policy: p.Policy, CnodeId: p.CnodeId})
	}
	return r.clients.SiteController.UpsertPortMap(req.SiteId, req.CNodeId, ports)
}

func (r *Router) postApplySwitchPolicyHandler(c *gin.Context, req *SiteStateRequest) (*pb.ApplySwitchPolicyResponse, error) {
	return r.clients.SiteController.ApplySwitchPolicy(req.SiteId)
}

func (r *Router) postPowerCycleNodeHandler(c *gin.Context, req *PowerCycleNodeRequest) (*pb.PowerCycleNodeResponse, error) {
	return r.clients.SiteController.PowerCycleNode(req.SiteId, req.Role, req.Reason, req.RequestedBy)
}
