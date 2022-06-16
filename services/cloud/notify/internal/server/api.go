package server

import (
	uuid "github.com/satori/go.uuid"
	"github.com/ukama/ukama/services/cloud/notify/internal/db"
	"gorm.io/datatypes"
)

type Notification struct {
	NotificationID uuid.UUID           `json:"notificationID"`
	NodeID         string              `json:"nodeID"`
	NodeType       string              `json:"nodeType"`
	Severity       db.SeverityType     `json:"severity,omitempty" type:"string"`
	Type           db.NotificationType `json:"notificationType,omitempty" type:"string"`
	ServiceName    string              `json:"serviceName,omitempty"`
	Time           uint32              `json:"time,omitempty"`
	Description    string              `json:"description,omitempty"`
	Details        datatypes.JSON      `json:"details,omitempty"`
}

type ReqPostNotification struct {
	LookingTo string `query:"looking_to" validate:"eq=post_notification,required"`
	Notification
}

type ReqDeleteNotification struct {
	LookingFor string    `query:"looking_to" validate:"eq=delete_notification,required"`
	Id         uuid.UUID `query:"notification_id" validate:"required"`
}

type RespNotificationList struct {
	Notifications []Notification `json:"notifications"`
}

type ReqListNotification struct {
	LookingFor string `query:"looking_for" validate:"eq=list_notification,required"`
}

type ReqGetNotificationTypeForNode struct {
	NodeID     string              `query:"node" validate:"required"`
	LookingFor string              `query:"looking_for" validate:"eq=notification,required"`
	Type       db.NotificationType `query:"type" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqDeleteNotificationForNode struct {
	NodeID     string              `query:"node" validate:"required"`
	LookingFor string              `query:"looking_to" validate:"eq=delete_notification,required"`
	Type       db.NotificationType `query:"type" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqListNotificationForNode struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=list_notification,required"`
}

type ReqGetNotificationTypeForService struct {
	ServiceName string              `query:"service" validate:"required"`
	LookingFor  string              `query:"looking_for" validate:"eq=notification,required"`
	Type        db.NotificationType `query:"type" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqDeleteNotificationForService struct {
	ServiceName string              `query:"service" validate:"required"`
	LookingFor  string              `query:"looking_to" validate:"eq=delete_notification,required"`
	Type        db.NotificationType `query:"type" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqListNotificationForService struct {
	ServiceName string `query:"service" validate:"required"`
	LookingFor  string `query:"looking_for" validate:"eq=list_notification,required"`
}
