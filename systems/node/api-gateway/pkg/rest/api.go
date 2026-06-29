/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type PingNodeRequest struct {
	NodeId string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
}

type RestartNodeRequest struct {
	NodeId string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
}

type ToggleSwitchPortRequest struct {
	NodeId string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
	Status bool   `json:"status"`
	Port   int32  `json:"port" validate:"required"`
}

type ToggleStateRequest struct {
	NodeId string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
	State  string `json:"state" validate:"required" example:"on|off" path:"state"`
}
type GetStatesRequest struct {
	NodeId string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
}

type GetStatesHistoryRequest struct {
	NodeId     string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
	PageNumber int32  `json:"page_number" query:"page_number"`
	PageSize   int32  `json:"page_size" query:"page_size"`
	StartTime  string `json:"start_time" query:"start_time"`
	EndTime    string `json:"end_time" query:"end_time"`
}
type EnforceStateTransitionRequest struct {
	NodeId string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
	Event  string `json:"event" validate:"required" example:"{{Event}}" path:"event"`
}
type ApplyConfigRequest struct {
	Commit string `json:"commit" path:"commit" example:"commit" validate:"required"`
}

type GetConfigVersionRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}

type UpdateSoftwareRequest struct {
	Name   string `json:"name" validate:"required" path:"name"`
	Tag    string `json:"tag" validate:"required" path:"tag"`
	NodeId string `json:"node_id" validate:"required" path:"node_id"`
}

type ListAppsRequest struct{}

type ListSoftwareRequest struct {
	NodeId  string `json:"node_id" form:"node_id" query:"node_id" binding:"required"`
	AppName string `json:"app_name" form:"app_name" query:"app_name" binding:"required"`
	Status  string `json:"status" form:"status" query:"status" binding:"required" validate:"eq=unknown|eq=update_available|eq=up_to_date|eq=update_in_progress|eq=update_failed"`
}
type SiteActionRequest struct {
	SiteId      string `json:"site_id" validate:"required" path:"site_id"`
	Reason      string `json:"reason"`
	RequestedBy string `json:"requestedBy"`
	State       string `json:"state" path:"state" validate:"required,oneof=on off"`
}

type SiteStateRequest struct {
	SiteId string `json:"site_id" validate:"required" path:"site_id"`
}

type SitePortMapRequest struct {
	SiteId  string             `json:"site_id" validate:"required" path:"site_id"`
	CNodeId string             `json:"cnode_id"`
	Ports   []SitePortMapEntry `json:"ports" validate:"required"`
}

type SitePortMapEntry struct {
	Port    int32  `json:"port" validate:"required"`
	Role    string `json:"role" validate:"required"`
	NodeId  string `json:"node_id"`
	Class   string `json:"class" validate:"required"`
	Policy  string `json:"policy" validate:"required"`
	CnodeId string `json:"cnode_id"`
}

type PowerCycleNodeRequest struct {
	SiteId      string `json:"site_id" validate:"required" path:"site_id"`
	Role        string `json:"role" validate:"required" path:"role"`
	Reason      string `json:"reason"`
	RequestedBy string `json:"requestedBy"`
}

type ToggleInternetSwitchRequest struct {
	SiteId string `json:"site_id" validate:"required" path:"site_id"`
	Status bool   `json:"status"`
	Port   int32  `json:"port" validate:"required"`
}
