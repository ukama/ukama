/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

 type SendEmailRes struct {
	 Message  string `json:"message"`
	 MailerId string `json:"mailer_id"`
 }
 
 type SendEmailReq struct {
	 To           []string `json:"to" validate:"required"`
	 TemplateName string   `json:"template_name" validate:"required"`
	 Values       map[string]interface{}
 }
 
 type GetEmailByIdReq struct {
	 MailerId string `json:"mailer_id" validate:"required" path:"mailer_id" binding:"required"`
 }
 
 type AddNodeNotificationReq struct {
	 NodeId      string `json:"node_id"`
	 Severity    string `json:"severity,omitempty" type:"string"`
	 Type        string `json:"notification_type,omitempty" validate:"eq=alert|eq=event"`
	 ServiceName string `json:"service_name,omitempty"`
	 Status      uint32 `json:"status,omitempty"`
	 Time        uint32 `json:"time,omitempty"`
	 Description string `json:"description,omitempty"`
	 Details     string `json:"details,omitempty"`
 }
 
 type GetNodeNotificationReq struct {
	 NotificationId string `json:"notification_id" path:"notification_id" validate:"required"`
 }
 
 type GetNodeNotificationsReq struct {
	 NodeId      string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	 ServiceName string `form:"service_name" json:"service_name" query:"service_name" binding:"required"`
	 Type        string `form:"notification_type" json:"notification_type" query:"notification_type" binding:"required"`
	 Count       uint32 `form:"count" json:"count" query:"count" binding:"required"`
	 Sort        bool   `form:"sort" json:"sort" query:"sort" binding:"required"`
 }
 
 type DelNodeNotificationsReq struct {
	 NodeId      string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	 ServiceName string `form:"service_name" json:"service_name" query:"service_name" binding:"required"`
	 Type        string `form:"notification_type" json:"notification_type" query:"notification_type" binding:"required"`
 }
 
 type GetEventNotificationByIdRequest struct {
	 Id string `path:"id" binding:"required"`
 }
 
 type GetEventNotificationRequest struct {
	 OrgId        string `json:"org_id" form:"org_id" query:"org_id"`
	 NetworkId    string `json:"network_id" form:"network_id" query:"network_id"`
	 SubscriberId string `json:"subscriber_id" form:"subscriber_id" query:"subscriber_id"`
	 UserId       string `json:"user_id" form:"user_id" query:"user_id"`
 }
 
 type UpdateEventNotificationStatusRequest struct {
	 Id     string `path:"id" binding:"required"`
	 IsRead bool   `form:"read" json:"read" query:"is_read" binding:"required"`
 }
 
 type GetRealTimeEventNotificationRequest struct {
	 OrgId        string   `json:"org_id" form:"org_id" query:"org_id"`
	 NetworkId    string   `json:"network_id" form:"network_id" query:"network_id"`
	 SubscriberId string   `json:"subscriber_id" form:"subscriber_id" query:"subscriber_id"`
	 UserId       string   `json:"user_id" form:"user_id" query:"user_id"`
	 Scopes       []string `json:"scopes" form:"scopes" query:"scope"`
 }