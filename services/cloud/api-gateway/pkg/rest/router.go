package rest

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/ukama/ukama/services/common/rest"
	"github.com/wI2L/fizz/openapi"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukama/services/cloud/api-gateway/cmd/version"
	pb "github.com/ukama/ukama/services/cloud/registry/pb/gen"
	"github.com/ukama/ukama/services/common/config"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/services/cloud/api-gateway/pkg"
	"github.com/ukama/ukama/services/cloud/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	userspb "github.com/ukama/ukama/services/cloud/users/pb/gen"
)

const ORG_URL_PARAMETER = "org"

type Router struct {
	f              *fizz.Fizz
	authMiddleware AuthMiddleware
	clients        *Clients
	config         *RouterConfig
}

type RouterConfig struct {
	metricsConfig config.Metrics
	httpEndpoints *pkg.HttpEndpoints
	debugMode     bool
	serverConf    *rest.HttpConfig
}

type Clients struct {
	Registry *client.Registry
	User     *client.Users
}

type AuthMiddleware interface {
	IsAuthenticated(c *gin.Context)
	IsAuthorized(c *gin.Context)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Registry = client.NewRegistry(endpoints.Registry, endpoints.TimeoutSeconds)
	c.User = client.NewUsers(endpoints.Users, endpoints.TimeoutSeconds)
	return c
}

func NewRouter(
	authMiddleware AuthMiddleware,
	clients *Clients,
	config *RouterConfig) *Router {

	r := &Router{
		authMiddleware: authMiddleware,
		clients:        clients,
		config:         config,
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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.ServiceName, version.Version, r.config.debugMode)
	authorized := r.f.Group("/", "Authorization", "Requires authorization", r.authMiddleware.IsAuthenticated,
		r.authMiddleware.IsAuthorized)

	authorized.Use()
	{
		const org = "/orgs/" + ":" + ORG_URL_PARAMETER

		authorized.GET(org, []fizz.OperationOption{}, tonic.Handler(r.orgHandler, http.StatusOK))

		// metrics
		metricsProxy := r.getMetricsProxyHandler()

		authorized.GET(org+"/metrics/:metric", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "For request format refer to nodes/metrics/openapi.json Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
				info.Summary = "Get metrics for a node"
				info.ID = "GetMetrics"
			}}, metricsProxy)

		// registry
		nodes := authorized.Group(org+"/nodes", "Nodes", "Nodes operations")
		nodes.GET("", nil, tonic.Handler(r.getNodesHandler, http.StatusOK))
		nodes.PUT("/:node", nil, tonic.Handler(r.addOrUpdateNodeHandler, http.StatusOK))
		nodes.GET("/:node", nil, tonic.Handler(r.getNodeHandler, http.StatusOK))
		nodes.DELETE("/:node", nil, tonic.Handler(r.deleteNodeHandler, http.StatusOK))

		nodes.GET("/metrics/openapi.json", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Summary = "Get metrics endpoint Open API specification"
			}}, metricsProxy)

		nodes.GET("/:node/metrics/:metric", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "For request format refer to nodes/metrics/openapi.json Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
				info.Summary = "Get metrics for a node"
				info.ID = "GetMetrics"
			}}, metricsProxy)

		nodes.GET("/:node/metrics/:metric/:path", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "For request format refer to nodes/metrics/openapi.json Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
				info.Summary = "Get metrics for a node"
			}}, metricsProxy)

		nodes.GET("/metrics", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "Get list of metrics for a node. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
				info.Summary = "Get metrics list for a node"
				info.ID = "GetMetricsList"
			}}, metricsProxy)

		// user's management
		hss := authorized.Group(org+"/users", "Network Users", "Operations on network users and SIM cards"+
			"Do not confuse with organization users")
		hss.GET("", formatDoc("Get list of users", ""), tonic.Handler(r.getUsersHandler, http.StatusOK))
		hss.POST("", formatDoc("Create new user", ""), tonic.Handler(r.postUsersHandler, http.StatusCreated))
		hss.DELETE("/:user", formatDoc("Delete user", ""), tonic.Handler(r.deleteUserHandler, http.StatusOK))
		hss.GET("/:user", formatDoc("Get user info", ""), tonic.Handler(r.getUserHandler, http.StatusOK))
		hss.PATCH("/:user", formatDoc("Update user's information and deactivates user",
			"All fields are optional. User could be deactivated by setting isDeactivated flag to true. If a user is deactivated, all his SIM cards are purged. This operation is not recoverable"),
			tonic.Handler(r.updateUserHandler, http.StatusOK))
		hss.PUT("/:user/sims/:iccid/services", formatDoc("Enable or disable services for a SIM card", ""),
			tonic.Handler(r.setSimStatusHandler, http.StatusOK))
		hss.GET("/:user/sims/:iccid/qr", formatDoc("Get e-sim installation QR code", ""),
			tonic.Handler(r.getSimQrHandler, http.StatusOK))

	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

