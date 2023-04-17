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
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg/client"

	oc "github.com/ory/client-go"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

var SESSION_KEY = "ukama_session"

type Router struct {
	f      *fizz.Fizz
	config *RouterConfig
	client *Clients
}

type RouterConfig struct {
	debugMode  bool
	serverConf *rest.HttpConfig
	auth       *config.Auth
	o          *ory.APIClient
	s          *config.Service
	k          string
}

type AuthManager interface {
	ValidateSession(ss, t string) (*oc.Session, error)
	LoginUser(email string, password string) (*oc.SuccessfulNativeLogin, error)
	UpdateRole(ss, t, orgId, role string, user *pkg.UserTraits) error
	AuthorizeUser(ss, t, orgId, role, relation, object string) (*oc.Session, error)
}

type Clients struct {
	au AuthManager
}

func NewClientsSet(a AuthManager) *Clients {
	c := &Clients{}
	c.au = a
	return c
}

func NewRouter(c *Clients, config *RouterConfig) *Router {
	r := &Router{
		config: config,
		client: c,
	}
	if !config.debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init()
	return r
}

func NewRouterConfig(svcConf *pkg.Config, k string) *RouterConfig {
	return &RouterConfig{
		serverConf: &svcConf.Server,
		debugMode:  svcConf.DebugMode,
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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
	v1 := r.f.Group("/v1", "Auth API GW", "Auth system version v1")

	v1.GET("/whoami", formatDoc("Get user info", ""), tonic.Handler(r.getUserInfo, http.StatusOK))
	v1.GET("/auth", formatDoc("Authenticate user", ""), tonic.Handler(r.authenticate, http.StatusOK))
	v1.POST("/login", formatDoc("Login user", ""), tonic.Handler(r.login, http.StatusOK))
	v1.PUT("/role", formatDoc("Update user role", ""), tonic.Handler(r.updateRole, http.StatusOK))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	opt := []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
	return opt
}

func (p *Router) getUserInfo(c *gin.Context, req *OptReqHeader) (*GetUserInfo, error) {
	st, err := pkg.SessionType(c, SESSION_KEY)
	if err != nil {
		return nil, err
	}
	var ss string

	if st == "cookie" {
		ss = pkg.GetCookieStr(c, SESSION_KEY)
	} else if st == "header" {
		ss = pkg.GetTokenStr(c)
		err := pkg.ValidateToken(c.Writer, ss, p.config.k)
		if err == nil {
			t, e := pkg.GetSessionFromToken(c.Writer, ss, p.config.k)
			if e != nil {
				return nil, e
			}
			ss = t.Session
		} else {
			return nil, err
		}
	}
	res, err := p.client.au.ValidateSession(ss, st)
	if err != nil {
		return nil, err
	}

	user, err := pkg.GetUserTraitsFromSession(req.OrgId, res)
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

func (p *Router) authenticate(c *gin.Context, req *OptReqHeader) error {
	st, err := pkg.SessionType(c, SESSION_KEY)
	if err != nil {
		return err
	}

	var ss string
	var orgId string
	if st == "cookie" {
		ss = pkg.GetCookieStr(c, SESSION_KEY)
	} else if st == "header" {
		_, orgId = pkg.GetMemberDetails(c)
		ss = pkg.GetTokenStr(c)
		err := pkg.ValidateToken(c.Writer, ss, p.config.k)

		if err == nil {
			t, e := pkg.GetSessionFromToken(c.Writer, ss, p.config.k)
			if e != nil {
				return e
			}
			ss = t.Session
		} else {
			return err
		}
	}
	meta := c.Request.Header.Get("meta")
	_, method, path, err := pkg.GetMetaHeaderValues(meta)
	if err != nil {
		return err
	}

	resp, err := p.client.au.ValidateSession(ss, st)
	if err != nil {
		return err
	}

	user, err := pkg.GetUserTraitsFromSession(req.OrgId, resp)
	if err != nil {
		return err
	}

	if user.Role != "" {
		_, err := p.client.au.AuthorizeUser(ss, st, user.Role, orgId, method, path)
		if err != nil {
			return err
		}

	} else {
		logrus.Errorf("No role found for organization %s", orgId)
		return errors.New("No role found for organization " + orgId)
	}
	logrus.Infof("user %s is %s in %s", user.Id, user.Role, orgId)

	return nil
}

func (p *Router) login(c *gin.Context, req *LoginReq) (*LoginRes, error) {
	res, err := p.client.au.LoginUser(req.Email, req.Password)
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

func (p *Router) updateRole(c *gin.Context, req *UpdateRoleReq) error {

	st, err := pkg.SessionType(c, SESSION_KEY)
	if err != nil {
		return err
	}
	var ss string
	if st == "cookie" {
		ss = pkg.GetCookieStr(c, SESSION_KEY)
	} else if st == "header" {
		ss = pkg.GetTokenStr(c)
		err := pkg.ValidateToken(c.Writer, ss, p.config.k)

		if err == nil {
			t, e := pkg.GetSessionFromToken(c.Writer, ss, p.config.k)
			if e != nil {
				return e
			}
			ss = t.Session
		} else {
			return err
		}
	}

	logrus.Infof("fetch user session %s %s", ss, st)
	res, err := p.client.au.ValidateSession(ss, st)
	if err != nil {
		return err
	}
	logrus.Info("parse response")
	user, err := pkg.GetUserTraitsFromSession(req.OrgId, res)
	if err != nil {
		return err
	}
	logrus.Infof("update role of user %s", user.Id)

	err = p.client.au.UpdateRole(ss, st, req.OrgId, string(req.Role), user)
	if err != nil {
		return err
	}

	return nil
}
