/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package api

import (
	"fmt"
	"strings"
)

type RegisterEPC struct {
	ExpectedInstanceID string `json:"expected_instance_id"`
	InstanceID         string `json:"epc_instance_id" validate:"required"`
}

type EPCResponse struct {
	InstanceID string `json:"epc_instance_id"`
	State      string `json:"state"`
}

type GetPDN struct {
	Imsi string `json:"imsi" path:"imsi" validate:"required"`
}

type PDNResponse struct {
	InstanceID string `json:"epc_instance_id"`
	PDNID      string `json:"pdn_id"`
	Imsi       string `json:"imsi"`
	ApnName    string `json:"apn_name"`
	IP         string `json:"ip"`
	State      string `json:"state"`
	DataState  string `json:"data_state"`
	Reason     string `json:"reason"`
}

// NormalizeAPN accepts the textual form used by the local EPC interface.
// A missing or invalid value is never replaced with a default here: EPC
// must send the APN it selected during admission.
func NormalizeAPN(value string) (string, error) {
	if len(value) == 0 || len(value) > 99 {
		return "", fmt.Errorf("invalid_apn")
	}
	for _, label := range strings.Split(value, ".") {
		if len(label) == 0 || len(label) > 63 ||
			label[0] == '-' || label[len(label)-1] == '-' {
			return "", fmt.Errorf("invalid_apn")
		}
		for _, ch := range label {
			if (ch >= 'a' && ch <= 'z') ||
				(ch >= 'A' && ch <= 'Z') ||
				(ch >= '0' && ch <= '9') || ch == '-' {
				continue
			}
			return "", fmt.Errorf("invalid_apn")
		}
	}
	return strings.ToLower(value), nil
}
