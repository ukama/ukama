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
