/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/notification/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/client"
	epb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
	mailerpb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	npb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
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
	m client.Mailer
	n client.Notify
	e client.EventNotification
	d client.Distributor
}

var upgrader = websocket.Upgrader{
	// Solve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewClientsSet(endpoints *pkg.GrpcEndpoints) *Clients {
	var err error

	c := &Clients{}

	c.m, err = client.NewMailer(endpoints.Mailer, endpoints.Timeout)
	if err != nil {
		log.Fatalf("failed to create mailer client: %v", err)
	}

	c.n, err = client.NewNotify(endpoints.Notify, endpoints.Timeout)
	if err != nil {
		log.Fatalf("failed to create notify client: %v", err)
	}

	c.e, err = client.NewEventNotification(endpoints.EventNotification, endpoints.Timeout)
	if err != nil {
		log.Fatalf("failed to create event-notify client: %v", err)
	}

	c.d, err = client.NewDistributor(endpoints.Distributor, endpoints.Timeout)
	if err != nil {
		log.Fatalf("failed to create distributor client: %v", err)
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
	log.Info("Listening on port ", rt.config.serverConf.Port)

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

		eNotif := auth.Group("/event-notification", "Event Notification", "Event to Notifications")
		eNotif.GET("", formatDoc("Get Notification By filter", "Get a specific notificationby filter"), tonic.Handler(r.getEventNotifications, http.StatusOK))
		eNotif.GET("/:id", formatDoc("Get Notification by Id", "Get a notification"), tonic.Handler(r.getEventNotification, http.StatusOK))
		eNotif.POST("/:id", formatDoc("Update Notifications", "Update matching notification"), tonic.Handler(r.updateEventNotification, http.StatusOK))

		dist := auth.Group("/distributor", "Event distribution", "real time even distribution")
		dist.GET("/live", formatDoc("Real-time Notifications", "Get notification as they are reproted"), tonic.Handler(r.liveEventNotificationHandler, http.StatusOK))

	}
}

func formatDoc(summary string, description string) []fizz.OperationOption {
	return []fizz.OperationOption{func(info *openapi.OperationInfo) {
		info.Summary = summary
		info.Description = description
	}}
}

func (r *Router) sendEmailHandler(c *gin.Context, req *SendEmailReq) (*mailerpb.SendEmailResponse, error) {
	payload := mailerpb.SendEmailRequest{
		To:           req.To,
		TemplateName: req.TemplateName,
		Values:       make(map[string]string),
	}

	// Convert map[string]interface{} to map[string]string
	for key, value := range req.Values {
		if strValue, ok := value.(string); ok {
			payload.Values[key] = strValue
		}
	}

	res, err := r.clients.m.SendEmail(&payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// getEmailByIdHandler handles the get email by ID API endpoint.
func (r *Router) getEmailByIdHandler(c *gin.Context, req *GetEmailByIdReq) (*mailerpb.GetEmailByIdResponse, error) {
	mailerId := req.MailerId
	res, err := r.clients.m.GetEmailById(mailerId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Router) postNotification(c *gin.Context, req *AddNodeNotificationReq) (*npb.AddResponse, error) {
	return r.clients.n.Add(req.NodeId, req.Severity,
		req.Type, req.ServiceName, req.Description, req.Details, req.Status, req.Time)
}

func (r *Router) getNotification(c *gin.Context, req *GetNodeNotificationReq) (*npb.GetResponse, error) {
	return r.clients.n.Get(req.NotificationId)
}

func (r *Router) getNotifications(c *gin.Context, req *GetNodeNotificationsReq) (*npb.ListResponse, error) {
	return r.clients.n.List(req.NodeId, req.ServiceName, req.Type, req.Count, req.Sort)
}

func (r *Router) deleteNotification(c *gin.Context, req *GetNodeNotificationReq) (*npb.DeleteResponse, error) {
	return r.clients.n.Delete(req.NotificationId)
}

func (r *Router) deleteNotifications(c *gin.Context, req *DelNodeNotificationsReq) (*npb.ListResponse, error) {
	return r.clients.n.Purge(req.NodeId, req.ServiceName, req.Type)
}

func (r *Router) getEventNotification(c *gin.Context, req *GetEventNotificationByIdRequest) (*epb.GetResponse, error) {
	return r.clients.e.Get(req.Id)
}

func (r *Router) getEventNotifications(c *gin.Context, req *GetEventNotificationRequest) (*epb.GetAllResponse, error) {
	return r.clients.e.GetAll(req.OrgId, req.NetworkId, req.SubscriberId, req.UserId, req.Role)
}

func (r *Router) updateEventNotification(c *gin.Context, req *UpdateEventNotificationStatusRequest) (*epb.UpdateStatusResponse, error) {
	return r.clients.e.UpdateStatus(req.Id, req.IsRead)
}

func (r *Router) liveEventNotificationHandler(c *gin.Context, req *GetRealTimeEventNotificationRequest) error {

	log.Infof("Requesting real time notifications %+v", req)

	//Upgrade get request to webSocket protocol
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("upgrade: %s", err.Error())
		return err
	}
	defer ws.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := r.clients.d.GetNotificationStream(ctx, req.OrgId, req.NetworkId, req.SubscriberId, req.UserId, req.Scopes)
	if err != nil {
		log.Errorf("error getting notification on stream:Error: %s", err.Error())
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Infof("EOF received from stream")
			return nil
		} else if err == nil {
			log.Infof("received data %+v for request %+v", resp, req)
			w, err := ws.NextWriter(1)
			if err != nil {
				log.Errorf("Error getting writer: %s", err.Error())
				break
			}

			bytes, err := json.Marshal(resp)
			if err != nil {
				log.Errorf("Failed to Marshal notification stream %+v for user %s Error: %v", resp, req.UserId, err)
				break
			}

			_, err = w.Write(bytes)
			if err != nil {
				log.Errorf("Failed to  write notification %+v for user %s to ws response. Error: %s", resp, req.UserId, err)
				break
			}

		} else {
			log.Errorf("Error while fetching the notification. %+v", err)
			break
		}
	}
	log.Infof("Closing real time notifications %+v", req)
	return err

}
