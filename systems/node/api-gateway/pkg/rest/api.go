/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type RestartNodeRequest struct {
	NodeId string `json:"node_id" validate:"required" example:"{{NodeId}}" path:"node_id"`
}

type RestartSiteRequest struct {
	SiteName  string `json:"site_name"  example:"site1" validate:"required" path:"site_name"`
	NetworkId string `json:"network_id" example:"{{NetworkId}}" validate:"required" path:"network_id"`
}

type RestartNodesRequest struct {
	NetworkId string   `json:"network_id" example:"{{NetworkId}}" validate:"required" path:"network_id"`
	NodeIds   []string `json:"node_ids" example:"{{NodeIds}}" validate:"required"`
}

type ApplyConfigRequest struct {
	Commit string `json:"commit" path:"commit" example:"commit" validate:"required"`
}

type GetConfigVersionRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}

type UpdateSoftwareRequest struct {
	Space string `json:"space" validate:"required" path:"space"`
	Name  string `json:"name" validate:"required" path:"name"`
	Tag   string `json:"tag" validate:"required" path:"tag"`
	NodeId string `json:"node_id" validate:"required" path:"node_id"`
}
