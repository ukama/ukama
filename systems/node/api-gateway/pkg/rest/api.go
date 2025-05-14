/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type PingNodeRequest struct {
	NodeId    string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
	RequestId string `json:"request_id" validate:"required"`
	Message   string `json:"message"`
	TimeStamp uint64 `json:"time_stamp"`
}

type RestartNodeRequest struct {
	NodeId string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
}

type GetStatesRequest struct {
	NodeId string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
}

type RestartSiteRequest struct {
	SiteId    string `json:"site_id"  example:"site-1" validate:"required" path:"site_name"`
	NetworkId string `json:"network_id" example:"{{NetworkId}}" validate:"required" path:"network_id"`
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

type RestartNodesRequest struct {
	NetworkId string   `json:"network_id" example:"{{NetworkId}}" validate:"required" path:"network_id"`
	NodeIds   []string `json:"node_ids" example:"{{NodeIds}}" validate:"required"`
}
type ToggleInternetSwitchRequest struct {
	SiteId string `json:"site_id" example:"{{SiteId}}" validate:"required" path:"site_id"`
	Status bool   `json:"status" example:"{{Status}}"`
	Port   int32  `json:"port" example:"{{Port}}" validate:"required"`
}
type ToggleRfRequest struct {
	NodeId string `json:"node_id" example:"{{NodeId}}" validate:"required" path:"node_id"`
	Status bool   `json:"status" example:"{{Status}}"`
}

type ApplyConfigRequest struct {
	Commit string `json:"commit" path:"commit" example:"commit" validate:"required"`
}

type GetConfigVersionRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}

type UpdateSoftwareRequest struct {
	Space  string `json:"space" validate:"required" path:"space"`
	Name   string `json:"name" validate:"required" path:"name"`
	Tag    string `json:"tag" validate:"required" path:"tag"`
	NodeId string `json:"node_id" validate:"required" path:"node_id"`
}
