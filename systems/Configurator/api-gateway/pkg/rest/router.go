package rest

import (
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/configurator/api-gateway/cmd/version"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/configurator/api-gateway/pkg"
	"github.com/ukama/ukama/systems/configurator/api-gateway/pkg/client"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
	contPb "github.com/ukama/ukama/systems/configurator/controller/pb/gen"
)

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
	Controller controller
}

type controller interface {
	RestartSite(siteName, networkId string) (*contPb.RestartSiteResponse, error)
	RestartNode(nodeId string) (*contPb.RestartNodeResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Controller = client.NewController(endpoints.Controller, endpoints.Timeout)

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
	auth := r.f.Group("/v1", "API gateway", "Configurator system version v1", func(ctx *gin.Context) {
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

		const cont = "/controllers"
		controller := auth.Group(cont, "Controllers", "Operations on controllers")
		controller.POST("/restartSite", formatDoc("restart site in org", "restaring a site within an org "), tonic.Handler(r.postRestartSiteHandler, http.StatusCreated))
		controller.POST("/restartNode", formatDoc("restart node in network", "restaring a node within an network "), tonic.Handler(r.postRestartNodeHandler, http.StatusCreated))

	}
}

func (r *Router) postRestartNodeHandler(c *gin.Context, req *RestartNodeRequest) (*contPb.RestartNodeResponse, error) {
	return r.clients.Controller.RestartNode(req.NodeId)
}

func (r *Router) postRestartSiteHandler(c *gin.Context, req *RestartSiteRequest) (*contPb.RestartSiteResponse, error) {
	return r.clients.Controller.RestartSite(req.SiteName, req.NetworkId)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
