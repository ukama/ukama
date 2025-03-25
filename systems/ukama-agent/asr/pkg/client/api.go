/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"time"
)

type SimCardInfo struct {
	Imsi           string `json:"imsi,omitempty"`
	Iccid          string `json:"iccid,omitempty"`
	Op             []byte `json:"op,omitempty"`
	Amf            []byte `json:"amf"`
	Key            []byte `json:"key,omitempty"`
	AlgoType       uint32 `json:"algo_type,omitempty"`
	UeDlAmbrBps    uint32 `json:"ue_dl_ambr_bps,omitempty"`
	UeUlAmbrBps    uint32 `json:"ue_ul_ambr_bps,omitempty"`
	Sqn            uint64 `json:"sqn,omitempty"`
	CsgIdPrsent    bool   `json:"c_sg_id_prsent,omitempty"`
	CsgId          uint32 `json:"csg_id,omitempty"`
	DefaultApnName string `json:"default_apn_name,omitempty"`
}

type Sim struct {
	SimCardInfo *SimCardInfo `json:"sim"`
}

type NetworkInfo struct {
	NetworkId     string    `json:"id"`
	Name          string    `json:"name"`
	OrgId         string    `json:"org_id"`
	IsDeactivated bool      `json:"is_deactivated"`
	CreatedAt     time.Time `json:"created_at"`
}
