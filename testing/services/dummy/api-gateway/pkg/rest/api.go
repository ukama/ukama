/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type UpdateReq struct {
	Iccid   string `form:"iccid" json:"iccid" validate:"required"`
	Profile string `form:"profile" json:"profile,omitempty" validate:"eq=normal|eq=min|eq=max" enum:"normal,min,max"`
}

type SimPoolUploadSimReq struct {
	Data string `example:"SUNDSUQsTVNJU0ROLFNtRHBBZGRyZXNzLEFjdGl2YXRpb25Db2RlLElzUGh5c2ljYWwsUXJDb2RlCjg5MTAzMDAwMDAwMDM1NDA4NTUsODgwMTcwMTI0ODQ3NTcxLDEwMDEuOS4wLjAuMSwxMDEwLFRSVUUsNDU5MDgxYQo4OTEwMzAwMDAwMDAzNTQwODQ1LDg4MDE3MDEyNDg0NzU3MiwxMDAxLjkuMC4wLjIsMTAxMCxUUlVFLDQ1OTA4MWIKODkxMDMwMDAwMDAwMzU0MDgzNSw4ODAxNzAxMjQ4NDc1NzMsMTAwMS45LjAuMC4zLDEwMTAsVFJVRSw0NTkwODFj" type:"array" format:"byte" form:"data" json:"data" binding:"required" validate:"required"`
}

type GetSims struct {
}

type GetSimByIccid struct {
	Iccid string `json:"iccid" path:"iccid" validate:"required"`
}
