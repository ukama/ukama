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
	pborg "github.com/ukama/ukama/systems/registry/org/pb/gen"
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
}

type Clients struct {
	Registry registry
	User     *client.Users
}

type registry interface {
	GetOrg(orgName string) (*pborg.GetByNameResponse, error)
	GetOrgs(ownerUUID string) (*pborg.GetByOwnerResponse, error)
	AddOrg(orgName string, owner string, certificate string) (*pborg.AddResponse, error)
	GetMember(orgName string, userUUID string) (*pborg.MemberResponse, error)
	GetMembers(orgName string) (*pborg.GetMembersResponse, error)
	AddMember(orgName string, userUUID string) (*pborg.MemberResponse, error)
	RemoveMember(orgName string, userUUID string) error
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Registry = client.NewRegistry(endpoints.Network, endpoints.Org, endpoints.Timeout)
	c.User = client.NewUsers(endpoints.Users, endpoints.Timeout)
	return c
}

func NewRouter(clients *Clients, config *RouterConfig) *Router {
	r := &Router{
		clients: clients,
		config:  config,
	}

	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init()
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		metricsConfig: svcConf.Metrics,
		httpEndpoints: &svcConf.HttpServices,
		serverConf:    &svcConf.Server,
		debugMode:     svcConf.DebugMode,
	}
}

func (rt *Router) Run() {
	logrus.Info("Listening on port ", rt.config.serverConf.Port)

	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init() {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "API gateway", "Registry system version v1")

	// org handlers
	const org = "/orgs"
	orgs := v1.Group(org, "Orgs", "Operations on Orgs")
	orgs.GET("/", formatDoc("Get Orgs", "Get all organization owned by a user"), tonic.Handler(r.getOrgsHandler, http.StatusOK))
	orgs.POST("/", formatDoc("Add Org", "Add a new organization"), tonic.Handler(r.postOrgHandler, http.StatusCreated))
	orgs.GET("/:org", formatDoc("Get Org", "Get a specific organization"), tonic.Handler(r.getOrgHandler, http.StatusOK))
	// update org
	// Deactivate org
	// Delete org
	orgs.GET("/:org/members", formatDoc("Get Members", "Get all members of an organization"), tonic.Handler(r.getMembersHandler, http.StatusOK))
	orgs.POST("/:org/members", formatDoc("Add Member", "Add a new member to an organization"), tonic.Handler(r.postMemberHandler, http.StatusCreated))
	orgs.GET("/:org/members/:user_uuid", formatDoc("Get Member", "Get a member of an organization"), tonic.Handler(r.getMemberHandler, http.StatusOK))
	// update member
	// Deactivate member
	orgs.DELETE("/:org/members/:user_uuid", formatDoc("Remove Member", "Remove a member from an organization"), tonic.Handler(r.removeMemberHandler, http.StatusOK))

	// Users Handlers
	const user = "/users"
	users := v1.Group(user, "Users", "Operations on Users")
	users.POST("/", formatDoc("Add User", "Add a new User to the registry"), tonic.Handler(r.postUserHandler, http.StatusCreated))
	users.GET("/:user_uuid", formatDoc("Get User", "Get a specific user"), tonic.Handler(r.getUserHandler, http.StatusOK))
}

// Org handlers

func (r *Router) getOrgHandler(c *gin.Context, req *GetOrgRequest) (*pborg.GetByNameResponse, error) {
	return r.clients.Registry.GetOrg(c.Param("org"))
}

func (r *Router) getOrgsHandler(c *gin.Context) (*pborg.GetByOwnerResponse, error) {
	ownerUUID, ok := c.GetQuery("user_uuid")
	if !ok {
		return nil, &rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: "user_uuid is a mandatory query parameter"}
	}

	return r.clients.Registry.GetOrgs(ownerUUID)
}

func (r *Router) postOrgHandler(c *gin.Context, req *AddOrgRequest) (*pborg.AddResponse, error) {
	return r.clients.Registry.AddOrg(req.OrgName, req.Owner, req.Certificate)
}

func (r *Router) getMembersHandler(c *gin.Context, req *GetOrgRequest) (*pborg.GetMembersResponse, error) {
	return r.clients.Registry.GetMembers(c.Param("org"))
}

func (r *Router) getMemberHandler(c *gin.Context, req *GetMemberRequest) (*pborg.MemberResponse, error) {
	return r.clients.Registry.GetMember(c.Param("org"), c.Param("user_uuid"))
}

func (r *Router) postMemberHandler(c *gin.Context, req *MemberRequest) (*pborg.MemberResponse, error) {
	return r.clients.Registry.AddMember(req.OrgName, req.UserUUID)
}

func (r *Router) removeMemberHandler(c *gin.Context, req *GetMemberRequest) error {
	return r.clients.Registry.RemoveMember(c.Param("org"), c.Param("user_uuid"))
}

func (r *Router) getUserHandler(c *gin.Context, req *GetUserRequest) (*userspb.GetResponse, error) {
	return r.clients.User.Get(c.Param("user_uuid"), c.GetString(USER_ID_KEY))
}

func (r *Router) postUserHandler(c *gin.Context, req *AddUserRequest) (*userspb.AddResponse, error) {
	return r.clients.User.AddUser(&userspb.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	},
		c.GetString(USER_ID_KEY))
}

func (r *Router) deleteUserHandler(c *gin.Context, req *DeleteUserRequest) error {
	return r.clients.User.Delete(req.UserId, c.GetString(USER_ID_KEY))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
