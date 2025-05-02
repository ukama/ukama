/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"time"
)

type SimInfo struct {
	Iccid          string `example:"8910300000003540855" json:"iccid" validate:"required"`
	SimType        string `example:"test" json:"sim_type" validate:"required"`
	Msisdn         string `example:"880170124847571" json:"msisdn" validate:"required"`
	SmDpAddress    string `example:"1001.9.0.0.1" json:"smdp_address" validate:"required"`
	ActivationCode string `example:"1010" json:"activation_code" validate:"required"`
	QrCode         string `example:"459081a" json:"qr_code" validate:"required"`
	IsPhysicalSim  bool   `example:"true" json:"is_physical_sim" validate:"required"`
}

type SimPoolStats struct {
	Total     uint64 `json:"count"`
	Available uint64 `json:"available"`
	Consumed  uint64 `json:"consumed"`
	Failed    uint64 `json:"failed"`
}

type SIM struct {
	Id                string    `json:"id" validate:"required"`
	SubscriberId      string    `json:"subscriber_id" validate:"required"`
	Iccid             string    `json:"iccid" validate:"required"`
	Msisdn            string    `json:"msisdn" validate:"required"`
	Package           Package   `json:"package" validate:"required"`
	FirstActivatedOn  time.Time `json:"first_activated_on" validate:"required"`
	LastActivationOn  time.Time `json:"last_activation_on" validate:"required"`
	DeactivationCount uint64    `json:"deactivation_count" validate:"required"`
	IsPhysical        string    `json:"type" validate:"required"`
	SimType           string    `json:"sim_type" validate:"required"`
	ActivationCount   uint64    `json:"activation_count" validate:"required"`
	AllocatedAt       time.Time `json:"allocated_at" validate:"required"`
	Status            string    `json:"status" validate:"required"`
}

type Package struct {
	PackageId       string    `json:"package_id" validate:"required"`
	StartDate       time.Time `json:"start_date" validate:"required"`
	EndDate         time.Time `json:"end_date" validate:"required"`
	DefaultDuration uint64    `json:"default_duration" validate:"required"`
}

type Subscriber struct {
	SubscriberId          string `json:"subscriber_id" validate:"required"`
	Name                  string `json:"name" validate:"required"`
	Email                 string `json:"email" validate:"email,required"`
	Phone                 string `json:"phone" validate:"required"`
	Dob                   string `json:"dob" validate:"required"`
	ProofOfIdentification string `json:"proof_of_identification" validate:"required"`
	IdSerial              string `json:"id_serial" validate:"required"`
	Address               string `json:"address" validate:"required"`
	SimList               []SIM  `json:"sims" validate:"required"`
}

type SimByIccidReq struct {
	Iccid string `example:"8910300000003540855" path:"iccid" validate:"required"`
}

type SimPoolTypeReq struct {
	SimType string `example:"test" form:"sim_type" json:"sim_type" path:"sim_type" binding:"required" validate:"required"`
}

type SimPoolStatReq struct {
	SimType string `example:"test" form:"sim_type" json:"sim_type" path:"sim_type" binding:"required" validate:"required"`
}

type SimPoolRemoveSimReq struct {
	Id []uint64 `example:"[1]" form:"id" json:"id" path:"sim_id" binding:"required" validate:"required"`
}

type SimPoolUploadSimReq struct {
	SimType string `example:"test" form:"sim_type" json:"sim_type" binding:"required" validate:"required"`
	Data    string `example:"SUNDSUQsTVNJU0ROLFNtRHBBZGRyZXNzLEFjdGl2YXRpb25Db2RlLElzUGh5c2ljYWwsUXJDb2RlCjg5MTAzMDAwMDAwMDM1NDA4NTUsODgwMTcwMTI0ODQ3NTcxLDEwMDEuOS4wLjAuMSwxMDEwLFRSVUUsNDU5MDgxYQo4OTEwMzAwMDAwMDAzNTQwODQ1LDg4MDE3MDEyNDg0NzU3MiwxMDAxLjkuMC4wLjIsMTAxMCxUUlVFLDQ1OTA4MWIKODkxMDMwMDAwMDAwMzU0MDgzNSw4ODAxNzAxMjQ4NDc1NzMsMTAwMS45LjAuMC4zLDEwMTAsVFJVRSw0NTkwODFj" type:"array" format:"byte" form:"data" json:"data" binding:"required" validate:"required"`
}

