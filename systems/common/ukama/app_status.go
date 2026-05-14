/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ukama

import (
	"database/sql/driver"
	"strconv"
)

type AppStatus uint8

const (
	AppStatusUnknown AppStatus = iota
	AppStatusPending
	AppStatusActive
	AppStatusInactive
	AppStatusRunning
	AppStatusStopped
)

func (s *AppStatus) Scan(value interface{}) error {
	*s = AppStatus(uint8(value.(int64)))

	return nil
}

func (s AppStatus) Value() (driver.Value, error) {
	return int64(s), nil
}


func (s AppStatus) String() string {
	t := map[AppStatus]string{0: "unknown", 1: "pending", 2: "active", 3: "inactive", 4: "running", 5: "stopped"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseAppStatus(value string) AppStatus {
	i, err := strconv.Atoi(value)
	if err == nil {
		return AppStatus(i)
	}

	t := map[string]AppStatus{"unknown": 0, "pending": 1, "active": 2, "inactive": 3, "running": 4, "stopped": 5}

	v, ok := t[value]
	if !ok {
		return AppStatus(0)
	}
	return AppStatus(v)
}
