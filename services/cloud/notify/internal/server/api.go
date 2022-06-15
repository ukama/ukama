package server

import (
	"github.com/gofrs/uuid"
	"github.com/ukama/ukama/services/cloud/notify/internal/db"
)

type ReqPostNotification struct {
	LookingTo string `query:"looking_to" validate:"eq=post_notification,required"`
	db.Notification
}

type ReqDeleteNotification struct {
	LookingFor string    `query:"looking_to" validate:"eq=delete_notification,required"`
	Id         uuid.UUID `query:"notification_id" validate:"required"`
}

type ReqListNotification struct {
	LookingFor string `query:"looking_for" validate:"eq=list_notification,required"`
}

type ReqGetNotificationTypeForNode struct {
	NodeID     string              `query:"node" validate:"required"`
	LookingFor string              `query:"looking_for" validate:"eq=notification,required"`
	Type       db.NotificationType `json:"notificationType" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqDeleteNotificationForNode struct {
	NodeID     string              `query:"node" validate:"required"`
	LookingFor string              `query:"looking_to" validate:"eq=delete_notification,required"`
	Type       db.NotificationType `query:"notification_type" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqListNotificationForNode struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=list_notification,required"`
}

type ReqGetNotificationTypeForService struct {
	ServiceName string              `query:"service" validate:"required"`
	LookingFor  string              `query:"looking_for" validate:"eq=notification,required"`
	Type        db.NotificationType `json:"notificationType" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqDeleteNotificationForService struct {
	NodeID     string              `query:"service" validate:"required"`
	LookingFor string              `query:"looking_to" validate:"eq=delete_notification,required"`
	Type       db.NotificationType `query:"notification_type" default:"alert" validate:"eq=alert|eq=event"`
}

type ReqListNotificationForService struct {
	ServiceName string `query:"service" validate:"required"`
	LookingFor  string `query:"looking_for" validate:"eq=list_notification,required"`
}
