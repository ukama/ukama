/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type UpdateReq struct {
	Iccid    string `form:"iccid" json:"iccid" validate:"required"`
	Profile  string `form:"profile" json:"profile,omitempty" validate:"eq=normal|eq=min|eq=max" enum:"normal,min,max"`
	Scenario string `form:"scenario" json:"scenario,omitempty" validate:"eq=default|eq=backhaul_down|eq=backhaul_downlink_down|eq=node_off|eq=node_on|eq=node_restart|eq=node_rf_off|eq=node_rf_on" enum:"default,backhaul_down,backhaul_downlink_down,node_off,node_on,node_restart,node_rf_off,node_rf_on"`
}

type SimPoolUploadSimReq struct {
	Data string `example:"SUNDSUQsTVNJU0ROLFNtRHBBZGRyZXNzLEFjdGl2YXRpb25Db2RlLElzUGh5c2ljYWwsUXJDb2RlCjg5MTAzMDAwMDAwMDM1NDA4NTUsODgwMTcwMTI0ODQ3NTcxLDEwMDEuOS4wLjAuMSwxMDEwLFRSVUUsNDU5MDgxYQo4OTEwMzAwMDAwMDAzNTQwODQ1LDg4MDE3MDEyNDg0NzU3MiwxMDAxLjkuMC4wLjIsMTAxMCxUUlVFLDQ1OTA4MWIKODkxMDMwMDAwMDAwMzU0MDgzNSw4ODAxNzAxMjQ4NDc1NzMsMTAwMS45LjAuMC4zLDEwMTAsVFJVRSw0NTkwODFj" type:"array" format:"byte" form:"data" json:"data" binding:"required" validate:"required"`
}

type GetSims struct {
}

type GetSimByIccid struct {
	Iccid string `json:"iccid" path:"iccid" validate:"required"`
}

type UpdateSiteMetricsReq struct {
	SiteId      string       `form:"siteId" json:"siteId,omitempty" validate:"required"`
	Profile     string       `form:"profile" json:"profile,omitempty" `
	PortUpdates []PortUpdate `form:"portUpdates" json:"portUpdates,omitempty" `
}

type PortUpdate struct {
	PortNumber int32 `form:"portNumber" json:"portNumber" validate:"required"`
	Status     bool  `form:"status" json:"status" validate:"required"`
}

type SiteConfig struct {
	AvgBackhaulSpeed float64 `form:"avgBackhaulSpeed" json:"avgBackhaulSpeed"`
	AvgLatency       float64 `form:"avgLatency" json:"avgLatency"`
	SolarEfficiency  float64 `form:"solarEfficiency" json:"solarEfficiency"`
}

type StartReq struct {
	SiteId     string      `form:"siteId" json:"siteId,omitempty" validate:"required"`
	Profile    string      `form:"profile" json:"profile,omitempty"`
	SiteConfig *SiteConfig `form:"siteConfig" json:"siteConfig,omitempty"`
}
