/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type AddOrgRequest struct {
	OrgName     string `path:"org" validate:"required"`
	OrgId       string `json:"org_id" validate:"required"`
	Ip          string `json:"ip"`
	Certificate string `json:"certificate" validate:"required"`
}

type UpdateOrgRequest struct {
	OrgName     string `path:"org" validate:"required"`
	Ip          string `json:"ip"`
	Certificate string `json:"certificate"`
}

type GetOrgsRequest struct{}
type GetOrgRequest struct {
	OrgName string `path:"org" validate:"required"`
}

type AddNodeRequest struct {
	OrgName string `path:"org" validate:"required"`
	NodeId  string `path:"node" validate:"required"`
}

type DeleteNodeRequest struct {
	OrgName string `path:"org" validate:"required"`
	NodeId  string `path:"node" validate:"required"`
}

type GetNodeRequest struct {
	OrgName string `path:"org" validate:"required"`
	NodeId  string `path:"node" validate:"required"`
}

type AddSystemRequest struct {
	OrgName     string `path:"org" validate:"required"`
	SysName     string `path:"system" validate:"required"`
	Ip          string `json:"ip" validate:"required"`
	Certificate string `json:"certificate" validate:"required"`
	Port        int32  `json:"port" validate:"required"`
	URL         string `json:"url"`
}

type UpdateSystemRequest struct {
	OrgName     string `path:"org" validate:"required"`
	SysName     string `path:"system" validate:"required"`
	Ip          string `json:"ip"`
	Certificate string `json:"certificate"`
	Port        int32  `json:"port"`
}

type GetSystemRequest struct {
	OrgName string `path:"org" validate:"required"`
	SysName string `path:"system" validate:"required"`
}

type DeleteSystemRequest struct {
	OrgName string `path:"org" validate:"required"`
	SysName string `path:"system" validate:"required"`
}
