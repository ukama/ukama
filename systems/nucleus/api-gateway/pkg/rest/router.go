package rest

import (
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/wI2L/fizz/openapi"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/nucleus/api-gateway/cmd/version"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/nucleus/api-gateway/pkg"
	"github.com/ukama/ukama/systems/nucleus/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	orgpb "github.com/ukama/ukama/systems/nucleus/orgs/pb/gen"
	userspb "github.com/ukama/ukama/systems/nucleus/users/pb/gen"
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
}

type registry interface {
	AddOrg(orgName string, owner string, certificate string) (*orgpb.AddResponse, error)
	GetOrg(orgName string) (*orgpb.GetByNameResponse, error)
	GetOrgs(ownerUUID string) (*orgpb.GetByOwnerResponse, error)
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
	}
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
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	},
		c.GetString(USER_ID_KEY))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
