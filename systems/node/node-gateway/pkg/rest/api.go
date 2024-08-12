/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type GetRunningAppsRequest struct {
	NodeId string `example:"{{NodeId}}" validate:"required" path:"node_id" `
}

type AddNotificationReq struct {
	NodeId      string `json:"node_id"`
	Severity    string `json:"severity,omitempty" type:"string"`
	Type        string `json:"notification_type,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	Status      uint32 `json:"status,omitempty"`
	Time        uint32 `json:"time,omitempty"`
	Description string `json:"description,omitempty"`
	Details     string `json:"details,omitempty"`
}


type GetNotificationReq struct {
	NotificationId string `json:"notification_id" path:"notification_id" validate:"required"`
}

type GetNotificationsReq struct {
	NodeId      string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	ServiceName string `form:"service_name" json:"service_name" query:"service_name" binding:"required"`
	Type        string `form:"notification_type" json:"notification_type" query:"notification_type" binding:"required"`
	Count       uint32 `form:"count" json:"count" query:"count" binding:"required"`
	Sort        bool   `form:"sort" json:"sort" query:"sort" binding:"required"`
}

type DelNotificationsReq struct {
	NodeId      string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	ServiceName string `form:"service_name" json:"service_name" query:"service_name" binding:"required"`
	Type        string `form:"notification_type" json:"notification_type" query:"notification_type" binding:"required"`
}

type Logs struct {
	AppName string `json:"app_name" validate:"required"`
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type AddLogsRequest struct {
	NodeId string `example:"{{NodeId}}" json:"node_id" path:"node_id"`
	Logs   []Logs `json:"logs" validate:"required"`
}

type StoreRunningAppsInfoRequest struct {
	NodeId    string   `validate:"required" example:"{{NodeId}}" path:"node_id" `
	Timestamp string   `json:"timestamp" validate:"required" example:"{{time}}"`
	System    []System `json:"system" validate:"required" example:"{{system}}"`
	Capps     []Capps  `json:"capps" validate:"required" example:"{{capps}}"`
}

type System struct {
	Name  string `json:"name" validate:"required" example:"{{name}}"`
	Value string `json:"value" validate:"required" example:"{{value}}"`
}

type Capps struct {
	Space     string      `json:"space" validate:"required" example:"{{space}}"`
	Name      string      `json:"name" validate:"required" example:"{{name}}"`
	Tag       string      `json:"tag" validate:"required" example:"{{tag}}"`
	Status    string      `json:"status" validate:"required" example:"{{status}}"`
	Resources []Resources `json:"resources" validate:"required" example:"{{resources}}"`
}

type Resources struct {
	Name  string `json:"name" validate:"required" example:"{{name}}"`
	Value string `json:"value" validate:"required" example:"{{value}}"`
}
