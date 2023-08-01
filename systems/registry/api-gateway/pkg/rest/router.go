package rest

import (
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/wI2L/fizz/openapi"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/registry/api-gateway/cmd/version"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/registry/api-gateway/pkg"
	"github.com/ukama/ukama/systems/registry/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
	mpb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
)

const USER_ID_KEY = "UserId"
const ORG_URL_PARAMETER = "org"

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	metricsConfig config.Metrics
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
	auth          *config.Auth
}

type Clients struct {
	Network network
	Node    node
	Member  member
}

type network interface {
	AddNetwork(orgName string, netName string) (*netpb.AddResponse, error)
	GetNetwork(netID string) (*netpb.GetResponse, error)
	GetNetworks(org string) (*netpb.GetByOrgResponse, error)

	AddSite(netID string, siteName string) (*netpb.AddSiteResponse, error)
	GetSite(netID string, siteName string) (*netpb.GetSiteResponse, error)
	GetSites(netID string) (*netpb.GetSitesByNetworkResponse, error)
}

type member interface {
	GetMember(userUUID string) (*mpb.MemberResponse, error)
	GetMembers() (*mpb.GetMembersResponse, error)
	AddMember(userUUID string, role string) (*mpb.MemberResponse, error)
	UpdateMember(userUUID string, isDeactivated bool, role string) error
	RemoveMember(userUUID string) error
}

