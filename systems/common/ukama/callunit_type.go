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

type CallUnitType uint8

const (
	CallUnitTypeUnknown CallUnitType = iota
	CallUnitTypeSec
	CallUnitTypeMin
	CallUnitTypeHours
)

func (s *CallUnitType) Scan(value interface{}) error {
	*s = CallUnitType(uint8(value.(int64)))
	return nil
}

func (s CallUnitType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s CallUnitType) String() string {
	t := map[CallUnitType]string{0: "unknown", 1: "seconds", 2: "minutes", 3: "hours"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseCallUnitType(value string) CallUnitType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return CallUnitType(i)
	}

	t := map[string]CallUnitType{"unknown": 0, "seconds": 1, "minutes": 2, "hours": 3}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return CallUnitType(0)
	}

	return CallUnitType(v)
}

func ReturnCallUnits(value CallUnitType) int64 {
	t := map[CallUnitType]int64{CallUnitTypeUnknown: 0, CallUnitTypeSec: 1,
		CallUnitTypeMin: 60, CallUnitTypeHours: 3600}

	v, ok := t[value]
	if !ok {
		return 0
	}

	return v
}
