package rest

import (
	"fmt"
	"net/http"

	"google.golang.org/protobuf/types/known/wrapperspb"

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
	GetOrg(orgName string) (*pborg.Organization, error)
	GetOrgs(ownerUUID string) (*pborg.GetByOwnerResponse, error)
	AddOrg(orgName string, owner string) (*pborg.Organization, error)
	IsAuthorized(userId string, org string) (bool, error)
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
	// const org = "/orgs/" + ":" + ORG_URL_PARAMETER
	const org = "/orgs"

	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "API gateway", "Registry system version v1")

	// org handler
	orgs := v1.Group(org, "Orgs", "Operations on Orgs")
	orgs.GET("/:org", formatDoc("Get Org", "Get a specific organization"), tonic.Handler(r.getOrgHandler, http.StatusOK))

	// network
}

func (r *Router) getOrgNameFromRoute(c *gin.Context) string {
	return c.Param("org")
}

// Org handlers

func (r *Router) getOrgHandler(c *gin.Context, req *GetOrgRequest) (*pborg.Organization, error) {
	orgName := r.getOrgNameFromRoute(c)

	return r.clients.Registry.GetOrg(orgName)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) deleteUserHandler(c *gin.Context, req *DeleteUserRequest) error {
	return r.clients.User.Delete(req.UserId, c.GetString(USER_ID_KEY))
}

func (r *Router) getUserHandler(c *gin.Context, req *GetUserRequest) (*userspb.GetResponse, error) {
	return r.clients.User.Get(req.UserId, c.GetString(USER_ID_KEY))
}

func boolToPbBool(data *bool) *wrapperspb.BoolValue {
	if data == nil {
		return nil
	}
	return &wrapperspb.BoolValue{Value: *data}
}
