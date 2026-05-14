/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import "encoding/json"

type ListHealthReportsRequest struct {
	ReportId   string `form:"reportId" json:"reportId" query:"reportId"`
	NodeId     string `form:"nodeId" json:"nodeId" query:"nodeId"`
	ReportedAt string `form:"reportedAt" json:"reportedAt" query:"reportedAt"`
	Timeframe  string `form:"timeframe" default:"all" json:"timeframe" query:"timeframe" binding:"required" validate:"eq=all|eq=latest"`
}

type ListAppsRequest struct {
	NodeId string `form:"nodeId" json:"nodeId" query:"nodeId"`
	AppName string `form:"appName" json:"appName" query:"appName"`
}

type ListInterfacesRequest struct {
	NodeId string `form:"nodeId" json:"nodeId" query:"nodeId"`
	InterfaceName string `form:"interfaceName" json:"interfaceName" query:"interfaceName"`
}

type AddNotificationReq struct {
	NodeId      string          `form:"node_id" json:"node_id"`
	Severity    string          `form:"severity" json:"severity,omitempty" validate:"eq=fatal|eq=critical|eq=high|eq=medium|eq=low|eq=clean|eq=log|eq=warning|eq=debug|eq=trace"`
	Type        string          `form:"type" json:"type,omitempty" validate:"eq=alert|eq=event"`
	ServiceName string          `form:"service_name" json:"service_name,omitempty"`
	Status      uint32          `form:"status" json:"status,omitempty"`
	Time        uint32          `form:"time" json:"time,omitempty"`
	Details     json.RawMessage `form:"details" json:"details,omitempty"`
}

type GetNotificationReq struct {
	NotificationId string `json:"notification_id" path:"notification_id" validate:"required"`
}

type GetNotificationsReq struct {
	NodeId      string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	ServiceName string `form:"service_name" json:"service_name" query:"service_name" binding:"required"`
	Type        string `form:"type" json:"type" query:"type" binding:"required"`
	Count       uint32 `form:"count" json:"count" query:"count" binding:"required"`
	Sort        bool   `form:"sort" json:"sort" query:"sort" binding:"required"`
}

type DelNotificationsReq struct {
	NodeId      string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	ServiceName string `form:"service_name" json:"service_name" query:"service_name" binding:"required"`
	Type        string `form:"type" json:"type" query:"type" binding:"required"`
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

type StoreHealthReportRequest struct {
	NodeId string          `path:"node_id" validate:"required" example:"{{NodeId}}"`
	raw    json.RawMessage `json:"-"`
}

// StoreHealthReportOpenAPIInput is only for OpenAPI (fizz.InputModel). Runtime binding
// uses StoreHealthReportRequest, which accepts the same JSON object as the request body.
type StoreHealthReportOpenAPIInput struct {
	NodeId         string                   `path:"node_id" example:"{{NodeId}}"`
	SchemaVersion string                   `json:"schemaVersion" example:"1"`
	NodeID        string                   `json:"nodeId,omitempty"`
	NodeType      string                   `json:"nodeType" example:"HomeNode"`
	ReportedAt    string                   `json:"reportedAt" example:"2023-12-12T00:00:00Z"`
	Capabilities  []string                 `json:"capabilities,omitempty"`
	System        map[string]interface{}   `json:"system,omitempty"`
	Interfaces    map[string]interface{}   `json:"interfaces,omitempty"`
	Apps          []map[string]interface{} `json:"apps,omitempty"`
	Events        []map[string]interface{} `json:"events,omitempty"`
}

// UnmarshalJSON captures the entire POST body as the health payload (HealthPayload JSON).
func (s *StoreHealthReportRequest) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		s.raw = []byte("{}")
		return nil
	}
	// copy so later mutations do not alias gin's buffer
	s.raw = append(json.RawMessage(nil), data...)
	return nil
}

// HealthPayloadBytes returns bytes forwarded to the health service (defaults to {}).
func (s *StoreHealthReportRequest) HealthPayloadBytes() []byte {
	if len(s.raw) == 0 {
		return []byte("{}")
	}
	return s.raw
}
