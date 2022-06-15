package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/cloud/notify/cmd/version"
	"github.com/ukama/ukama/services/cloud/notify/internal"
	"github.com/ukama/ukama/services/cloud/notify/internal/db"
	"github.com/ukama/ukama/services/cloud/notify/internal/notify"
	"github.com/ukama/ukama/services/common/rest"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/services/common/ukama"

	"github.com/wI2L/fizz"
)

type Router struct {
	fizz *fizz.Fizz
	port int
	repo db.NotificationRepo
	n    *notify.Notify
	s    *sr.ServiceRouter
}

func (r *Router) Run(close chan error) {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		close <- err
	}
	close <- nil
}

func NewRouter(config *internal.Config, svcR *sr.ServiceRouter, repo db.NotificationRepo) *Router {

	f := rest.NewFizzRouter(&config.Server, internal.ServiceName, version.Version, internal.IsDebugMode)

	r := &Router{fizz: f,
		port: config.Server.Port,
		repo: repo,
	}

	if svcR != nil {
		r.s = svcR
	}

	r.init()

	return r
}

func (r *Router) init() {
	notif := r.fizz.Group("notification", "Notification", "Notifications")
	notif.POST("", nil, tonic.Handler(r.PostNewNotification, http.StatusOK))
	notif.DELETE("", nil, tonic.Handler(r.DeleteNotification, http.StatusOK))
	notif.GET("/list", nil, tonic.Handler(r.ListNotification, http.StatusOK))

	node := notif.Group("node", "Node", "Node")
	node.DELETE("", nil, tonic.Handler(r.DeleteNotificationForNode, http.StatusOK))
	node.GET("", nil, tonic.Handler(r.GetNotificationForNode, http.StatusOK))
	node.GET("/list", nil, tonic.Handler(r.ListNotificationForNode, http.StatusOK))

	service := notif.Group("service", "Service", "Service")
	service.DELETE("", nil, tonic.Handler(r.DeleteNotificationForService, http.StatusOK))
	service.GET("", nil, tonic.Handler(r.GetNotificationForService, http.StatusOK))
	service.GET("/list", nil, tonic.Handler(r.ListNotificationForService, http.StatusOK))

}

func (r *Router) DeleteNotification(c *gin.Context, req *ReqDeleteNotification) error {
	return nil
}

func (r *Router) ListNotification(c *gin.Context, req *ReqListNotification) error {
	return nil
}

func (r *Router) PostNewNotification(c *gin.Context, req *ReqPostNotification) error {
	logrus.Debugf("Handling new notification opertaion on %+v.", req)

	/* validate nodeid */
	_, err := ukama.ValidateNodeId(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid node:" + err.Error(),
		}
	}
	return nil
}

func (r *Router) DeleteNotificationForNode(c *gin.Context, req *ReqDeleteNotificationForNode) error {
	return nil
}

func (r *Router) ListNotificationForNode(c *gin.Context, req *ReqListNotificationForNode) (*db.Notification, error) {
	return nil, nil
}

func (r *Router) GetNotificationForNode(c *gin.Context, req *ReqGetNotificationTypeForNode) (*db.Notification, error) {
	return nil, nil
}

func (r *Router) DeleteNotificationForService(c *gin.Context, req *ReqDeleteNotificationForService) error {
	return nil
}

func (r *Router) ListNotificationForService(c *gin.Context, req *ReqListNotificationForService) (*db.Notification, error) {
	return nil, nil
}

func (r *Router) GetNotificationForService(c *gin.Context, req *ReqGetNotificationTypeForService) (*db.Notification, error) {
	return nil, nil
}