type SimPoolAddSimReq struct {
	SimInfo []SimInfo `form:"sim_info" json:"sim_info" binding:"required"`
}

type SubscriberAddReq struct {
	Name                  string `example:"John" json:"name" validate:"required"`
	Email                 string `example:"john@example.com" json:"email" validate:"required"`
	NetworkId             string `example:"{{NetworkUUID}}" json:"network_id"`
	Gender                string `example:"male" json:"gender"`
	Phone                 string `example:"4151231234" json:"phone"`
	IdSerial              string `example:"123456789" json:"id_serial"`
	Dob                   string `example:"Mon, 02 Jan 2006 15:04:05 MST" json:"dob"`
	ProofOfIdentification string `example:"passport" json:"proof_of_identification"`
	Address               string `example:"Mr John Smith. 132, My Street, Kingston, New York 12401" json:"address"`
}

type SubscriberGetReqByEmail struct {
	Email string `example:"{{SubscriberEmail}}" form:"subscriber_email" json:"subscriber_email" path:"subscriber_email" binding:"required" validate:"required"`
}

type SubscriberGetReq struct {
	SubscriberId string `example:"{{SubscriberUUID}}" form:"subscriber_id" json:"subscriber_id" path:"subscriber_id" binding:"required" validate:"required"`
}

type SubscriberDeleteReq struct {
	SubscriberId string `example:"{{SubscriberUUID}}" form:"subscriber_id" json:"subscriber_id" path:"subscriber_id" binding:"required" validate:"required"`
}

type SubscriberByNetworkReq struct {
	NetworkId string `example:"{{NetworkUUID}}" form:"network_id" json:"network_id" path:"network_id" binding:"required" validate:"required"`
}

type SubscriberUpdateReq struct {
	SubscriberId          string `example:"{{SubscriberUUID}}" path:"subscriber_id" validate:"required"`
	Name                  string `example:"John" json:"name" validate:"required"`
	Phone                 string `example:"4151231234" json:"phone"`
	Address               string `example:"Mr John Smith. 132, My Street, Kingston, New York 12401" json:"address"`
	ProofOfIdentification string `example:"passport" json:"proof_of_identification"`
	IdSerial              string `example:"123456789" json:"id_serial"`
}

type SimListReq struct {
	NetworkId string `example:"{{NetworkUUID}}" form:"network_id" json:"network_id" path:"network_id" binding:"required" validate:"required"`
}

type AllocateSimReq struct {
	SubscriberId  string `example:"{{SubscriberUUID}}" json:"subscriber_id" validate:"required"`
	SimToken      string `example:"pj/9A5Hk8VkkZxOJyu0+9fWs7J6HCOjhmD5jEsIvOfZqmFFFMyStgC3Va4l1b6I5+2ibKOsJjR9KGug=" json:"sim_token"`
	PackageId     string `example:"{{PackageUUID}}" json:"package_id" validate:"required"`
	NetworkId     string `example:"{{NetworkUUID}}" json:"network_id" validate:"required"`
	SimType       string `example:"test" json:"sim_type" validate:"required"`
	TrafficPolicy uint32 `json:"traffic_policy"`
}

type SetActivePackageForSimReq struct {
	SimId     string `example:"{{SimUUID}}" path:"sim_id" validate:"required"`
	PackageId string `example:"{{PackageUUID}}" path:"package_id" validate:"required"`
}
type SimReq struct {
	SimId string `example:"{{SimUUID}}" form:"sim_id" json:"sim_id" path:"sim_id" binding:"required" validate:"required"`
}

type SimByNetworkReq struct {
	NetworkId string `example:"{{NetworkUUID}}" form:"network_id" json:"network_id" path:"network_id" binding:"required" validate:"required"`
}

