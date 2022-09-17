package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go/config"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/services/common/rest"
	"github.com/ukama/ukama/systems/init/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg/client"
)

const NodeIdParamName = "node"

type Router struct {
	f       *fizz.Fizz
	port    int
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
	l *client.Lookup
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.ServiceName, version.Version, r.config.debugMode)
	init := r.f.Group("/", "Init", "Init system")

	lookup := init.Group("lookup", "lookup", "looking for credentials")
	lookup.GET("/:org", nil, tonic.Handler(r.Handler, http.StatusOK))
	lookup.PUT("/:org", tonic.Handler(r.Handler, http.StatusCreated))
	lookup.GET("/:org/:node", tonic.Handler(r.Handler, http.StatusOK))
	lookup.PUT("/:org/:node", tonic.Handler(r.Handler, http.StatusCreated))
	lookup.DELTE("/:org/:node", tonic.Handler(r.Handler, http.StatusOK))
	lookup.GET("/:org/:system", tonic.Handler(r.Handler, http.StatusOK))
	lookup.PUT("/:org/:system", tonic.Handler(r.Handler, http.StatusCreated))
	lookup.DELTE("/:org/:system", tonic.Handler(r.Handler, http.StatusOK))
}

func (r *Router) Handler(c *gin.Context, req *GetNodeRequest) error {

}
