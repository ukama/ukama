/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ukama/ukama/systems/common/config"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/registry/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/registry/api-gateway/pkg"
	"github.com/ukama/ukama/systems/registry/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	sitepb "github.com/ukama/ukama/systems/registry/site/pb/gen"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	log "github.com/sirupsen/logrus"
	invpb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
	mpb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
)

const (
	Undefined = -1
)

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	metricsConfig config.Metrics
	httpEndpoints *pkg.HttpEndpoints
	grpcEndpoints *pkg.GrpcEndpoints
	descriptions  *pkg.ServiceDescriptions
	debugMode     bool
	serverConf    *rest.HttpConfig
	auth          *config.Auth
}

type Clients struct {
	Network    network
	Node       node
	Member     member
	Invitation invitation
	Site       site
}

type network interface {
	AddNetwork(ctx context.Context, netName string, allowedCountries, allowedNetworks []string, budget, overdraft float64, trafficPolicy uint32, paymentLinks bool) (*netpb.AddResponse, error)
	GetNetwork(ctx context.Context, netID string) (*netpb.GetResponse, error)
	GetNetworks(ctx context.Context) (*netpb.GetNetworksResponse, error)
	SetNetworkDefault(ctx context.Context, netID string) (*netpb.SetDefaultResponse, error)
	GetDefault(ctx context.Context) (*netpb.GetDefaultResponse, error)
	RemoveNetwork(ctx context.Context, netID string) (*netpb.DeleteResponse, error)
}

type site interface {
	AddSiteContext(context.Context, *sitepb.AddRequest) (*sitepb.AddResponse, error)
	AddSite(ctx context.Context, networkId, name, backhaulId, powerId, accessId, switchId, location, spectrumId string, isDeactivated bool, latitude, longitude string, installDate string) (*sitepb.AddResponse, error)
	GetSite(ctx context.Context, siteId string) (*sitepb.GetResponse, error)
	List(ctx context.Context, networkId string, isDeactivate bool) (*sitepb.ListResponse, error)
	UpdateSite(ctx context.Context, siteId, name string) (*sitepb.UpdateResponse, error)
	RemoveSite(ctx context.Context, siteId string) (*sitepb.DeleteResponse, error)
}

type invitation interface {
	AddInvitation(ctx context.Context, name, email, role string) (*invpb.AddResponse, error)
	GetInvitationById(ctx context.Context, invitationId string) (*invpb.GetResponse, error)
	UpdateInvitation(ctx context.Context, invitationId string, status string, email string) (*invpb.UpdateStatusResponse, error)
	RemoveInvitation(ctx context.Context, invitationId string) (*invpb.DeleteResponse, error)
	GetAllInvitations(ctx context.Context) (*invpb.GetAllResponse, error)
	GetInvitationsByEmail(ctx context.Context, email string) (*invpb.GetByEmailResponse, error)
}

type member interface {
	GetMember(ctx context.Context, userUUID string) (*mpb.MemberResponse, error)
	GetMemberByUserId(ctx context.Context, userUUID string) (*mpb.GetMemberByUserIdResponse, error)
	GetMembers(ctx context.Context) (*mpb.GetMembersResponse, error)
	AddMember(ctx context.Context, userUUID string, role string) (*mpb.MemberResponse, error)
	UpdateMember(ctx context.Context, userUUID string, isDeactivated bool, role string) error
	RemoveMember(ctx context.Context, userUUID string) error
}

