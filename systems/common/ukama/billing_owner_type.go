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

type OwnerType uint8

const (
	OwnerTypeUnknown OwnerType = iota
	OwnerTypeOrg
	OwnerTypeSubscriber
)

func (s *OwnerType) Scan(value interface{}) error {
	*s = OwnerType(uint8(value.(int64)))

	return nil
}

func (o OwnerType) Value() (driver.Value, error) {
	return int64(o), nil
}

func (o OwnerType) String() string {
	t := map[OwnerType]string{0: "unknown", 1: "org", 2: "subscriber"}

	v, ok := t[o]
	if !ok {
		return t[0]
	}

	return v
}

func ParseOwnerType(value string) OwnerType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return OwnerType(i)
	}

	t := map[string]OwnerType{"unknown": 0, "org": 1, "subscriber": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return OwnerType(0)
	}

	return OwnerType(v)
}
