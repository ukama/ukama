/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

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

type CDR struct {
	Session    string `json:"session" validate:"required"`
	Imsi       string `json:"imsi" validate:"required"`
	ApnName    string `json:"apn_name" validate:"required"`
	Ip         string `json:"ip" validate:"required"`
	StartTime  uint64 `json:"start_time" validate:"required"`
	EndTime    uint64 `json:"end_time" validate:"required"`
	TxBytes    uint64 `json:"tx_bytes" validate:"required"`
	RXBytes    uint64 `json:"rx_bytes" validate:"required"`
	TotalBytes uint64 `json:"total_bytes" validate:"required"`
}

type GetCDRById struct {
	ID uint64 `json:"id" validate:"required"`
}

type GetCDRByImsi struct {
	Imsi string `json:"imsi" path:"imii" validate:"required"`
}

type PolicyByImsi struct {
	Imsi string `json:"imsi" validate:"required"`
}

type AddPolicyByImsi struct {
	Imsi   string `json:"imsi" validate:"required"`
	Policy Policy `json:"policy" validate:"required"`
}

type Policy struct {
	Ulbr uint64 `json:"ulbr" validate:"required"`
	Dlbr uint64 `json:"dlbr" validate:"required"`
	Data uint64 `json:"data" validate:"required"`
}

type Reroute struct {
	Ip string `json:"ip" validate:"required"`
}