type node interface {
	AddNode(ctx context.Context, nodeId, name, state string) (*nodepb.AddNodeResponse, error)
	GetNode(ctx context.Context, nodeId string) (*nodepb.GetNodeResponse, error)
	GetNodes(ctx context.Context) (*nodepb.GetNodesResponse, error)
	List(ctx context.Context, req *nodepb.ListRequest) (*nodepb.ListResponse, error)
	GetNetworkNodes(ctx context.Context, networkId string) (*nodepb.GetByNetworkResponse, error)
	GetSiteNodes(ctx context.Context, siteId string) (*nodepb.GetBySiteResponse, error)
	GetNodesByState(ctx context.Context, connectivity, state string) (*nodepb.GetNodesResponse, error)
	UpdateNodeState(ctx context.Context, nodeId string, state string) (*nodepb.UpdateNodeResponse, error)
	UpdateNode(ctx context.Context, nodeId string, name string, latitude string, longitude string) (*nodepb.UpdateNodeResponse, error)
	DeleteNode(ctx context.Context, nodeId string) (*nodepb.DeleteNodeResponse, error)
	AttachNodes(ctx context.Context, node, l, r string) (*nodepb.AttachNodesResponse, error)
	DetachNode(ctx context.Context, nodeId string) (*nodepb.DetachNodeResponse, error)
	AddNodeToSite(ctx context.Context, nodeId, networkId, siteId string) (*nodepb.AddNodeToSiteResponse, error)
	ReleaseNodeFromSite(ctx context.Context, nodeId string) (*nodepb.ReleaseNodeFromSiteResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Network = client.NewNetworkRegistry(endpoints.Network, endpoints.Timeout)
	c.Node = client.NewNode(endpoints.Node, endpoints.Timeout)
	c.Member = client.NewMemberRegistry(endpoints.Member, endpoints.Timeout)
	c.Invitation = client.NewInvitationRegistry(endpoints.Invitation, endpoints.Timeout)
	c.Site = client.NewSiteRegistry(endpoints.Site, endpoints.Timeout)

	return c
}

func NewRouter(clients *Clients, config *RouterConfig, authfunc func(*gin.Context, string) error) *Router {
	r := &Router{
		clients: clients,
		config:  config,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init(authfunc)
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		httpEndpoints: &svcConf.Http,
		grpcEndpoints: &svcConf.Services,
		descriptions:  &svcConf.Descriptions,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
		auth:          svcConf.Auth,
	}
}

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")

	desc := r.config.descriptions
	if desc == nil {
		desc = &pkg.ServiceDescriptions{}
	}

	if r.config.grpcEndpoints != nil {
		rest.RegisterStatusEndpoint(r.f, pkg.SystemName, map[string]rest.StatusTarget{
			"network":    {Host: r.config.grpcEndpoints.Network, Description: desc.Network},
			"member":     {Host: r.config.grpcEndpoints.Member, Description: desc.Member},
			"node":       {Host: r.config.grpcEndpoints.Node, Description: desc.Node},
			"invitation": {Host: r.config.grpcEndpoints.Invitation, Description: desc.Invitation},
			"site":       {Host: r.config.grpcEndpoints.Site, Description: desc.Site},
		}, r.config.grpcEndpoints.Timeout)
	}

	auth := r.f.Group("/v1", "API gateway", "Registry system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")
			return
		}
		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)
		err := f(ctx, r.config.auth.AuthAPIGW)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}
	})
	auth.Use()
	{
		const mem = "/members"
		member := auth.Group(mem, "Members", desc.Member)
		member.GET("", formatDoc("Get Members", "Get all members of an organization"), tonic.Handler(r.getMembersHandler, http.StatusOK))
		member.GET("/user/:user_id", formatDoc("Get Member", "Get member by user id"), tonic.Handler(r.getMemberByUserIdHandler, http.StatusOK))
		member.POST("", formatDoc("Add Member", "Add a new member to an organization"), tonic.Handler(r.postMemberHandler, http.StatusCreated))
		member.GET("/:member_id", formatDoc("Get Member", "Get a member of an organization"), tonic.Handler(r.getMemberHandler, http.StatusOK))
		member.PATCH("/:member_id", formatDoc("Update Member", "Update a member of an organization"), tonic.Handler(r.patchMemberHandler, http.StatusOK))
		member.DELETE("/:member_id", formatDoc("Remove Member", "Remove a member from an organization"), tonic.Handler(r.removeMemberHandler, http.StatusOK))

		// Invitation routes
		const inv = "/invitations"
		invitations := auth.Group(inv, "Invitations", desc.Invitation)
		invitations.POST("", formatDoc("Add Invitation", "Add a new invitation to an organization"), tonic.Handler(r.postInvitationHandler, http.StatusCreated))
		invitations.GET("/:invitation_id", formatDoc("Get Invitation", "Get a specific invitation"), tonic.Handler(r.getInvitationHandler, http.StatusOK))
		invitations.PATCH("/:invitation_id", formatDoc("Update Invitation", "Update a specific invitation"), tonic.Handler(r.patchInvitationHandler, http.StatusOK))
		invitations.DELETE("/:invitation_id", formatDoc("Remove Invitation", "Remove a invitation from an organization"), tonic.Handler(r.removeInvitationHandler, http.StatusOK))
		invitations.GET("/", formatDoc("Get Invitations", "Get all invitations of an organization"), tonic.Handler(r.getAllInvitationsHandler, http.StatusOK))
		invitations.GET("/user/:email", formatDoc("Get Invitations by email", "Get invitations by email"), tonic.Handler(r.getInvitationsByEmailHandler, http.StatusOK))

		// Network routes
		// Networks
		const net = "/networks"
		networks := auth.Group(net, "Networks", desc.Network)
		networks.GET("", formatDoc("Get Networks", "Get all Networks of an organization"), tonic.Handler(r.getNetworksHandler, http.StatusOK))
		networks.GET("/default", formatDoc("Get Default Network", "Get default Networks of an organization"), tonic.Handler(r.getDefaultNetworkHandler, http.StatusOK))
		networks.POST("", formatDoc("Add Network", "Add a new network to an organization"), tonic.Handler(r.postNetworkHandler, http.StatusCreated))
		networks.GET("/:net_id", formatDoc("Get Network", "Get a specific network"), tonic.Handler(r.getNetworkHandler, http.StatusOK))
		networks.PATCH("/:net_id", formatDoc("Set Network Default", "Set a specific network default"), tonic.Handler(r.setNetworkDefaultHandler, http.StatusOK))
		networks.DELETE("/:net_id", formatDoc("Remove Network", "Remove a network of an organization"), tonic.Handler(r.removeNetworkHandler, http.StatusOK))
		// Admins
		// Vendors

		// Sites

		const site = "/sites"
		sites := auth.Group(site, "Sites", desc.Site)
		sites.GET("", formatDoc("Get Sites", "Get all sites of a network"), tonic.Handler(r.getSitesHandler, http.StatusOK))
		sites.POST("", formatDoc("Add Site", "Add a new site to a network"), tonic.Handler(r.postSiteHandler, http.StatusCreated))
		sites.GET("/:site_id", formatDoc("Get Site", "Get a site of a network"), tonic.Handler(r.getSiteHandler, http.StatusOK))
		sites.PATCH("/:site_id", formatDoc("Update Site", "Update a site of a network"), tonic.Handler(r.updateSiteHandler, http.StatusOK))
		sites.DELETE("/:site_id", formatDoc("Remove Site", "Remove a site of a network"), tonic.Handler(r.removeSiteHandler, http.StatusOK))

		// Node routes
		const node = "/nodes"
		nodes := auth.Group(node, "Nodes", desc.Node)
		/** Deprecated: Use List API instead */
		nodes.GET("", formatDoc("Get Nodes", "Get all or free Nodes"), tonic.Handler(r.getNodes, http.StatusOK))
		nodes.GET("/state", formatDoc("Get Nodes by state", "Get all nodes by state"), tonic.Handler(r.getNodesByState, http.StatusOK))
		nodes.GET("/:node_id", formatDoc("Get Node", "Get a specific node"), tonic.Handler(r.getNodeHandler, http.StatusOK))
		nodes.GET("sites/:site_id", formatDoc("Get Nodes For Site", "Get all nodes of a site"), tonic.Handler(r.getSiteNodesHandler, http.StatusOK))
		nodes.GET("networks/:net_id", formatDoc("Get Nodes For Network", "Get all nodes of a network"), tonic.Handler(r.getNetworkNodesHandler, http.StatusOK))
		//
		nodes.GET("/list", formatDoc("List Nodes", "List all by filters"), tonic.Handler(r.list, http.StatusOK))
		nodes.POST("", formatDoc("Add Node", "Add a new Node to an organization"), tonic.Handler(r.postAddNodeHandler, http.StatusCreated))
		nodes.PUT("/:node_id", formatDoc("Update Node", "Update node name or state"), tonic.Handler(r.putUpdateNodeHandler, http.StatusOK))
		nodes.PATCH("/:node_id", formatDoc("Update Node State", "Update node state"), tonic.Handler(r.patchUpdateNodeStateHandler, http.StatusOK))
		nodes.DELETE("/:node_id", formatDoc("Delete Node", "Remove node from org"), tonic.Handler(r.deleteNodeHandler, http.StatusOK))
		nodes.POST("/:node_id/attach", formatDoc("Attach Node", "Group nodes"), tonic.Handler(r.postAttachedNodesHandler, http.StatusCreated))
		nodes.DELETE("/:node_id/attach", formatDoc("Dettach Node", "Move node out of group"), tonic.Handler(r.deleteAttachedNodeHandler, http.StatusOK))
		nodes.POST("/:node_id/sites", formatDoc("Add To Site", "Add node to site"), tonic.Handler(r.postNodeToSiteHandler, http.StatusCreated))
		nodes.DELETE("/:node_id/sites", formatDoc("Release From Site", "Release node from site"), tonic.Handler(r.deleteNodeFromSiteHandler, http.StatusOK))
	}
}

