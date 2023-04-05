package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/auth/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg"
	"github.com/ukama/ukama/systems/common/auth"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

var SESSION_KEY = "ukama_session"

type Router struct {
	f *fizz.Fizz

	config *RouterConfig
}

type RouterConfig struct {
	debugMode  bool
	serverConf *rest.HttpConfig
	r          *rest.RestClient
	auth       *config.Auth
	s          *config.Service
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
		r:          svcConf.R,
		s:          svcConf.Service,
		auth:       svcConf.Auth,
	}
}

func (rt *Router) Run() {
	logrus.Info("Listening on port ", rt.config.serverConf.Port)

	rt.f.Generator().SetSecuritySchemes(map[string]*openapi.SecuritySchemeOrRef{
		"ukama_session": {
			SecurityScheme: &openapi.SecurityScheme{
				Type: "oauth2",
				In:   "header",
				Name: SESSION_KEY,
				Flows: &openapi.OAuthFlows{
					Implicit: &openapi.OAuthFlow{
						AuthorizationURL: rt.config.auth.AuthAppUrl + "?redirect=" + rt.config.s.Uri + "/swagger/#/",
					},
				},
			},
		},
	})

	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		logrus.Error(err)
	}
}

func (r *Router) init() {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode)
	v1 := r.f.Group("/v1", "Auth system", "Auth system version v1")

	v1.GET("/whoami", formatDoc("Get user info", ""), tonic.Handler(r.getUserInfo, http.StatusOK))
	v1.GET("/auth", formatDoc("Authenticate user", ""), tonic.Handler(r.authenticate, http.StatusOK))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	opt := []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
	return opt
}

func (p *Router) getUserInfo(c *gin.Context) (*GetUserInfo, error) {
	sessionStr, err := auth.GetSessionFromCookie(c, SESSION_KEY)
	if err != nil {
		return nil, err
	}

	res, err := auth.GetUserBySession(sessionStr, p.config.r)
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
	sessionStr, err := auth.GetSessionFromCookie(c, SESSION_KEY)
	if err != nil {
		return nil, err
	}
	res, err := auth.GetUserBySession(sessionStr, p.config.r)
	if err != nil {
		return nil, err
	}

	if res.Identity.Id == "" {
		return nil, errors.New("user not found")
	}

	return &Authenticate{
		IsValidSession: true,
	}, nil
}
