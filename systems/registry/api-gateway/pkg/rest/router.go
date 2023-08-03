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
	"github.com/sirupsen/logrus"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	userspb "github.com/ukama/ukama/systems/registry/users/pb/gen"
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
	Registry registry
	User     *client.Users
	Node     *client.Node
}

type registry interface {
	AddOrg(orgName string, owner string, certificate string) (*orgpb.AddResponse, error)
	GetOrg(orgName string) (*orgpb.GetByNameResponse, error)
	GetOrgs(ownerUUID string) (*orgpb.GetByOwnerResponse, error)
	GetMember(orgName string, userUUID string) (*orgpb.MemberResponse, error)
	GetMembers(orgName string) (*orgpb.GetMembersResponse, error)
	AddMember(orgName string, userUUID string, role string) (*orgpb.MemberResponse, error)
	UpdateMember(orgName string, userUUID string, isDeactivated bool, role string) error
	RemoveMember(orgName string, userUUID string) error

	AddNetwork(orgName string, netName string) (*netpb.AddResponse, error)
	GetNetwork(netID string) (*netpb.GetResponse, error)
	GetNetworks(org string) (*netpb.GetByOrgResponse, error)

	AddSite(netID string, siteName string) (*netpb.AddSiteResponse, error)
	GetSite(netID string, siteName string) (*netpb.GetSiteResponse, error)
	GetSites(netID string) (*netpb.GetSitesByNetworkResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Registry = client.NewRegistry(endpoints.Network, endpoints.Org, endpoints.Timeout)
	c.User = client.NewUsers(endpoints.Users, endpoints.Timeout)
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
	logrus.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
	auth := r.f.Group("/v1", "API gateway", "Registry system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			logrus.Info("Bypassing auth")
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
		const org = "/orgs"
		orgs := auth.Group(org, "Orgs", "Operations on Orgs")
		orgs.GET("", formatDoc("Get Orgs", "Get all organization owned by a user"), tonic.Handler(r.getOrgsHandler, http.StatusOK))
		orgs.POST("", formatDoc("Add Org", "Add a new organization"), tonic.Handler(r.postOrgHandler, http.StatusCreated))
		orgs.GET("/:org", formatDoc("Get Org", "Get a specific organization"), tonic.Handler(r.getOrgHandler, http.StatusOK))
		// update org
		// Deactivate org
		// Delete org
		orgs.GET("/:org/members", formatDoc("Get Members", "Get all members of an organization"), tonic.Handler(r.getMembersHandler, http.StatusOK))
		orgs.POST("/:org/members", formatDoc("Add Member", "Add a new member to an organization"), tonic.Handler(r.postMemberHandler, http.StatusCreated))
		orgs.GET("/:org/members/:user_uuid", formatDoc("Get Member", "Get a member of an organization"), tonic.Handler(r.getMemberHandler, http.StatusOK))
		orgs.PATCH("/:org/members/:user_uuid", formatDoc("Update Member", "Update a member of an organization"), tonic.Handler(r.patchMemberHandler, http.StatusOK))
		orgs.DELETE("/:org/members/:user_uuid", formatDoc("Remove Member", "Remove a member from an organization"), tonic.Handler(r.removeMemberHandler, http.StatusOK))
		orgs.GET("/:org/nodes", formatDoc("Get Org Nodes", "Get all or free nodes of an organization"), tonic.Handler(r.getOrgNodesHandler, http.StatusOK))

		// Users routes
		const user = "/users"
		users := auth.Group(user, "Users", "Operations on Users")
		users.POST("", formatDoc("Add User", "Add a new User to the registry"), tonic.Handler(r.postUserHandler, http.StatusCreated))
		users.GET("/:user_id", formatDoc("Get User", "Get a specific user"), tonic.Handler(r.getUserHandler, http.StatusOK))
		users.GET("/auth/:auth_id", formatDoc("Get User By AuthId", "Get a specific user by authId"), tonic.Handler(r.getUserByAuthIdHandler, http.StatusOK))
		users.GET("/whoami/:user_id", formatDoc("Get detailed User", "Get a specific user details with all linked orgs"), tonic.Handler(r.whoamiHandler, http.StatusOK))
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
	return r.clients.Node.AddNode(req.NodeId, req.Name, req.OrgId)
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

// Org handlers
func (r *Router) getOrgHandler(c *gin.Context, req *GetOrgRequest) (*orgpb.GetByNameResponse, error) {
	return r.clients.Registry.GetOrg(c.Param("org"))
}

func (r *Router) getOrgsHandler(c *gin.Context, req *GetOrgsRequest) (*orgpb.GetByOwnerResponse, error) {
	ownerUUID, ok := c.GetQuery("user_uuid")
	if !ok {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: "user_uuid is a mandatory query parameter"}
	}

	return r.clients.Registry.GetOrgs(ownerUUID)
}

func (r *Router) postOrgHandler(c *gin.Context, req *AddOrgRequest) (*orgpb.AddResponse, error) {
	return r.clients.Registry.AddOrg(req.OrgName, req.Owner, req.Certificate)
}

func (r *Router) getMembersHandler(c *gin.Context, req *GetOrgRequest) (*orgpb.GetMembersResponse, error) {
	return r.clients.Registry.GetMembers(c.Param("org"))
}

func (r *Router) getMemberHandler(c *gin.Context, req *GetMemberRequest) (*orgpb.MemberResponse, error) {
	return r.clients.Registry.GetMember(c.Param("org"), c.Param("user_uuid"))
}

func (r *Router) postMemberHandler(c *gin.Context, req *MemberRequest) (*orgpb.MemberResponse, error) {
	return r.clients.Registry.AddMember(req.OrgName, req.UserUuid, req.Role)
}

func (r *Router) patchMemberHandler(c *gin.Context, req *UpdateMemberRequest) error {
	return r.clients.Registry.UpdateMember(req.OrgName, req.UserUuid, req.IsDeactivated, req.Role)
}

func (r *Router) removeMemberHandler(c *gin.Context, req *GetMemberRequest) error {
	return r.clients.Registry.RemoveMember(c.Param("org"), c.Param("user_uuid"))
}

// Users handlers
func (r *Router) getUserHandler(c *gin.Context, req *GetUserRequest) (*userspb.GetResponse, error) {
	return r.clients.User.Get(c.Param("user_id"), c.GetString(USER_ID_KEY))
}

func (r *Router) getUserByAuthIdHandler(c *gin.Context, req *GetUserByAuthIdRequest) (*userspb.GetResponse, error) {
	return r.clients.User.GetByAuthId(c.Param("auth_id"), c.GetString(USER_ID_KEY))
}

func (r *Router) whoamiHandler(c *gin.Context, req *GetUserRequest) (*userspb.WhoamiResponse, error) {
	return r.clients.User.Whoami(c.Param("user_id"), c.GetString(USER_ID_KEY))
}

func (r *Router) postUserHandler(c *gin.Context, req *AddUserRequest) (*userspb.AddResponse, error) {
	return r.clients.User.AddUser(&userspb.User{
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		AuthId: req.AuthId,
	},
		c.GetString(USER_ID_KEY))
}

// Network handlers

func (r *Router) getNetworkHandler(c *gin.Context, req *GetNetworkRequest) (*netpb.GetResponse, error) {
	return r.clients.Registry.GetNetwork(req.NetworkId)
}

func (r *Router) getNetworksHandler(c *gin.Context, req *GetNetworksRequest) (*netpb.GetByOrgResponse, error) {
	orgName, ok := c.GetQuery("org")
	if !ok {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: "org is a mandatory query parameter"}
	}

	return r.clients.Registry.GetNetworks(orgName)
}

func (r *Router) postNetworkHandler(c *gin.Context, req *AddNetworkRequest) (*netpb.AddResponse, error) {
	return r.clients.Registry.AddNetwork(req.OrgName, req.NetName)
}

func (r *Router) getSiteHandler(c *gin.Context, req *GetSiteRequest) (*netpb.GetSiteResponse, error) {
	return r.clients.Registry.GetSite(req.NetworkId, req.SiteName)
}

func (r *Router) getSitesHandler(c *gin.Context, req *GetNetworkRequest) (*netpb.GetSitesByNetworkResponse, error) {
	return r.clients.Registry.GetSites(req.NetworkId)
}

func (r *Router) postSiteHandler(c *gin.Context, req *AddSiteRequest) (*netpb.AddSiteResponse, error) {
	return r.clients.Registry.AddSite(req.NetworkId, req.SiteName)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
