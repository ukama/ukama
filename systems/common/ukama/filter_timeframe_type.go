/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ukama

import "database/sql/driver"

type FilterTimeframesType uint8

const (
	FilterTimeframesTypeUnknown FilterTimeframesType = iota
	FilterTimeframesTypeAll
	FilterTimeframesTypeLatest
)

func (s FilterTimeframesType) String() string {
	return []string{"unknown", "all", "latest"}[s]
}

func (s *FilterTimeframesType) Scan(value interface{}) error {
	*s = FilterTimeframesType(uint8(value.(int64)))
	return nil
}

func (s FilterTimeframesType) Value() (driver.Value, error) {
	return int64(s), nil
}

func ParseFilterTimeframesType(value string) FilterTimeframesType {
	switch value {
	case "all":
		return FilterTimeframesTypeAll
	case "latest":
		return FilterTimeframesTypeLatest
	}
	return FilterTimeframesTypeUnknown
}

func ReturnFilterTimeframesType(value FilterTimeframesType) int64 {
	t := map[FilterTimeframesType]int64{FilterTimeframesTypeUnknown: 0, FilterTimeframesTypeAll: 1,
		FilterTimeframesTypeLatest: 2}

	v, ok := t[value]
	if !ok {
		return 0
	}

	return v
}