type node interface {
	AddNode(nodeId, name, orgId, state string) (*nodepb.AddNodeResponse, error)
	GetNode(nodeId string) (*nodepb.GetNodeResponse, error)
	GetOrgNodes(orgId string, free bool) (*nodepb.GetByOrgResponse, error)
	GetSiteNodes(siteId string) (*nodepb.GetBySiteResponse, error)
	GetAllNodes(free bool) (*nodepb.GetNodesResponse, error)
	UpdateNodeState(nodeId string, state string) (*nodepb.UpdateNodeResponse, error)
	UpdateNode(nodeId string, name string) (*nodepb.UpdateNodeResponse, error)
	DeleteNode(nodeId string) (*nodepb.DeleteNodeResponse, error)
	AttachNodes(node, l, r string) (*nodepb.AttachNodesResponse, error)
	DetachNode(nodeId string) (*nodepb.DetachNodeResponse, error)
	AddNodeToSite(nodeId, networkId, siteId string) (*nodepb.AddNodeToSiteResponse, error)
	ReleaseNodeFromSite(nodeId string) (*nodepb.ReleaseNodeFromSiteResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Network = client.NewNetworkRegistry(endpoints.Network, endpoints.Timeout)
	c.Node = client.NewNode(endpoints.Node, endpoints.Timeout)
	c.Member = client.NewMemberRegistry(endpoints.Member, endpoints.Timeout)

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
		httpEndpoints: &svcConf.HttpServices,
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
			return
		}
		if err == nil {
			return
		}
	})
	auth.Use()
	{
		// org routes
		// const org = "/orgs"
		// orgs := auth.Group(org, "Orgs", "Operations on Orgs")
		// orgs.GET("", formatDoc("Get Orgs", "Get all organization owned by a user"), tonic.Handler(r.getOrgsHandler, http.StatusOK))
		// orgs.POST("", formatDoc("Add Org", "Add a new organization"), tonic.Handler(r.postOrgHandler, http.StatusCreated))
		// orgs.GET("/:org", formatDoc("Get Org", "Get a specific organization"), tonic.Handler(r.getOrgHandler, http.StatusOK))
		// update org
		// Deactivate org
		// Delete org
		const mem = "/member"
		member := auth.Group(mem, "Members", "Operations on Members")
		member.GET("", formatDoc("Get Members", "Get all members of an organization"), tonic.Handler(r.getMembersHandler, http.StatusOK))
		member.POST("", formatDoc("Add Member", "Add a new member to an organization"), tonic.Handler(r.postMemberHandler, http.StatusCreated))
		member.GET("/:user_uuid", formatDoc("Get Member", "Get a member of an organization"), tonic.Handler(r.getMemberHandler, http.StatusOK))
		member.PATCH("/:user_uuid", formatDoc("Update Member", "Update a member of an organization"), tonic.Handler(r.patchMemberHandler, http.StatusOK))
		member.DELETE("/:user_uuid", formatDoc("Remove Member", "Remove a member from an organization"), tonic.Handler(r.removeMemberHandler, http.StatusOK))

		// Users routes
		// const user = "/users"
		// users := auth.Group(user, "Users", "Operations on Users")
		// users.POST("", formatDoc("Add User", "Add a new User to the registry"), tonic.Handler(r.postUserHandler, http.StatusCreated))
		// users.GET("/:user_id", formatDoc("Get User", "Get a specific user"), tonic.Handler(r.getUserHandler, http.StatusOK))
		// users.GET("/auth/:auth_id", formatDoc("Get User By AuthId", "Get a specific user by authId"), tonic.Handler(r.getUserByAuthIdHandler, http.StatusOK))
		// users.GET("/whoami/:user_id", formatDoc("Get detailed User", "Get a specific user details with all linked orgs"), tonic.Handler(r.whoamiHandler, http.StatusOK))
		// user orgs-member
		// update user
		// Deactivate user
		// Delete user
		// users.DELETE("/:user_id", formatDoc("Remove User", "Remove a user from the registry"), tonic.Handler(r.removeUserHandler, http.StatusOK))

		// Network routes
		// Networks
		const net = "/networks"
		networks := auth.Group(net, "Networks", "Operations on Networks")
		networks.GET("", formatDoc("Get Networks", "Get all Networks of an organization"), tonic.Handler(r.getNetworksHandler, http.StatusOK))
		networks.POST("", formatDoc("Add Network", "Add a new network to an organization"), tonic.Handler(r.postNetworkHandler, http.StatusCreated))
		networks.GET("/:net_id", formatDoc("Get Network", "Get a specific network"), tonic.Handler(r.getNetworkHandler, http.StatusOK))
		// update network
		// networks.DELETE("/:net_id", formatDoc("Remove Network", "Remove a network of an organization"), tonic.Handler(r.removeNetworkHandler, http.StatusOK))
		// Admins
		// Vendors

		// Sites
		networks.GET("/:net_id/sites", formatDoc("Get Sites", "Get all sites of a network"), tonic.Handler(r.getSitesHandler, http.StatusOK))
		networks.POST("/:net_id/sites", formatDoc("Add Site", "Add a new site to a network"), tonic.Handler(r.postSiteHandler, http.StatusCreated))
		networks.GET("/:net_id/sites/:site", formatDoc("Get Site", "Get a site of a network"), tonic.Handler(r.getSiteHandler, http.StatusOK))
		// update sites
		// delete sites

		// Node routes
		const node = "/nodes"
		nodes := auth.Group(node, "Nodes", "Operations on Nodes")
		nodes.GET("", formatDoc("Get Nodes", "Get all or free Nodes"), tonic.Handler(r.getAllNodesHandler, http.StatusOK))
		nodes.GET("/:node_id", formatDoc("Get Node", "Get a specific node"), tonic.Handler(r.getNodeHandler, http.StatusOK))
		nodes.GET("sites/:site_id", formatDoc("Get Nodes For Site", "Get all nodes of a site"), tonic.Handler(r.getSiteNodesHandler, http.StatusOK))
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

// Node handlers
func (r *Router) getOrgNodesHandler(c *gin.Context, req *GetOrgNodesRequest) (*nodepb.GetByOrgResponse, error) {
	return r.clients.Node.GetOrgNodes(req.OrgId, req.Free)
}

func (r *Router) getSiteNodesHandler(c *gin.Context, req *GetSiteNodesRequest) (*nodepb.GetBySiteResponse, error) {
	return r.clients.Node.GetSiteNodes(req.SiteId)
}

func (r *Router) getAllNodesHandler(c *gin.Context, req *GetNodesRequest) (*nodepb.GetNodesResponse, error) {
	return r.clients.Node.GetAllNodes(req.Free)
}

func (r *Router) getNodeHandler(c *gin.Context, req *GetNodeRequest) (*nodepb.GetNodeResponse, error) {
	return r.clients.Node.GetNode(req.NodeId)
}

func (r *Router) postAddNodeHandler(c *gin.Context, req *AddNodeRequest) (*nodepb.AddNodeResponse, error) {
	return r.clients.Node.AddNode(req.NodeId, req.Name, req.OrgId, req.State)
}

func (r *Router) postAttachedNodesHandler(c *gin.Context, req *AttachNodesRequest) (*nodepb.AttachNodesResponse, error) {
	return r.clients.Node.AttachNodes(req.ParentNode, req.AmpNodeL, req.AmpNodeR)
}

func (r *Router) deleteAttachedNodeHandler(c *gin.Context, req *DetachNodeRequest) (*nodepb.DetachNodeResponse, error) {
	return r.clients.Node.DetachNode(req.NodeId)
}

func (r *Router) putUpdateNodeHandler(c *gin.Context, req *UpdateNodeRequest) (*nodepb.UpdateNodeResponse, error) {
	return r.clients.Node.UpdateNode(req.NodeId, req.Name)
}

func (r *Router) patchUpdateNodeStateHandler(c *gin.Context, req *UpdateNodeStateRequest) (*nodepb.UpdateNodeResponse, error) {
	return r.clients.Node.UpdateNodeState(req.NodeId, req.State)
}

func (r *Router) postNodeToSiteHandler(c *gin.Context, req *AddNodeToSiteRequest) (*nodepb.AddNodeToSiteResponse, error) {
	return r.clients.Node.AddNodeToSite(req.NodeId, req.NetworkId, req.SiteId)
}

func (r *Router) deleteNodeFromSiteHandler(c *gin.Context, req *ReleaseNodeFromSiteRequest) (*nodepb.ReleaseNodeFromSiteResponse, error) {
	return r.clients.Node.ReleaseNodeFromSite(req.NodeId)
}

func (r *Router) deleteNodeHandler(c *gin.Context, req *DeleteNodeRequest) (*nodepb.DeleteNodeResponse, error) {
	return r.clients.Node.DeleteNode(req.NodeId)
}

/* Member */
func (r *Router) getMembersHandler(c *gin.Context, req *GetMembersRequest) (*mpb.GetMembersResponse, error) {
	return r.clients.Member.GetMembers()
}

func (r *Router) getMemberHandler(c *gin.Context, req *GetMemberRequest) (*mpb.MemberResponse, error) {
	return r.clients.Member.GetMember(req.UserUuid)
}

func (r *Router) postMemberHandler(c *gin.Context, req *MemberRequest) (*mpb.MemberResponse, error) {
	return r.clients.Member.AddMember(req.UserUuid, req.Role)
}

func (r *Router) patchMemberHandler(c *gin.Context, req *UpdateMemberRequest) error {
	return r.clients.Member.UpdateMember(req.UserUuid, req.IsDeactivated, req.Role)
}

func (r *Router) removeMemberHandler(c *gin.Context, req *RemoveMemberRequest) error {
	return r.clients.Member.RemoveMember(req.UserUuid)
}

// Network handlers

func (r *Router) getNetworkHandler(c *gin.Context, req *GetNetworkRequest) (*netpb.GetResponse, error) {
	return r.clients.Network.GetNetwork(req.NetworkId)
}

func (r *Router) getNetworksHandler(c *gin.Context, req *GetNetworksRequest) (*netpb.GetByOrgResponse, error) {
	orgName, ok := c.GetQuery("org")
	if !ok {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: "org is a mandatory query parameter"}
	}

	return r.clients.Network.GetNetworks(orgName)
}

func (r *Router) postNetworkHandler(c *gin.Context, req *AddNetworkRequest) (*netpb.AddResponse, error) {
	return r.clients.Network.AddNetwork(req.OrgName, req.NetName)
}

func (r *Router) getSiteHandler(c *gin.Context, req *GetSiteRequest) (*netpb.GetSiteResponse, error) {
	return r.clients.Network.GetSite(req.NetworkId, req.SiteName)
}

func (r *Router) getSitesHandler(c *gin.Context, req *GetNetworkRequest) (*netpb.GetSitesByNetworkResponse, error) {
	return r.clients.Network.GetSites(req.NetworkId)
}

func (r *Router) postSiteHandler(c *gin.Context, req *AddSiteRequest) (*netpb.AddSiteResponse, error) {
	return r.clients.Network.AddSite(req.NetworkId, req.SiteName)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
