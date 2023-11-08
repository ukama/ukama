/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

type ReqGetNode struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=fact_node_info,required"`
}

type RespGetNode struct {
	NodeID string `json:"node"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type ReqBuildNode struct {
	LookingTo string `query:"looking_to" validate:"eq=create_node,required"`
	Type      string `query:"type" validate:"eq=HNODE|eq=TNODE|eq=ANODE|eq=hnode|eq=tnode|eq=anode,required"`
	Count     int    `query:"count" default:"1" type:"integer"`
}

type RespBuildNode struct {
	NodeIDList []string `json:"NodeID"`
}

type ReqDeleteNode struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_to" validate:"eq=fact_delete,required"`
}

type ReqGetNodeList struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=fact_node_list,required"`
}

type RespGetNodeList struct {
	NodeList []RespGetNode `json:"nodes"`
}