func (r *Router) getNetworkNodesHandler(c *gin.Context, req *GetNetworkNodesRequest) (*nodepb.GetByNetworkResponse, error) {
	return r.clients.Node.GetNetworkNodes(c.Request.Context(), req.NetworkId)
}

func (r *Router) getSiteNodesHandler(c *gin.Context, req *GetSiteNodesRequest) (*nodepb.GetBySiteResponse, error) {
	return r.clients.Node.GetSiteNodes(c.Request.Context(), req.SiteId)
}

func (r *Router) getNodes(c *gin.Context) (*nodepb.GetNodesResponse, error) {
	return r.clients.Node.GetNodes(c.Request.Context())
}
func (r *Router) list(c *gin.Context, req *ListNodesRequest) (*nodepb.ListResponse, error) {
	listReq := &nodepb.ListRequest{
		Type:         req.Type,
		SiteId:       req.SiteId,
		NodeId:       req.NodeId,
		NetworkId:    req.NetworkId,
		State:        Undefined,
		Connectivity: Undefined,
	}

	if req.State != "" {
		nodeState := ukama.ParseNodeState(req.State)
		listReq.State = cpb.NodeState(nodeState)
	}

	if req.Connectivity != "" {
		listReq.Connectivity = cpb.NodeConnectivity(ukama.ParseNodeConnectivity(req.Connectivity))
	}
	return r.clients.Node.List(c.Request.Context(), listReq)
}

