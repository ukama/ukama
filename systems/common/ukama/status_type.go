/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukama

import (
	"database/sql/driver"
	"strconv"
	"strings"
)

type StatusType uint8

const (
	StatusTypeUnknown StatusType = iota
	StatusTypePending
	StatusTypeProcessing
	StatusTypeCompleted
	StatusTypeFailed
)

func (s *StatusType) Scan(value interface{}) error {
	*s = StatusType(uint8(value.(int64)))

	return nil
}

func (s StatusType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s StatusType) String() string {
	t := map[StatusType]string{0: "unknown", 1: "pending", 2: "processing", 3: "completed", 4: "failed"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseStatusType(value string) StatusType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return StatusType(i)
	}

	t := map[string]StatusType{"unknown": 0, "pending": 1, "processing": 2, "completed": 3, "failed": 4}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return StatusType(0)
	}

	return StatusType(v)
}
