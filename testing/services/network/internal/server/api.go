/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import "github.com/ukama/ukama/testing/services/network/internal/db"

type ReqActionOnNode struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_to" validate:"eq=vnode_power_on|eq=vnode_power_off,required"`
	Org       string `query:"org" validate:"required"`
}

type ReqPowerOnNode struct {
	ReqActionOnNode
}

type ReqPowerOffNode struct {
	ReqActionOnNode
}

type ReqGetNode struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=vnode_info,required"`
}

type RespGetNode struct {
	Node    db.VNode `json:"node"`
	Runtime string   `json:"runtime"`
}

type ReqGetNodeList struct {
	LookingFor string `query:"looking_for" validate:"eq=vnode_list,required"`
}

type RespGetNodeList struct {
	NodeList []db.VNode `json:"nodes"`
}

type ReqDeleteNode struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_to" validate:"eq=vnode_delete,required"`
}
