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

type MessageUnitType uint8

const (
	MessageUnitTypeUnknown MessageUnitType = iota
	MessageUnitTypeInt
)

func (s *MessageUnitType) Scan(value interface{}) error {
	*s = MessageUnitType(uint8(value.(int64)))
	return nil
}

func (s MessageUnitType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s MessageUnitType) String() string {
	t := map[MessageUnitType]string{0: "unknown", 1: "int"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseMessageType(value string) MessageUnitType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return MessageUnitType(i)
	}

	t := map[string]MessageUnitType{"unknown": 0, "int": 1}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return MessageUnitType(0)
	}

	return MessageUnitType(v)
}

func ReturnMessageUnits(value MessageUnitType) int64 {
	t := map[MessageUnitType]int64{MessageUnitTypeUnknown: 0, MessageUnitTypeInt: 1}

	v, ok := t[value]
	if !ok {
		return 0
	}
	return v
}
