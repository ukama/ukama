package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/api/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/db"
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
		// network routes
		network := auth.Group("/networks", "Network", "Networks")
		network.POST("", formatDoc("Creste Network", "Create a new network"), tonic.Handler(r.postNetwork, http.StatusCreated))
		network.GET("/:network_id", formatDoc("Get Network", "Get a specific network"), tonic.Handler(r.getNetwork, http.StatusOK))
	}
}

func (r *Router) postNetwork(c *gin.Context, req *AddNetworkReq) (*client.NetworkInfo, error) {
	res, status, err := r.clients.CreateNetwork(req.OrgName, req.NetName)
	if err != nil {
		if !errors.Is(err, &client.ErrorResponse{}) {
			return nil, rest.HttpError{
				HttpCode: http.StatusRequestTimeout,
				Message:  err.Error(),
			}
		}

		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	c.Status(resourceStatusToHTTPStatus(*status))

	return res, err
}

// func (r *Router) getNetwork(c *gin.Context, req *GetNetworkReq) (*npb.GetResponse, error) {
func (r *Router) getNetwork(c *gin.Context, req *GetNetworkReq) (*client.NetworkInfo, error) {
	res, status, err := r.clients.GetNetwork(req.NetworkId)
	if err != nil {
		if !errors.Is(err, &client.ErrorResponse{}) {
			return nil, rest.HttpError{
				HttpCode: http.StatusRequestTimeout,
				Message:  err.Error(),
			}
		}

		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  err.Error(),
		}
	}

	c.Status(resourceStatusToHTTPStatus(*status))

	return res, err
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func resourceStatusToHTTPStatus(status db.ResourceStatus) int {
	switch status {
	case db.ResourceStatusPending:
		return http.StatusPartialContent

	case db.ResourceStatusCompleted:
		return http.StatusOK

	default:
		return http.StatusUnprocessableEntity
	}
}
