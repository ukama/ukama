/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ukama

import (
	"strconv"
	"strings"
)

type SoftwareStatusType uint8

const (
	SoftwareStatusTypeUnknown SoftwareStatusType = 0
	SoftwareStatusTypeUpdateAvailable SoftwareStatusType = 1
	SoftwareStatusTypeUpToDate SoftwareStatusType = 2
	SoftwareStatusTypeUpdateInProgress SoftwareStatusType = 3
	SoftwareStatusTypeUpdateFailed SoftwareStatusType = 4
)

func (s SoftwareStatusType) String() string {
	t := map[SoftwareStatusType]string{0: "unknown", 1: "update_available", 2: "up_to_date", 3: "update_in_progress", 4: "update_failed"}
	v, ok := t[s]
	if !ok {
		return t[0]
	}
	return v
}

func ParseSoftwareStatusType(value string) SoftwareStatusType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return SoftwareStatusType(i)
	}
	t := map[string]SoftwareStatusType{"unknown": 0, "update_available": 1, "up_to_date": 2, "update_in_progress": 3, "update_failed": 4}
	v, ok := t[strings.ToLower(value)]
	if !ok {
		return SoftwareStatusType(0)
	}
	return SoftwareStatusType(v)
}

func (e SoftwareStatusType) Value() (uint8, error) {
	return uint8(e), nil
}

func (e *SoftwareStatusType) Scan(value interface{}) error {
	*e = SoftwareStatusType(uint8(value.(int64)))
	return nil
}