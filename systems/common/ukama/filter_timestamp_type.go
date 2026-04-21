/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ukama

import "database/sql/driver"

type FilterTimestampType uint8

const (
	FilterTimestampTypeUnknown FilterTimestampType = iota
	FilterTimestampTypeAll
	FilterTimestampTypeLatest
)

func (s FilterTimestampType) String() string {
	return []string{"unknown", "all", "latest"}[s]
}

func (s *FilterTimestampType) Scan(value interface{}) error {
	*s = FilterTimestampType(uint8(value.(int64)))
	return nil
}

func (s FilterTimestampType) Value() (driver.Value, error) {
	return int64(s), nil
}

func ParseFilterTimestampType(value string) FilterTimestampType {
	switch value {
	case "all":
		return FilterTimestampTypeAll
	case "latest":
		return FilterTimestampTypeLatest
	}
	return FilterTimestampTypeUnknown
}

func ReturnFilterTimestampType(value FilterTimestampType) int64 {
	t := map[FilterTimestampType]int64{FilterTimestampTypeUnknown: 0, FilterTimestampTypeAll: 1,
		FilterTimestampTypeLatest: 2}

	v, ok := t[value]
	if !ok {
		return 0
	}

	return v
}