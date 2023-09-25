package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/api/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/wI2L/fizz/openapi"

	log "github.com/sirupsen/logrus"
)

var REDIRECT_URI = "https://subscriber.dev.ukama.com/swagger/#/"

type Router struct {
	f       *fizz.Fizz
	clients client.Client
	config  *RouterConfig
}

type RouterConfig struct {
	debugMode  bool
	serverConf *rest.HttpConfig
	auth       *config.Auth
}

func NewRouter(clients client.Client, config *RouterConfig, authfunc func(*gin.Context, string) error) *Router {
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
		serverConf: &svcConf.Server,
		debugMode:  svcConf.DebugMode,
		auth:       svcConf.Auth,
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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName,
		version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")

	auth := r.f.Group("/v1", "Ukama API GW ", "API system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			log.Info("Bypassing auth")

			return
		}

		s := fmt.Sprintf("%s, %s, %s", pkg.SystemName, ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Request.Header.Set("Meta", s)

		err := f(ctx, r.config.auth.AuthServerUrl)
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
		// networks routes
		networks := auth.Group("/networks", "Network", "Networks")
		networks.POST("", formatDoc("Create Network", "Create a new network"), tonic.Handler(r.postNetwork, http.StatusPartialContent))
		networks.GET("/:network_id", formatDoc("Get Network", "Get a specific network"), tonic.Handler(r.getNetwork, http.StatusOK))

		// package routes
		packages := auth.Group("/packages", "Package", "Packages")
		packages.POST("", formatDoc("Add Package", "Add a new package"), tonic.Handler(r.postPackage, http.StatusPartialContent))
		packages.GET("/:package_id", formatDoc("Get Package", "Get a specific package"), tonic.Handler(r.getPackage, http.StatusOK))

		// sims routes
		sims := auth.Group("/sims", "Sim", "sims")
		sims.POST("", formatDoc("Configure Sim", "Configure a new sim"), tonic.Handler(r.postSim, http.StatusPartialContent))
		sims.GET("/:iccid", formatDoc("Get Sim", "Get a specific sim"), tonic.Handler(r.getSim, http.StatusOK))
	}
}

func (r *Router) postNetwork(c *gin.Context, req *AddNetworkReq) (*client.NetworkInfo, error) {
	return r.clients.CreateNetwork(req.OrgName, req.NetName, req.AllowedCountries,
		req.AllowedNetworks, req.Budget, req.Overdraft, req.TrafficPolicy, req.PaymentLinks)
}

func (r *Router) getNetwork(c *gin.Context, req *GetNetworkReq) (*client.NetworkInfo, error) {
	return r.clients.GetNetwork(req.NetworkId)
}

func (r *Router) postPackage(c *gin.Context, req *AddPackageReq) (*client.PackageInfo, error) {
	return r.clients.AddPackage(req.Name, req.OrgId, req.OwnerId, req.From, req.To, req.BaserateId,
		req.Active, req.Flatrate, req.SmsVolume, req.VoiceVolume, req.DataVolume, req.VoiceUnit, req.DataUnit,
		req.SimType, req.Apn, req.Type, req.Markup, req.Amount, req.Overdraft, req.TrafficPolicy, req.Networks)
}

func (r *Router) getPackage(c *gin.Context, req *GetPackageReq) (*client.PackageInfo, error) {
	return r.clients.GetPackage(req.PackageId)
}

func (r *Router) postSim(c *gin.Context, req *AddSimReq) (*client.SimInfo, error) {
	return r.clients.ConfigureSim(req.SubscriberId, req.NetworkId,
		req.PackageId, req.SimType, req.SimToken)
}

func (r *Router) getSim(c *gin.Context, req *GetSimReq) (*client.SimInfo, error) {
	return r.clients.GetSim(req.Iccid)
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}
