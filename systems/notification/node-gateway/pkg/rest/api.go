package rest

import (
	"github.com/ukama/ukama/systems/common/rest"
)

type AddNotificationReq struct {
	rest.BaseRequest
	NodeId      string `json:"node_id"`
	Severity    string `json:"severity,omitempty" type:"string"`
	Type        string `json:"notification_type,omitempty" validate:"eq=alert|eq=event"`
	ServiceName string `json:"service_name,omitempty"`
	Status      uint32 `json:"status,omitempty"`
	Time        uint32 `json:"time,omitempty"`
	Description string `json:"description,omitempty"`
	Details     string `json:"details,omitempty"`
}

type GetNotificationReq struct {
	rest.BaseRequest
	NotificationId string `json:"notification_id" path:"notification_id" validate:"required"`
}

type GetNotificationsReq struct {
	rest.BaseRequest
	NodeId      string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	ServiceName string `form:"service_name" json:"service_name" query:"service_name" binding:"required"`
	Type        string `form:"notification_type" json:"notification_type" query:"notification_type" binding:"required"`
	Count       uint32 `form:"count" json:"count" query:"count" binding:"required"`
	Sort        bool   `form:"sort" json:"sort" query:"sort" binding:"required"`
}

type DelNotificationsReq struct {
	rest.BaseRequest
	NodeId      string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	ServiceName string `form:"service_name" json:"service_name" query:"service_name" binding:"required"`
	Type        string `form:"notification_type" json:"notification_type" query:"notification_type" binding:"required"`
}
