package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/mailer/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/mailer/api-gateway/pkg"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

type Router struct {
	f      *fizz.Fizz
	config *RouterConfig
	client *Clients
}

type RouterConfig struct {
	debugMode  bool
	serverConf *rest.HttpConfig
	mailer     *config.Mailer
	s          *config.Service
}

type MailerManager interface {
	SendEmail(to string, message string) (SendEmailRes, error)
}

type Clients struct {
	ma MailerManager
}

func NewClientsSet(a MailerManager) *Clients {
	c := &Clients{}
	c.ma = a
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

func NewRouterConfig(svcConf *pkg.Config) *RouterConfig {
	return &RouterConfig{
		serverConf: &svcConf.Server,
		debugMode:  svcConf.DebugMode,
		s:          svcConf.Service,

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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, "http://localhost:8080")

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
	res, err := p.client.ma.SendEmail(req.To, req.Message)
	if err != nil {
		return nil, err
	}

	return &SendEmailRes{
		Message: res.Message,
	}, nil
}
