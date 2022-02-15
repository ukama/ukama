package rest

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ukama/ukamaX/common/rest"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukamaX/cloud/api-gateway/cmd/version"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"github.com/ukama/ukamaX/common/config"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukamaX/cloud/api-gateway/pkg"
	"github.com/ukama/ukamaX/cloud/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	hsspb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	nodeMetr "github.com/ukama/ukamaX/metrics/node-metrics"
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
	Hss      *client.Hss
}

type AuthMiddleware interface {
	IsAuthenticated(c *gin.Context)
	IsAuthorized(c *gin.Context)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Registry = client.NewRegistry(endpoints.Registry, endpoints.TimeoutSeconds)
	c.Hss = client.NewHss(endpoints.Hss, endpoints.TimeoutSeconds)
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
	r.init(config.serverConf.Port)
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

func (r *Router) init(port int) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.ServiceName, version.Version, r.config.debugMode)

	authorized := r.f.Group("/", "Authorization", "Requires authorization", r.authMiddleware.IsAuthenticated,
		r.authMiddleware.IsAuthorized)

	authorized.Use()
	{
		const org = "/orgs/" + ":" + ORG_URL_PARAMETER

		authorized.GET(org, []fizz.OperationOption{}, tonic.Handler(r.orgHandler, http.StatusOK))

		// registry
		nodes := authorized.Group(org+"/nodes", "Nodes", "Nodes operations")
		nodes.GET("", nil, tonic.Handler(r.nodesHandler, http.StatusOK))

		// metrics
		nodes.GET("/:node/metrics/:metric", []fizz.OperationOption{
			func(info *openapi.OperationInfo) {
				info.Description = "Get metrics for a node. Response has Prometheus data format https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors"
			}}, tonic.Handler(r.metricHandler, http.StatusOK))
		nodes.GET("/:node/metrics/list", nil, tonic.Handler(r.metricListHandler, http.StatusOK))

		// hss
		hss := authorized.Group(org+"/users", "Network Users", "Operations on network users and SIM cards"+
			"Do not confuse with organization users")
		hss.GET("", nil, tonic.Handler(r.getUsersHandler, http.StatusOK))
		hss.POST("", []fizz.OperationOption{}, tonic.Handler(r.postUsersHandler, http.StatusCreated))
		hss.DELETE("/:user", nil, tonic.Handler(r.deleteUserHandler, http.StatusOK))
	}
}

func (r *Router) getOrgNameFromRoute(c *gin.Context) string {
	return c.Param("org")
}

func (r *Router) orgHandler(c *gin.Context) (*pb.Organization, error) {
	orgName := r.getOrgNameFromRoute(c)
	return r.clients.Registry.GetOrg(orgName)
}

func (r *Router) nodesHandler(c *gin.Context) (*NodesList, error) {
	orgName := r.getOrgNameFromRoute(c)
	nl, err := r.clients.Registry.GetNodes(orgName)
	if err != nil {
		return nil, err
	}

	return MapNodesList(nl), nil
}

func (r *Router) metricListHandler(c *gin.Context) ([]string, error) {
	return nodeMetr.MetricTypes, nil
}

func (r *Router) metricHandler(c *gin.Context, in *GetNodeMetricsInput) error {
	exist := false
	for _, m := range nodeMetr.MetricTypes {
		if strings.EqualFold(m, in.Metric) {
			exist = true
		}
	}
	if !exist {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  "Metric not found"}
	}

	metr := nodeMetr.Metrics{
		PrometheusUrl: r.config.httpEndpoints.NodeMetrics,
		Timeout:       uint(r.config.httpEndpoints.TimeoutSeconds),
	}

	c.Stream(func(w io.Writer) bool {
		httpCode, err := metr.GetMetric(strings.ToLower(in.Metric), in.NodeID, &nodeMetr.Interval{
			Start: in.From,
			End:   in.To,
			Step:  in.Step,
		}, w)

		c.Status(httpCode)
		if err != nil {
			logrus.Errorf("Error while getting metrics: %v", err)
			return true
		}

		return false
	})

	return nil
}

type PingResponse struct {
	Message string `json:"message"`
}

func (r *Router) getUsersHandler(c *gin.Context) (*hsspb.ListUsersResponse, error) {
	orgName := r.getOrgNameFromRoute(c)
	return r.clients.Hss.GetUsers(orgName)
}

func (r *Router) postUsersHandler(c *gin.Context, req *UserRequest) (*hsspb.AddUserResponse, error) {
	return r.clients.Hss.AddUser(req.Org, &hsspb.User{
		Imsi:      req.Imsi,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	})
}

func (r *Router) deleteUserHandler(c *gin.Context) error {
	orgName := r.getOrgNameFromRoute(c)
	userId := c.Param("user")
	return r.clients.Hss.Delete(orgName, userId)
}
