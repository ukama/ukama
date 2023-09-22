package rest

import (
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"

	"github.com/loopfz/gadgeto/tonic"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/node/node-gateway/cmd/version"
	"github.com/ukama/ukama/systems/node/node-gateway/pkg/client"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"

	"github.com/ukama/ukama/systems/node/node-gateway/pkg"

	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/node/health/pb/gen"
	healthPb "github.com/ukama/ukama/systems/node/health/pb/gen"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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
	Health health
}


type health interface {
	StoreRunningAppsInfo(req *healthPb.StoreRunningAppsInfoRequest) (*healthPb.StoreRunningAppsInfoResponse, error)
	GetRunningAppsInfo(nodeId string) (*healthPb.GetRunningAppsResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	c.Health = client.NewHealth(endpoints.Health, endpoints.Timeout)

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
	auth := r.f.Group("/v1", "API gateway", "node system version v1", func(ctx *gin.Context) {
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

		const heal = "/health"
		health := auth.Group(heal, "health", "Operations on health")
		 health.POST("/:node_id", formatDoc("Store running sofwtare infos", ""), tonic.Handler(r.postNodeInfosHandler, http.StatusCreated))
		health.GET("/:node_id", formatDoc("Get running sofwtare infos", ""), tonic.Handler(r.getNodeInfosHandler, http.StatusOK))
	}
}


func (r *Router) postNodeInfosHandler(c *gin.Context, req *StoreRunningAppsInfoRequest) (*healthPb.StoreRunningAppsInfoResponse, error) {
    var genSystems []*gen.System
    for _, sys := range req.System {
        genSystem := &gen.System{
            Name:  sys.Name,
            Value: sys.Value,
        }
        genSystems = append(genSystems, genSystem)
    }
    return r.clients.Health.StoreRunningAppsInfo(&healthPb.StoreRunningAppsInfoRequest{
        NodeId:    req.NodeId,
        Timestamp: req.Timestamp,
        System:    genSystems,
    })
}




func (r *Router) getNodeInfosHandler(c *gin.Context, req *GetRunningAppsRequest) (*healthPb.GetRunningAppsResponse, error) {
	resp, err := r.clients.Health.GetRunningAppsInfo(req.NodeId)
	
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resp, nil
}
func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

