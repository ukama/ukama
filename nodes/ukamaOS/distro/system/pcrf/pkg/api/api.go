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
	Imsi    string `json:"imsi" validate:"required"`
	ApnName string `json:"apn_name" validate:"required"`
	Ip      string `json:"ip" validate:"required"`
}

type EndSession struct {
	Imsi string `json:"imsi" validate:"required"`
}

type GetSessionByID struct {
	ID uint64 `json:"id" validate:"required"`
}

type GetSessionByImsi struct {
	Imsi string `json:"imsi" validate:"required"`
}

type CDR struct {
	Session    int    `json:"session" validate:"required"`
	Imsi       string `json:"imsi" validate:"required"`
	ApnName    string `json:"apn_name" validate:"required"`
	Ip         string `json:"ip" validate:"required"`
	StartTime  uint64 `json:"start_time" validate:"required"`
	EndTime    uint64 `json:"end_time" validate:"required"`
	TxBytes    uint64 `json:"tx_bytes" validate:"required"`
	RxBytes    uint64 `json:"rx_bytes" validate:"required"`
	TotalBytes uint64 `json:"total_bytes" validate:"required"`
}

type GetCDRBySessionId struct {
	ID uint64 `json:"id" validate:"required"`
}

type GetCDRByImsi struct {
	Imsi string `json:"imsi" path:"imii" validate:"required"`
}

type GetPolicyByImsi struct {
	Imsi string `json:"imsi" validate:"required"`
}

type GetPolicyByID struct {
	ID string `json:"id" validate:"required"`
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
	Id uint64 `json:"id" validate:"required"`
	Ip string `json:"ip" path:"ip" validate:"required"`
}

type AddPolicyByImsi struct {
	Imsi   string `json:"imsi" path:"imsi" validate:"required"`
	Policy Policy `json:"policy" validate:"required"`
}

type Policy struct {
	Uuid      uuid.UUID `json:"uuid" validate:"required"`
	Ulbr      uint64    `json:"ulbr" validate:"required"`
	Dlbr      uint64    `json:"dlbr" validate:"required"`
	Data      uint64    `json:"data" validate:"required"`
	StartTime int64     `json:"start_time" validate:"required"`
	EndTime   int64     `json:"end_time" validate:"required"`
}

type RequestSubscriber struct {
	Imsi string `json:"imsi" validate:"required"`
}

type CreateSubscriber struct {
	Imsi   string `json:"imsi" path:"imsi" validate:"required"`
	Policy Policy `json:"policy" validate:"required"`
}

type GetFlowsForImsi struct {
	Imsi string `json:"imsi" path:"imsi" validate:"required"`
}

type SubscriberResponse struct {
	ID       int       `json:"id" validate:"required"`
	Imsi     string    `json:"imsi" path:"imsi" validate:"required"`
	PolicyID uuid.UUID `json:"policy_id"`
}

type PolicyResponse struct {
	ID        uuid.UUID `json:"id" path:"id"`
	Data      uint64    `json:"data" path:"data"`
	Dlbr      uint64    `json:"dlbr" path:"dlbr"`
	Ulbr      uint64    `json:"ulbr" path:"ulbr"`
	StartTime int64     `json:"start_time" path:"start_time"`
	EndTime   int64     `json:"end_time" path:"end_time"`
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
