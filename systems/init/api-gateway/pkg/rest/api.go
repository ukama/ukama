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
	ApiGwIp     string `json:"api_gw_ip" validate:"required"`
	Certificate string `json:"certificate" validate:"required"`
	ApiGwPort   int32  `json:"api_gw_port" validate:"required"`
	ApiGwUrl    string `json:"api_gw_url"`
	NodeGwIp    string `json:"node_gw_ip" default:"0.0.0.0"`
	NodeGwPort  int32  `json:"node_gw_port" default:"8080"`
}

type UpdateSystemRequest struct {
	OrgName     string `path:"org" validate:"required"`
	SysName     string `path:"system" validate:"required"`
	Ip          string `json:"ip"`
	Certificate string `json:"certificate"`
	Port        int32  `json:"port"`
	NodeGwIp    string `json:"node_gw_ip"`
	NodeGwPort  int32  `json:"node_gw_port"`
	NodeGwURL   string `json:"node_gw_url"`
}

type GetSystemRequest struct {
	OrgName string `path:"org" validate:"required"`
	SysName string `path:"system" validate:"required"`
}

type DeleteSystemRequest struct {
	OrgName string `path:"org" validate:"required"`
	SysName string `path:"system" validate:"required"`
}