//
func (r *Router) getMetricsProxyHandler() gin.HandlerFunc {
	nodeMetricsUrl, err := url.Parse(r.config.httpEndpoints.NodeMetrics)
	if err != nil {
		logrus.Fatalf("Failed to parse node metrics endpoint: %s", err)
	}
	director := func(req *http.Request) {
		logrus.Infof("Request %s proxied", req.URL.String())
		req.URL.Scheme = nodeMetricsUrl.Scheme
		req.URL.Host = nodeMetricsUrl.Host

		idx := max(strings.Index(req.URL.Path, "/node"), strings.Index(req.URL.Path, "/orgs"))
		req.URL.Path = req.URL.Path[idx:]

		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	return gin.WrapH(proxy)
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func (r *Router) getOrgNameFromRoute(c *gin.Context) string {
	return c.Param("org")
}

func (r *Router) orgHandler(c *gin.Context) (*pb.Organization, error) {
	orgName := r.getOrgNameFromRoute(c)
	return r.clients.Registry.GetOrg(orgName)
}

func (r *Router) getNodesHandler(c *gin.Context) (*NodesList, error) {
	orgName := r.getOrgNameFromRoute(c)
	nl, err := r.clients.Registry.GetNodes(orgName)
	if err != nil {
		return nil, err
	}

	return MapNodesList(nl), nil
}

func (r *Router) getNodeHandler(c *gin.Context, req *GetNodeRequest) (*NodeExtended, error) {
	nodeName := c.Param("node")
	resp, err := r.clients.Registry.GetNode(nodeName)
	if err != nil {
		return nil, err
	}
	return mapExtendeNode(resp.Node), nil
}

func (r *Router) addOrUpdateNodeHandler(c *gin.Context, req *AddNodeRequest) (*pb.Node, error) {
	node, isCreated, err := r.clients.Registry.AddOrUpdate(req.OrgName, req.NodeId, req.NodeName)

	if isCreated {
		c.Status(http.StatusCreated)
	}

	return node, err
}

func (r *Router) deleteNodeHandler(c *gin.Context, req *DeleteNodeRequest) (*pb.DeleteNodeResponse, error) {
	return r.clients.Registry.DeleteNode(req.NodeId)

}

func (r *Router) getUsersHandler(c *gin.Context) (*userspb.ListResponse, error) {
	orgName := r.getOrgNameFromRoute(c)
	return r.clients.User.GetUsers(orgName, c.GetString(USER_ID_KEY))
}

func (r *Router) postUsersHandler(c *gin.Context, req *UserRequest) (*userspb.AddResponse, error) {
	return r.clients.User.AddUser(req.Org, &userspb.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	},
		req.SimToken,
		c.GetString(USER_ID_KEY))
}

func (r *Router) updateUserHandler(c *gin.Context, req *UpdateUserRequest) (*userspb.User, error) {
	if req.IsDeactivated {
		err := r.clients.User.DeactivateUser(req.UserId, c.GetString(USER_ID_KEY))
		if err != nil {
			return nil, err
		}
	}

	_, err := r.clients.User.UpdateUser(req.UserId, &userspb.UserAttributes{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}, c.GetString(USER_ID_KEY))

	if err != nil {
		return nil, err
	}

	resUser, err := r.clients.User.Get(req.UserId, c.GetString(USER_ID_KEY))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated user")
	}

	return resUser.GetUser(), nil
}

func (r *Router) deleteUserHandler(c *gin.Context, req *DeleteUserRequest) error {
	return r.clients.User.Delete(req.UserId, c.GetString(USER_ID_KEY))
}

func (r *Router) getUserHandler(c *gin.Context, req *GetUserRequest) (*userspb.GetResponse, error) {
	return r.clients.User.Get(req.UserId, c.GetString(USER_ID_KEY))
}

func (r *Router) setSimStatusHandler(c *gin.Context, req *SetSimStatusRequest) (*userspb.Sim, error) {
	return r.clients.User.SetSimStatus(&userspb.SetSimStatusRequest{
		Iccid:   req.Iccid,
		Carrier: simServicesToPbService(req.Carrier),
		Ukama:   simServicesToPbService(req.Ukama),
	}, c.GetString(USER_ID_KEY))
}

func (r *Router) getSimQrHandler(c *gin.Context, req *GetSimQrRequest) (*userspb.GetQrCodeResponse, error) {
	return r.clients.User.GetQr(req.Iccid, c.GetString(USER_ID_KEY))
}

func boolToPbBool(data *bool) *wrapperspb.BoolValue {
	if data == nil {
		return nil
	}
	return &wrapperspb.BoolValue{Value: *data}
}

func simServicesToPbService(data *SimServices) *userspb.SetSimStatusRequest_SetServices {
	if data == nil {
		return nil
	}

	return &userspb.SetSimStatusRequest_SetServices{
		Data:  boolToPbBool(data.Data),
		Voice: boolToPbBool(data.Voice),
		Sms:   boolToPbBool(data.Sms),
	}
}
