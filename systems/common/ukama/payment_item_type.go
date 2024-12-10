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

type ItemType uint8

const (
	ItemTypeUnknown ItemType = iota
	ItemTypePackage
	ItemTypeInvoice
)

func (s *ItemType) Scan(value interface{}) error {
	*s = ItemType(uint8(value.(int64)))

	return nil
}

func (s ItemType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s ItemType) String() string {
	t := map[ItemType]string{0: "unknown", 1: "package", 2: "invoice"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseItemType(value string) ItemType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return ItemType(i)
	}

	t := map[string]ItemType{"unknown": 0, "package": 1, "invoice": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return ItemType(0)
	}

	return ItemType(v)
}
