package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/auth/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg"

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
	s          *config.Service
	k          string
}

type MailerManager interface {
	sendEmail(to string, templateName string, values map[string]any) (*oc.SuccessfulNativeLogin, error)
}

type Clients struct {
	au MailerManager
}

func NewClientsSet(a MailerManager) *Clients {
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

		k: k,
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
	v1 := r.f.Group("/v1", "Mailer API GW", "Mailer system version v1")
	v1.POST("/sendEmail", formatDoc("send email", ""), tonic.Handler(r.sendEmail, http.StatusOK))
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	opt := []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
	return opt
}

func (p *Router) sendEmail(c *gin.Context, req *SendEmailReq) (*SendEmailRes, error) {
	res, err := p.client.au.sendEmail(req.To, req.TemplateName, req.Values)
	if err != nil {
		return nil, err
	}

	return &SendEmailRes{
		Success: res.Success,
	}, nil
}
