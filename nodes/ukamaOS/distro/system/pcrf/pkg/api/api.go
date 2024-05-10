/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package api

import "github.com/ukama/ukama/systems/common/uuid"

type CreateSession struct {
	Imsi    []uint8 `json:"imsi" validate:"required"`
	ImsiStr string  `json:"imsi_str,omitempty"`
	ApnName string  `json:"apn_name"`
	Ip      uint32  `json:"pdn_address" validate:"required"`
	IpStr   string  `json:"ip_str,omitempty"`
}

type EndSession struct {
	Imsi    []uint8 `json:"imsi" validate:"required"`
	ImsiStr string  `json:"imsi_str,omitempty"`
}

type GetSessionByID struct {
	ID uint64 `json:"id" path:"id" validate:"required"`
}

type GetSessionByImsi struct {
	Imsi string `json:"imsi" validate:"required"`
}

type CDR struct {
	Session       int    `json:"session" validate:"required"`
	NodeId        string `json:"node_id" validate:"required"`
	Imsi          string `json:"imsi" validate:"required"`
	Policy        string `json:"policy" validate:"required"`
	ApnName       string `json:"apn_name" validate:"required"`
	Ip            string `json:"ip" validate:"required"`
	StartTime     uint64 `json:"start_time" validate:"required"`
	EndTime       uint64 `json:"end_time" validate:"required"`
	LastUpdatedAt uint64 `json:"last_updated_at" validate:"required"`
	TxBytes       uint64 `json:"tx_bytes" validate:"required"`
	RxBytes       uint64 `json:"rx_bytes" validate:"required"`
	TotalBytes    uint64 `json:"total_bytes" validate:"required"`
}

type GetCDRBySessionId struct {
	ID uint64 `json:"id" path:"id" validate:"required"`
}

type GetCDRByImsi struct {
	Imsi string `json:"imsi" path:"imsi" validate:"required"`
}

type GetPolicyByImsi struct {
	Imsi string `json:"imsi" path:"imsi" validate:"required"`
}

type GetPolicyByID struct {
	ID string `json:"id" path:"id" validate:"required"`
}

type GetReRouteByImsi struct {
	Imsi string `json:"imsi" path:"imsi" validate:"required"`
}

type GetReRoutes struct {
}

type ReRoute struct {
	Ip string `json:"ip" validate:"required"`
}

type UpdateRerouteById struct {
	Id uint64 `json:"id" path:"id" validate:"required"`
	Ip string `json:"ip" validate:"required"`
}

type Policy struct {
	Uuid      uuid.UUID `json:"uuid" validate:"required"`
	Ulbr      uint64    `json:"ulbr" validate:"required"`
	Dlbr      uint64    `json:"dlbr" validate:"required"`
	Data      uint64    `json:"total_data" validate:"required"`
	Consumed  uint64    `json:"consumed_data"`
	Burst     uint64    `json:"burst" validate:"required"`
	StartTime int64     `json:"start_time" validate:"required"`
	EndTime   int64     `json:"end_time" validate:"required"`
}

type Spr struct {
	Imsi    string       `json:"imsi" path:"imsi" validate:"required"`
	Policy  Policy       `json:"policy" validate:"required"`
	ReRoute string       `json:"reroute" validate:"required"`
	Usage   UsageDetails `json:"usage" validate:"required"`
}

type UsageDetails struct {
	Data uint64 `json:"data" path:"data" validate:"required"`
	Time uint64 `json:"time" path:"time" validate:"required"`
}

type AddPolicy Policy

type RequestSubscriber struct {
	Imsi string `json:"imsi" path:"imsi" validate:"required"`
}

type CreateSubscriber struct {
	Imsi    string `json:"imsi" path:"imsi" validate:"required"`
	Policy  Policy `json:"policy" validate:"required"`
	ReRoute string `json:"reroute" validate:"required"`
}

type UpdateSubscriber CreateSubscriber

type GetFlowsForImsi struct {
	Imsi string `json:"imsi" path:"imsi" validate:"required"`
}

type SubscriberResponse struct {
	ID       int       `json:"id" validate:"required"`
	Imsi     string    `json:"imsi" path:"imsi" validate:"required"`
	PolicyID uuid.UUID `json:"policy_id"`
	ReRoute  string    `json:"reroute"`
}

type PolicyResponse struct {
	ID        uuid.UUID `json:"uuid" path:"id"`
	Burst     uint64    `json:"burst" path:"burst"`
	Data      uint64    `json:"total_data" path:"data"`
	Consumed  uint64    `json:"consumed_data" path:"consumed"`
	Dlbr      uint64    `json:"dlbr" path:"dlbr"`
	Ulbr      uint64    `json:"ulbr" path:"ulbr"`
	StartTime int64     `json:"start_time" path:"start_time"`
	EndTime   int64     `json:"end_time" path:"end_time"`
	CreatedAt int64     `json:"created_at" path:"created_at"`
	UpdatedAt int64     `json:"updated_at" path:"updated_at"`
}

type UsageRequest struct {
	Imsi string `json:"imsi" path:"imsi" validate:"required"`
}

type UsageResponse struct {
	Imsi string `json:"imsi" path:"imsi" validate:"required"`
	Data uint64 `json:"data" path:"data"`
	Time int64  `json:"time" path:"time"`
}
type SessionResponse struct {
	ID         int
	Imsi       string
	NodeId     string
	PolicyID   string
	ApnName    string
	UeIpaddr   string
	StartTime  uint64
	EndTime    uint64
	TxBytes    uint64
	RxBytes    uint64
	TotalBytes uint64
	TxMeterId  uint32
	RxMeterId  uint32
	State      string
	Sync       string
	UpdatedAt  uint64
}

type MeterResponse struct {
	ID        int
	Rate      uint64
	BurstSize uint64
	Type      string
}

type FlowResponse struct {
	ID        int
	Cookie    uint64
	Table     uint64
	Priority  uint64
	UeIpaddr  string
	ReRouting string
	MeterID   uint32
}

type ReRouteResponse struct {
	ID int
	Ip string
}
