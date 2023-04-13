package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	ory "github.com/ory/client-go"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/auth/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

var SESSION_KEY = "ukama_session"

type Router struct {
	f      *fizz.Fizz
	config *RouterConfig
}

type RouterConfig struct {
	debugMode  bool
	serverConf *rest.HttpConfig
	r          *rest.RestClient
	auth       *config.Auth
	o          *ory.APIClient
	s          *config.Service
	k          string
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

func NewRouterConfig(svcConf *pkg.Config, oc *ory.APIClient, k string) *RouterConfig {
	return &RouterConfig{
		serverConf: &svcConf.Server,
		debugMode:  svcConf.DebugMode,
		r:          svcConf.R,
		s:          svcConf.Service,
		o:          oc,
		auth:       svcConf.Auth,
		k:          k,
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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect="+r.config.auth.AuthAPIGW+"/swagger/#/")
	v1 := r.f.Group("/v1", "Auth API GW", "Auth system version v1")

	v1.GET("/whoami", formatDoc("Get user info", ""), tonic.Handler(r.getUserInfo, http.StatusOK))
	v1.GET("/auth", formatDoc("Authenticate user", ""), tonic.Handler(r.authenticate, http.StatusOK))
	v1.POST("/login", formatDoc("Login user", ""), tonic.Handler(r.login, http.StatusOK))
	v1.GET("/getSession/:token", formatDoc("Get user session based on session token", ""), tonic.Handler(r.getSession, http.StatusOK))
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

	res, err := pkg.ValidateSession(sessionStr, p.config.o)

	if err != nil {
		return nil, err
	}

	user, err := pkg.GetUserTraitsFromSession(res)
	if err != nil {
		return nil, err
	}
	return &GetUserInfo{
		Id:         user.Id,
		Name:       user.Name,
		Email:      user.Email,
		Role:       user.Role,
		FirstVisit: user.FirstVisit,
	}, nil
}

func (p *Router) authenticate(c *gin.Context) (*ory.Session, error) {
	sessionStr, err := pkg.GetSessionFromCookie(c, SESSION_KEY)
	if err != nil {
		return nil, err
	}

	res, err := pkg.ValidateSession(sessionStr, p.config.o)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Router) login(c *gin.Context, req *LoginReq) (*LoginRes, error) {
	res, err := pkg.LoginUser(req.Email, req.Password, p.config.o)
	if err != nil {
		return nil, err
	}

	token, err := pkg.GenerateJWT(res.SessionToken, res.Session.GetExpiresAt().Format(http.TimeFormat), res.Session.GetAuthenticatedAt().Format(http.TimeFormat), p.config.k)
	if err != nil {
		return nil, err
	}
	return &LoginRes{
		Token: token,
	}, nil
}

func (p *Router) getSession(c *gin.Context, req *GetSessionReq) (*ory.Session, error) {
	session, err := pkg.GetSessionFromToken(c.Writer, req.Token, p.config.k)
	if err != nil {
		return nil, err
	}
	res, err := pkg.CheckSession(session.Session, p.config.o)
	if err != nil {
		return nil, err
	}
	return res, nil
}
