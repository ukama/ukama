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
	t := map[DataUnitType]string{0: "unknown", 1: "Bytes", 2: "KiloBytes", 3: "MegaBytes", 4: "GigaBytes"}

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

	t := map[string]DataUnitType{"unknown": 0, "Bytes": 1, "KiloBytes": 2, "MegaBytes": 3, "GigaBytes": 4}

	v, ok := t[value]
	if !ok {
		return DataUnitType(0)
	}

	return DataUnitType(v)
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