func (r *Router) getNodesByState(c *gin.Context, req *GetNodesByStateRequest) (*nodepb.GetNodesResponse, error) {
	return r.clients.Node.GetNodesByState(c.Request.Context(), req.Connectivity, req.State)
}

func (r *Router) getNodeHandler(c *gin.Context, req *GetNodeRequest) (*nodepb.GetNodeResponse, error) {
	return r.clients.Node.GetNode(c.Request.Context(), req.NodeId)
}

func (r *Router) postAddNodeHandler(c *gin.Context, req *AddNodeRequest) (*nodepb.AddNodeResponse, error) {
	return r.clients.Node.AddNode(c.Request.Context(), req.NodeId, req.Name, req.State)
}

func (r *Router) postAttachedNodesHandler(c *gin.Context, req *AttachNodesRequest) (*nodepb.AttachNodesResponse, error) {
	return r.clients.Node.AttachNodes(c.Request.Context(), req.ParentNode, req.AmpNodeL, req.AmpNodeR)
}

func (r *Router) deleteAttachedNodeHandler(c *gin.Context, req *DetachNodeRequest) (*nodepb.DetachNodeResponse, error) {
	return r.clients.Node.DetachNode(c.Request.Context(), req.NodeId)
}

func (r *Router) putUpdateNodeHandler(c *gin.Context, req *UpdateNodeRequest) (*nodepb.UpdateNodeResponse, error) {
	return r.clients.Node.UpdateNode(c.Request.Context(), req.NodeId, req.Name, req.Latitude, req.Longitude)
}

