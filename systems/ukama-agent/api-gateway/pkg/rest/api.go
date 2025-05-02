/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type ReqData struct {
	Iccid        string `json:"iccid" path:"iccid" validate:"required"`
	Imsi         string `json:"imsi,omitempty"`
	SimPackageId string `json:"sim_package_id,omitempty"`
	PackageId    string `json:"package_id,omitempty"`
	NetworkId    string `json:"network_id,omitempty"`
}

type UsageForPeriodRequest struct {
	Iccid     string `json:"iccid" path:"iccid" validate:"required"`
	StartTime uint64 `query:"start_time" validate:"required"`
	EndTime   uint64 `query:"end_time" validate:"required"`
}

type UsageRequest struct {
	Iccid string `json:"iccid" path:"iccid" validate:"required"`
}

type ActivateReq ReqData

type DeactivateReq ReqData

type UpdatePackageReq ReqData

type ReadSubscriberReq struct {
	Iccid string `path:"iccid" validate:"required"`
}

type QueryUsageRequest struct {
	Iccid   string `form:"iccid" json:"iccid" query:"iccid" binding:"required"`
	NodeId  string `form:"node_id" json:"node_id" query:"node_id" binding:"required"`
	Session uint64 `form:"session" json:"session" query:"session" binding:"required"`
	From    uint64 `form:"from" json:"from" query:"from" binding:"required"`
	To      uint64 `form:"to" json:"to" query:"to" binding:"required"`
	Count   uint32 `form:"count" json:"count" query:"count" binding:"required"`
	Sort    bool   `form:"sort" json:"sort" query:"sort" binding:"required"`
}
