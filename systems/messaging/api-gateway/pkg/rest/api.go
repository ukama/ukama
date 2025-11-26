/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type GetNodeIPRequest struct {
	NodeId string `path:"node_id" validate:"required"`
}

type SetNodeIPRequest struct {
	NodeId   string `path:"node_id" validate:"required"`
	NodeIp   string `json:"node_ip" validate:"required"`
	MeshIp   string `json:"mesh_ip" validate:"required"`
	NodePort int32  `json:"node_port" validate:"required"`
	MeshPort int32  `json:"mesh_port" validate:"required"`
	Org      string `json:"org" validate:"required"`
	Network  string `json:"network"`
}

type DeleteNodeIPRequest struct {
	NodeId string `path:"node_id" validate:"required"`
}

type ListNodeIPsRequest struct {
}

type NodeOrgMapListRequest struct {
}
