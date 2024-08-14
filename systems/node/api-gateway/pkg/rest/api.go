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

type RestartSiteRequest struct {
	SiteId string `json:"site_id" example:"{{SiteId}}" validate:"required" path:"site_id"`
}

type RestartNodesRequest struct {
	NodeIds []string `json:"node_ids" example:"{{NodeIds}}" validate:"required"`
}

type ToggleInternetSwitchRequest struct {
	SiteId string `json:"site_id" example:"{{SiteId}}" validate:"required" path:"site_id"`
	Status bool   `json:"status" example:"{{Status}}" validate:"required"`
	Port   int32  `json:"port" example:"{{Port}}" validate:"required"`
}
type GetStateByNodeIdRequest struct {
    NodeId string `path:"node_id" binding:"required"`
}

type GetStateHistoryRequest struct {
    NodeId string `path:"node_id" binding:"required"`
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
