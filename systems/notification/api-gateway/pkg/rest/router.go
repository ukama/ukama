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
	"github.com/wI2L/fizz/openapi"

	emailPkg "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	npb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
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
	m client.Mailer
	n client.Notify
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
	r.f = rest.NewFizzRouter(r.config.serverConf, pkg.SystemName,
		version.Version, r.config.debugMode, r.config.auth.AuthAppUrl+"?redirect=true")

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
		// mailer routes
		mailer := auth.Group("/mailer", "Mailer", "Mailer")
		mailer.POST("/sendEmail", formatDoc("Send email notification", ""), tonic.Handler(r.sendEmailHandler, http.StatusOK))
		mailer.GET("/:mailer_id", formatDoc("Get email by id", ""), tonic.Handler(r.getEmailByIdHandler, http.StatusOK))

		// notify routes
		notif := auth.Group("/notifications", "Notification", "Notifications")
		notif.POST("", formatDoc("Insert Notification", "Insert a new notification"), tonic.Handler(r.postNotification, http.StatusCreated))
		notif.GET("", formatDoc("Get Notifications", "Get a list of notifications"), tonic.Handler(r.getNotifications, http.StatusOK))
		notif.GET("/:notification_id", formatDoc("Get Notification", "Get a specific notification"), tonic.Handler(r.getNotification, http.StatusOK))
		notif.DELETE("", formatDoc("Delete Notifications", "Delete matching notifications"), tonic.Handler(r.deleteNotifications, http.StatusOK))
		notif.DELETE("/:notification_id", formatDoc("Delete Notification", "Delete a specific notification"), tonic.Handler(r.deleteNotification, http.StatusOK))
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
		MailId:  res.MailId,
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
		MailId:  res.MailId,
		To:      res.To,
		Subject: res.Subject,
		Body:    res.Body,
	}, nil

}

func (r *Router) postNotification(c *gin.Context, req *AddNotificationReq) (*npb.AddResponse, error) {
	return r.clients.n.Add(req.NodeId, req.Severity,
		req.Type, req.ServiceName, req.Description, req.Details, req.Time)
}

func (r *Router) getNotification(c *gin.Context, req *GetNotificationReq) (*npb.GetResponse, error) {
	return r.clients.n.Get(req.NotificationId)
}

func (r *Router) getNotifications(c *gin.Context, req *GetNotificationsReq) (*npb.ListResponse, error) {
	return r.clients.n.List(req.NodeId, req.ServiceName, req.Type, req.Count, req.Sort)
}

func (r *Router) deleteNotification(c *gin.Context, req *GetNotificationReq) (*npb.DeleteResponse, error) {
	return r.clients.n.Delete(req.NotificationId)
}

func (r *Router) deleteNotifications(c *gin.Context, req *DelNotificationsReq) (*npb.ListResponse, error) {
	return r.clients.n.Purge(req.NodeId, req.ServiceName, req.Type)
}
