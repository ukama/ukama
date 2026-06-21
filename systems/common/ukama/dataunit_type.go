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

type DataUnitType uint8

const (
	DataUnitTypeUnknown DataUnitType = iota
	DataUnitTypeB
	DataUnitTypeKB
	DataUnitTypeMB
	DataUnitTypeGB
)

func (s *DataUnitType) Scan(value interface{}) error {
	*s = DataUnitType(uint8(value.(int64)))
	return nil
}

func (s DataUnitType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s DataUnitType) String() string {
	t := map[DataUnitType]string{0: "unknown", 1: "bytes", 2: "kilobytes", 3: "megabytes", 4: "gigabytes"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseDataUnitType(value string) DataUnitType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return DataUnitType(i)
	}

	s := map[string]DataUnitType{"unknown": 0, "bytes": 1, "kilobytes": 2, "megabytes": 3, "gigabytes": 4}
	t := map[string]DataUnitType{"unknown": 0, "b": 1, "kb": 2, "mb": 3, "gb": 4}

	v, ok := s[strings.ToLower(value)]
	if ok {
		return DataUnitType(v)
	}

	v, ok = t[strings.ToLower(value)]
	if ok {
		return DataUnitType(v)
	}

	return DataUnitType(0)
}

func ReturnDataUnits(value DataUnitType) int64 {
	t := map[DataUnitType]int64{DataUnitTypeUnknown: 0, DataUnitTypeB: 1, DataUnitTypeKB: 1,
		DataUnitTypeMB: 1, DataUnitTypeGB: 1024}

	v, ok := t[value]
	if !ok {
		return 0
	}
	return v
}

func ReturnDataUnitsInBytes(value DataUnitType) uint64 {
	t := map[DataUnitType]int64{DataUnitTypeUnknown: 0, DataUnitTypeB: 1, DataUnitTypeKB: 1024,
		DataUnitTypeMB: 1024 * 1024, DataUnitTypeGB: 1024 * 1024 * 1024}

	v, ok := t[value]
	if !ok {
		return 0
	}

	return uint64(v)
}
