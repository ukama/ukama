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

type SimType uint8

const (
	SimTypeUnknown SimType = iota
	SimTypeTest
	SimTypeOperatorData
	SimTypeUkamaData
)

func (s *SimType) Scan(value interface{}) error {
	*s = SimType(uint8(value.(int64)))
	return nil
}

func (s SimType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s SimType) String() string {
	t := map[SimType]string{0: "unknown", 1: "test", 2: "operator_data", 3: "ukama_data"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseSimType(value string) SimType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return SimType(i)
	}

	t := map[string]SimType{"unknown": 0, "test": 1, "operator_data": 2, "ukama_data": 3}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return SimType(0)
	}

	return SimType(v)
}

type SimStatus uint8

const (
	SimStatusUnknown SimStatus = iota
	SimStatusActive
	SimStatusInactive
	SimStatusTerminated
)

func (s *SimStatus) Scan(value interface{}) error {
	*s = SimStatus(uint8(value.(int64)))
	return nil
}

func (s SimStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s SimStatus) String() string {
	t := map[SimStatus]string{0: "unknown", 1: "active", 2: "inactive", 3: "terminated"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseSimStatus(value string) SimStatus {
	i, err := strconv.Atoi(value)
	if err == nil {
		return SimStatus(i)
	}

	t := map[string]SimStatus{"unknown": 0, "active": 1, "inactive": 2, "terminated": 3}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return SimStatus(0)
	}

	return SimStatus(v)
}