func (r *Router) patchUpdateNodeStateHandler(c *gin.Context, req *UpdateNodeStateRequest) (*nodepb.UpdateNodeResponse, error) {
	return r.clients.Node.UpdateNodeState(c.Request.Context(), req.NodeId, req.State)
}

func (r *Router) postNodeToSiteHandler(c *gin.Context, req *AddNodeToSiteRequest) (*nodepb.AddNodeToSiteResponse, error) {
	return r.clients.Node.AddNodeToSite(c.Request.Context(), req.NodeId, req.NetworkId, req.SiteId)
}

func (r *Router) deleteNodeFromSiteHandler(c *gin.Context, req *ReleaseNodeFromSiteRequest) (*nodepb.ReleaseNodeFromSiteResponse, error) {
	return r.clients.Node.ReleaseNodeFromSite(c.Request.Context(), req.NodeId)
}

func (r *Router) deleteNodeHandler(c *gin.Context, req *DeleteNodeRequest) (*nodepb.DeleteNodeResponse, error) {
	return r.clients.Node.DeleteNode(c.Request.Context(), req.NodeId)
}

/* Member */
func (r *Router) getMembersHandler(c *gin.Context, req *GetMembersRequest) (*mpb.GetMembersResponse, error) {
	return r.clients.Member.GetMembers(c.Request.Context())
}

func (r *Router) getMemberByUserIdHandler(c *gin.Context, req *GetMemberByUserRequest) (*mpb.GetMemberByUserIdResponse, error) {
	return r.clients.Member.GetMemberByUserId(c.Request.Context(), req.UserId)
}

func (r *Router) getMemberHandler(c *gin.Context, req *GetMemberRequest) (*mpb.MemberResponse, error) {
	return r.clients.Member.GetMember(c.Request.Context(), req.MemberId)
}

func (r *Router) postMemberHandler(c *gin.Context, req *MemberRequest) (*mpb.MemberResponse, error) {
	return r.clients.Member.AddMember(c.Request.Context(), req.UserUuid, req.Role)
}

func (r *Router) patchMemberHandler(c *gin.Context, req *UpdateMemberRequest) error {
	return r.clients.Member.UpdateMember(c.Request.Context(), req.MemberId, req.IsDeactivated, req.Role)
}

func (r *Router) removeMemberHandler(c *gin.Context, req *RemoveMemberRequest) error {
	return r.clients.Member.RemoveMember(c.Request.Context(), req.MemberId)
}

// Network handlers

func (r *Router) setNetworkDefaultHandler(c *gin.Context, req *GetNetworkRequest) (*netpb.SetDefaultResponse, error) {
	return r.clients.Network.SetNetworkDefault(c.Request.Context(), req.NetworkId)
}

func (r *Router) getNetworkHandler(c *gin.Context, req *GetNetworkRequest) (*netpb.GetResponse, error) {
	return r.clients.Network.GetNetwork(c.Request.Context(), req.NetworkId)
}

func (r *Router) getNetworksHandler(c *gin.Context) (*netpb.GetNetworksResponse, error) {

	return r.clients.Network.GetNetworks(c.Request.Context())
}

