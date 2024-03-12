package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/test-sys/node-gateway/cmd/version"
	"github.com/ukama/ukama/systems/test-sys/node-gateway/pkg"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

type Router struct {
	f      *fizz.Fizz
	config *RouterConfig
}

type RouterConfig struct {
	debugMode  bool
	serverConf *rest.HttpConfig
}

func NewRouter(config *RouterConfig) *Router {
	r := &Router{
		config: config,
	}
	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r.init()
	return r
}

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		serverConf: &svcConf.Server,
		debugMode:  svcConf.DebugMode,
	}
}

func (rt *Router) Run() {
	log.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}

func (r *Router) init() {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode,"")

	r.f.GET("/ping", formatDoc("Ping", "Returns 'pong' if the server is up"), tonic.Handler(r.pingHandler, http.StatusOK))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) pingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
