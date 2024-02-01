/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package types

import (
	"database/sql/driver"
	"strconv"
)

type SyncStatus uint8

const (
	SyncStatusUnknown SyncStatus = iota
	SyncStatusPending
	SyncStatusProcessing
	SyncStatusCompleted
	SyncStatusFailed
)

func (s *SyncStatus) Scan(value interface{}) error {
	*s = SyncStatus(uint8(value.(int64)))

	return nil
}

func (s SyncStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s SyncStatus) String() string {
	t := map[SyncStatus]string{0: "unknown", 1: "pending", 2: "processing", 3: "completed", 4: "failed"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseStatus(value string) SyncStatus {
	i, err := strconv.Atoi(value)
	if err == nil {
		return SyncStatus(i)
	}

	t := map[string]SyncStatus{"unknown": 0, "pending": 1, "processing": 2, "completed": 3, "failed": 4}

	v, ok := t[value]
	if !ok {
		return SyncStatus(0)
	}

	return SyncStatus(v)
}