func (r *Router) getDefaultNetworkHandler(c *gin.Context) (*netpb.GetDefaultResponse, error) {

	return r.clients.Network.GetDefault(c.Request.Context())
}

func (r *Router) postNetworkHandler(c *gin.Context, req *AddNetworkRequest) (*netpb.AddResponse, error) {
	return r.clients.Network.AddNetwork(c.Request.Context(), req.NetName, req.AllowedCountries, req.AllowedNetworks,
		req.Budget, req.Overdraft, req.TrafficPolicy, req.PaymentLinks)
}

func (r *Router) removeNetworkHandler(c *gin.Context, req *GetNetworkRequest) (*netpb.DeleteResponse, error) {
	return r.clients.Network.RemoveNetwork(c.Request.Context(), req.NetworkId)
}

func (r *Router) getSiteHandler(c *gin.Context, req *GetSiteRequest) (*sitepb.GetResponse, error) {
	return r.clients.Site.GetSite(c.Request.Context(), req.SiteId)
}

func (r *Router) getSitesHandler(c *gin.Context, req *GetSitesListRequest) (*sitepb.ListResponse, error) {
	return r.clients.Site.List(c.Request.Context(), req.NetworkId, req.IsDeactivated)

}

func (r *Router) updateSiteHandler(c *gin.Context, req *UpdateSiteRequest) (*sitepb.UpdateResponse, error) {
	return r.clients.Site.UpdateSite(c.Request.Context(),
		req.SiteId,
		req.Name,
	)
}

func (r *Router) removeSiteHandler(c *gin.Context, req *GetSiteRequest) (*sitepb.DeleteResponse, error) {
	return r.clients.Site.RemoveSite(c.Request.Context(), req.SiteId)
}

func (r *Router) postSiteHandler(c *gin.Context, req *AddSiteRequest) (*sitepb.AddResponse, error) {

	return r.clients.Site.AddSiteContext(c.Request.Context(), &sitepb.AddRequest{
		NetworkId: req.NetworkId, Name: req.Name, BackhaulId: req.BackhaulId,
		PowerId: req.PowerId, AccessId: req.AccessId, SwitchId: req.SwitchId,
		Location: req.Location, SpectrumId: req.SpectrumId, IsDeactivated: req.IsDeactivated,
		Latitude: req.Latitude, Longitude: req.Longitude, InstallDate: req.InstallDate,
	})
}

func (r *Router) postInvitationHandler(c *gin.Context, req *AddInvitationRequest) (*invpb.AddResponse, error) {
	return r.clients.Invitation.AddInvitation(c.Request.Context(), req.Name, strings.ToLower(req.Email), req.Role)
}

func (r *Router) getInvitationHandler(c *gin.Context, req *GetInvitationRequest) (*invpb.GetResponse, error) {
	return r.clients.Invitation.GetInvitationById(c.Request.Context(), req.InvitationId)
}

func (r *Router) patchInvitationHandler(c *gin.Context, req *UpdateInvitationRequest) (*invpb.UpdateStatusResponse, error) {
	return r.clients.Invitation.UpdateInvitation(c.Request.Context(), req.InvitationId, req.Status, strings.ToLower(req.Email))
}

func (r *Router) removeInvitationHandler(c *gin.Context, req *RemoveInvitationRequest) (*invpb.DeleteResponse, error) {
	return r.clients.Invitation.RemoveInvitation(c.Request.Context(), req.InvitationId)
}

func (r *Router) getAllInvitationsHandler(c *gin.Context) (*invpb.GetAllResponse, error) {
	return r.clients.Invitation.GetAllInvitations(c.Request.Context())
}

func (r *Router) getInvitationsByEmailHandler(c *gin.Context, req *GetInvitationsByEmailReq) (*invpb.GetByEmailResponse, error) {
	return r.clients.Invitation.GetInvitationsByEmail(c.Request.Context(), strings.ToLower(req.Email))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
