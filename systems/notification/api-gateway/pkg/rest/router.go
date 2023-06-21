package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/notification/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/client"
	emailPkg "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	"github.com/wI2L/fizz/openapi"
)

var REDIRECT_URI = "https://subscriber.dev.ukama.com/swagger/#/"

type Router struct {
	f       *fizz.Fizz
	clients *Clients
	config  *RouterConfig
}

type RouterConfig struct {
	debugMode  bool
	serverConf *rest.HttpConfig
	auth       *config.Auth
}

type Clients struct {
	m notification
}

type notification interface {
	SendEmail(*emailPkg.SendEmailRequest) (*emailPkg.SendEmailResponse, error)
	GetEmailById(*emailPkg.GetEmailByIdRequest) (*emailPkg.GetEmailByIdResponse, error)
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	c := &Clients{}
	var err error
	c.m, err = client.NewMailer(endpoints.Mailer, endpoints.Timeout)
	if err != nil {
		logrus.Fatalf("failed to create mailer client: %v", err)
	}

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
		serverConf: &svcConf.Server,
		debugMode:  svcConf.DebugMode,
		auth:       svcConf.Auth,
	}
}

func (rt *Router) Run() {
	logrus.Info("Listening on port ", rt.config.serverConf.Port)
	err := rt.f.Engine().Run(fmt.Sprint(":", rt.config.serverConf.Port))
	if err != nil {
		panic(err)
	}
}
func (r *Router) init(f func(*gin.Context, string) error) {
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName, version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")
	auth := r.f.Group("/v1", "Notification API GW ", "Notification system version v1", func(ctx *gin.Context) {
		if r.config.auth.BypassAuthMode {
			logrus.Info("Bypassing auth")
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
		mailer := auth.Group("/mailer", "Mailer", "Mailer")
		mailer.POST("/sendEmail", formatDoc("Send email notification", ""), tonic.Handler(r.sendEmailHandler, http.StatusOK))
		mailer.GET("/:mailer_id", formatDoc("Get email by id", ""), tonic.Handler(r.getEmailByIdHandler, http.StatusOK))

	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) sendEmailHandler(c *gin.Context, req *SendEmailReq) (message emailPkg.SendEmailResponse, err error) {
	payload := emailPkg.SendEmailRequest{
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
		Values:  req.Values,
	}
	res, err := r.clients.m.SendEmail(&payload)
	if err != nil {
		return emailPkg.SendEmailResponse{}, err
	}


	return emailPkg.SendEmailResponse{
		Message: res.Message,
	}, nil
}

func (r *Router) getEmailByIdHandler(c *gin.Context, req *GetEmailByIdReq) (message emailPkg.GetEmailByIdResponse, err error) {
	payload := emailPkg.GetEmailByIdRequest{
		MailId: req.MailerId,
	}
	res, err := r.clients.m.GetEmailById(&payload)
	if err != nil {
		return emailPkg.GetEmailByIdResponse{}, err
	}

	return emailPkg.GetEmailByIdResponse{
		MailId: res.MailId,
		To:      res.To,
		Subject: res.Subject,
		Body:    res.Body,		
	}, nil

}