type ActivateDeactivateSimReq struct {
	SimId  string `example:"{{SimUUID}}" path:"sim_id" validate:"required"`
	Status string `example:"active" json:"status" binding:"required" validate:"required"`
}

type GetSimsBySubReq struct {
	SubscriberId string `example:"{{SubscriberUUID}}" form:"subscriber_id" json:"subscriber_id" path:"subscriber_id" binding:"required" validate:"required"`
}

type PostPkgToSimReq struct {
	SimId     string `example:"{{SimUUID}}" json:"sim_id" validate:"required"`
	PackageId string `example:"{{PackageUUID}}" json:"package_id" validate:"required"`
	StartDate string `example:"" json:"start_date" validate:"required"`
}

type AddPkgToSimReq struct {
	// SimId     string `example:"{{SimUUID}}" json:"sim_id" path:"sim_id" binding:"required" validate:"required"`
	SimId     string `example:"{{SimUUID}}" json:"sim_id" path:"sim_id" validate:"required"`
	PackageId string `example:"{{PackageUUID}}" json:"package_id" validate:"required"`
	StartDate string `example:"" json:"start_date" validate:"required"`
}

type RemovePkgFromSimReq struct {
	SimId     string `example:"{{SimUUID}}" json:"sim_id" path:"sim_id" binding:"required" validate:"required"`
	PackageId string `example:"{{PackageUUID}}" json:"package_id" path:"package_id" binding:"required" validate:"required"`
}

type GetUsagesReq struct {
	SimId   string `form:"sim_id" json:"sim_id" query:"sim_id" binding:"required"`
	SimType string `form:"sim_type" json:"sim_type" query:"sim_type" binding:"required"`
	Type    string `form:"cdr_type" json:"cdr_type" query:"cdr_type" binding:"required"`
	From    string `form:"from" json:"from" query:"from" binding:"required"`
	To      string `form:"to" json:"to" query:"to" binding:"required"`
	Region  string `form:"region" json:"region" query:"region" binding:"required"`
}

type ListSimsReq struct {
	Iccid         string `form:"iccid" json:"iccid" query:"iccid" binding:"required"`
	Imsi          string `form:"Imsi" json:"Imsi" query:"Imsi" binding:"required"`
	SubscriberId  string `form:"subscriber_id" json:"subscriber_id" query:"subscriber_id" binding:"required"`
	NetworkId     string `form:"network_id" json:"network_id" query:"network_id" binding:"required"`
	SimType       string `form:"sim_type" json:"sim_type" query:"sim_type" binding:"required"`
	SimStatus     string `form:"sim_status" json:"sim_status" query:"sim_status" binding:"required"`
	TrafficPolicy uint32 `form:"traffic_policy" json:"traffic_policy" query:"traffic_policy" binding:"required"`
	IsPhysical    bool   `form:"is_physical_sim" json:"is_physical_sim" query:"is_physical_sim" binding:"required"`
	Count         uint32 `form:"count" json:"count" query:"count" binding:"required"`
	Sort          bool   `form:"sort" json:"sort" query:"sort" binding:"required"`
}

type ListPackagesForSimReq struct {
	SimId         string `example:"{{SimUUID}}" form:"sim_id" json:"sim_id" path:"sim_id" binding:"required" validate:"required"`
	DataPlanId    string `form:"data_plan_id" json:"data_plan_id" query:"data_plan_id" binding:"required"`
	FromStartDate string `form:"from_start_date" json:"from_start_date" query:"from_start_date" binding:"required"`
	ToStartDate   string `form:"to_start_date" json:"to_start_date" query:"to_start_date" binding:"required"`
	FromEndDate   string `form:"from_end_date" json:"from_end_date" query:"from_end_date" binding:"required"`
	ToEndDate     string `form:"to_end_date" json:"to_end_date" query:"to_end_date" binding:"required"`
	IsActive      bool   `form:"is_active" json:"is_active" query:"is_active" binding:"required"`
	AsExpired     bool   `form:"as_expired" json:"as_expired" query:"as_expired" binding:"required"`
	Count         uint32 `form:"count" json:"count" query:"count" binding:"required"`
	Sort          bool   `form:"sort" json:"sort" query:"sort" binding:"required"`
}
