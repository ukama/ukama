package server

import (
	uuid "github.com/satori/go.uuid"
	"github.com/ukama/ukama/services/cloud/notify/internal/db"
	"gorm.io/datatypes"
)

type Notification struct {
	NotificationID uuid.UUID      `json:"notificationID"`
	NodeID         string         `json:"nodeID"`
	NodeType       string         `json:"nodeType"`
	Severity       string         `json:"severity,omitempty" type:"string"`
	Type           string         `json:"notificationType,omitempty" validate:"eq=alert|eq=event"`
	ServiceName    string         `json:"serviceName,omitempty"`
	Time           uint32         `json:"time,omitempty"`
	Description    string         `json:"description,omitempty"`
	Details        datatypes.JSON `json:"details,omitempty"`
}

type ReqPostNotification struct {
	Notification
}

type ReqDeleteNotification struct {
	Id string `query:"notification_id" validate:"required"`
}

type RespNotificationList struct {
	Notifications []Notification `json:"notifications"`
}

type ReqGetNotificationTypeForNode struct {
	NodeID string              `query:"node" validate:"required"`
	Type   db.NotificationType `query:"type" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqDeleteNotificationForNode struct {
	NodeID string              `query:"node" validate:"required"`
	Type   db.NotificationType `query:"type" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqListNotificationForNode struct {
	NodeID string `query:"node" validate:"required"`
	Count  int    `query:"count" validate:"gte=1,lte=10" default:"5"`
}

type ReqGetNotificationTypeForService struct {
	ServiceName string              `query:"service" validate:"required"`
	Type        db.NotificationType `query:"type" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqDeleteNotificationForService struct {
	ServiceName string              `query:"service" validate:"required"`
	Type        db.NotificationType `query:"type" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqListNotificationForService struct {
	ServiceName string `query:"service" validate:"required"`
	Count       int    `query:"count" validate:"gte=1,lte=10" default:"5"`
}
