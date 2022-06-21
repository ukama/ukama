package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	uuid "github.com/satori/go.uuid"
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

	if repo != nil {
		r.n = notify.NewNotify(repo)
	}

	r.init()

	return r
}

func (r *Router) init() {
	notif := r.fizz.Group("notification", "Notification", "Notifications")
	notif.POST("", nil, tonic.Handler(r.postNewNotification, http.StatusCreated))
	notif.DELETE("", nil, tonic.Handler(r.deleteNotification, http.StatusOK))
	notif.GET("/list", nil, tonic.Handler(r.listNotification, http.StatusOK))

	node := notif.Group("node", "Node", "Node")
	node.DELETE("", nil, tonic.Handler(r.deleteNotificationForNode, http.StatusOK))
	node.GET("", nil, tonic.Handler(r.getNotificationForNode, http.StatusOK))
	node.GET("/list", nil, tonic.Handler(r.listNotificationForNode, http.StatusOK))

	service := notif.Group("service", "Service", "Service")
	service.DELETE("", nil, tonic.Handler(r.deleteNotificationForService, http.StatusOK))
	service.GET("", nil, tonic.Handler(r.getNotificationForService, http.StatusOK))
	service.GET("/list", nil, tonic.Handler(r.listNotificationForService, http.StatusOK))

}

/* Handle new notification */
func (r *Router) postNewNotification(c *gin.Context, req *ReqPostNotification) error {
	logrus.Debugf("Handling new notification: %+v.", req)

	/* validate nodeid */
	_, err := ukama.ValidateNodeId(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid node: " + err.Error(),
		}
	}

	nf := NewDbNotification(req)

	err = r.n.NewNotificationHandler(nf)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to register new notification:" + err.Error(),
		}
	}

	return nil
}

/* delete notification */
func (r *Router) deleteNotification(c *gin.Context, req *ReqDeleteNotification) error {
	logrus.Debugf("Handling delete notification: %+v.", req)

	id, err := uuid.FromString(req.Id)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid UUID supplied as notification ID. Error: " + err.Error(),
		}
	}

	err = r.n.DeleteNotification(id)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to delete notifications: Error" + err.Error(),
		}
	}

	return nil
}

/* List notification */
func (r *Router) listNotification(c *gin.Context, req *ReqListNotification) (*RespNotificationList, error) {
	logrus.Debugf("Handling list notifications: %+v.", req)
	var resp *RespNotificationList
	list, err := r.n.ListNotification()
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to get notification list. Error:" + err.Error(),
		}
	}

	if list != nil {
		resp = getNotificationList(list)
	}
	logrus.Debugf("Response with notifications: %+v.", resp)
	return resp, nil
}

func (r *Router) getNotificationForNode(c *gin.Context, req *ReqGetNotificationTypeForNode) (*RespNotificationList, error) {
	logrus.Debugf("Handling list notifications for NodeId : %+v.", req)

	var resp *RespNotificationList
	list, err := r.n.GetSpecificNotification(nil, &req.NodeID, string(req.Type))
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to get notification list for node " + req.NodeID + "Error:" + err.Error(),
		}
	}

	if list != nil {
		resp = getNotificationList(list)
	}
	logrus.Debugf("Response with notifications for node: %+v.", resp)

	return resp, nil
}

func (r *Router) deleteNotificationForNode(c *gin.Context, req *ReqDeleteNotificationForNode) error {
	logrus.Debugf("Handling delete notification for node: %+v.", req)

	err := r.n.DeleteSpecificNotification(nil, &req.NodeID, string(req.Type))
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to delete notification for node" + req.NodeID + "Error" + err.Error(),
		}
	}

	return nil
}

func (r *Router) listNotificationForNode(c *gin.Context, req *ReqListNotificationForNode) (*RespNotificationList, error) {
	logrus.Debugf("Handling list notifications for node: %+v.", req)
	var resp *RespNotificationList
	list, err := r.n.ListSpecificNotification(nil, &req.NodeID, req.Count)
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to get notification list for node" + req.NodeID + ": Error" + err.Error(),
		}
	}

	if list != nil {
		resp = getNotificationList(list)
	}
	logrus.Debugf("Response with notifications for node: %+v.", resp)

	return resp, nil
}

func (r *Router) getNotificationForService(c *gin.Context, req *ReqGetNotificationTypeForService) (*RespNotificationList, error) {
	logrus.Debugf("Handling list notifications for Service : %+v.", req)

	var resp *RespNotificationList
	list, err := r.n.GetSpecificNotification(&req.ServiceName, nil, string(req.Type))
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to get notification list for Service " + req.ServiceName + "Error:" + err.Error(),
		}
	}

	if list != nil {
		resp = getNotificationList(list)
	}
	logrus.Debugf("Response with notifications for service: %+v.", resp)

	return resp, nil
}

func (r *Router) deleteNotificationForService(c *gin.Context, req *ReqDeleteNotificationForService) error {
	logrus.Debugf("Handling delete notification for service: %+v.", req)

	err := r.n.DeleteSpecificNotification(&req.ServiceName, nil, string(req.Type))
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to delete notification for service" + req.ServiceName + "Error" + err.Error(),
		}
	}

	return nil
}

func (r *Router) listNotificationForService(c *gin.Context, req *ReqListNotificationForService) (*RespNotificationList, error) {
	logrus.Debugf("Handling list notifications for service: %+v.", req)
	var resp *RespNotificationList
	list, err := r.n.ListSpecificNotification(&req.ServiceName, nil, req.Count)
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to get notification list for service" + req.ServiceName + ": Error" + err.Error(),
		}
	}

	if list != nil {
		resp = getNotificationList(list)
	}
	logrus.Debugf("Response with notifications for service: %+v.", resp)

	return resp, nil
}

func NewDbNotification(r *ReqPostNotification) *db.Notification {
	n := &db.Notification{
		NotificationID: uuid.NewV4(),
		NodeID:         r.NodeID,
		NodeType:       r.NodeType,
		Severity:       db.SeverityType(r.Severity),
		Type:           db.NotificationType(r.Type),
		ServiceName:    r.ServiceName,
		Time:           r.Time,
		Description:    r.Description,
		Details:        r.Details,
	}
	return n
}

func getNotificationList(list *[]db.Notification) *RespNotificationList {
	var resp = &RespNotificationList{}

	size := len(*list)
	logrus.Debugf("Notification Count is %d", size)

	resp.Notifications = make([]Notification, size)
	for idx, nt := range *list {
		resp.Notifications[idx] = Notification{
			NotificationID: nt.NotificationID,
			NodeID:         nt.NodeID,
			NodeType:       nt.NodeType,
			Severity:       nt.Severity.String(),
			Type:           nt.Type.String(),
			ServiceName:    nt.ServiceName,
			Time:           nt.Time,
			Description:    nt.Description,
			Details:        nt.Details,
		}
	}
	return resp
}
