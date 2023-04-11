package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/auth/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

var SESSION_KEY = "ukama_session"
var REDIRECT_URI = "https://auth.dev.ukama.com/swagger/#/"

type Router struct {
	f              *fizz.Fizz
	config         *RouterConfig
	authRestClient *providers.AuthRestClient
}

type RouterConfig struct {
	debugMode  bool
	serverConf *rest.HttpConfig
	r          *rest.RestClient
	auth       *config.Auth
	s          *config.Service
}

func NewRouter(config *RouterConfig, authRestClient *providers.AuthRestClient) *Router {

	r := &Router{
		config:         config,
		authRestClient: authRestClient,
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
		r:          svcConf.R,
		s:          svcConf.Service,
		auth:       svcConf.Auth,
	}
}

func (rt *Router) Run() {
	logrus.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		logrus.Error(err)
	}
}

func (r *Router) init() {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect="+REDIRECT_URI)
	auth := r.f.Group("/v1", "Auth API GW", "Auth system version v1", func(ctx *gin.Context) {
		res, err := r.authRestClient.AuthenticateUser(ctx, r.config.auth.AuthAPIGW)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}
		if res.StatusCode() != http.StatusOK {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res.String())
			return
		}
	})
	auth.Use()
	{
		auth.GET("/whoami", formatDoc("Get user info", ""), tonic.Handler(r.getUserInfo, http.StatusOK))
		auth.GET("/auth", formatDoc("Authenticate user", ""), tonic.Handler(r.authenticate, http.StatusOK))
	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	opt := []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
	return opt
}

func (p *Router) getUserInfo(c *gin.Context) (*GetUserInfo, error) {
	sessionStr, err := pkg.GetSessionFromCookie(c, SESSION_KEY)
	if err != nil {
		return nil, err
	}

	res, err := pkg.GetUserBySession(sessionStr, p.config.r)
	if err != nil {
		return nil, err
	}

	return &GetUserInfo{
		Id:    res.Identity.Id,
		Name:  res.Identity.Traits.Name,
		Email: res.Identity.Traits.Email,
	}, nil
}

func (p *Router) authenticate(c *gin.Context) (*Authenticate, error) {
	sessionStr, err := pkg.GetSessionFromCookie(c, SESSION_KEY)
	if err != nil {
		return nil, err
	}
	res, err := pkg.GetUserBySession(sessionStr, p.config.r)
	if err != nil {
		return nil, err
	}

	if res.Identity.Id == "" {
		return &Authenticate{
			IsValidSession: false,
		}, nil
	}

	return &Authenticate{
		IsValidSession: true,
	}, nil
}
