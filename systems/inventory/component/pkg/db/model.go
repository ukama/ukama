/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"strconv"

	"github.com/ukama/ukama/systems/common/uuid"
)

type Component struct {
	Id            uuid.UUID `gorm:"primaryKey;type:uuid"`
	Company       string
	Inventory     string
	UserId        string
	Category      ComponentCategory
	Type          string
	Description   string
	DatasheetURL  string
	ImagesURL     string
	PartNumber    string
	Manufacturer  string
	Managed       string
	Warranty      uint32
	Specification string
}

type ComponentCategory uint8

const (
	ACCESS   ComponentCategory = 0
	BACKHAUL ComponentCategory = 1
	POWER    ComponentCategory = 2
	SWITCH   ComponentCategory = 3
)

func (c *ComponentCategory) Scan(value interface{}) error {
	*c = ComponentCategory(uint8(value.(int64)))

	return nil
}

func (c ComponentCategory) Value() (uint8, error) {
	return uint8(c), nil
}

func (c ComponentCategory) String() string {
	t := map[ComponentCategory]string{0: "access", 1: "backhaul", 2: "power", 3: "switch"}

	v, ok := t[c]
	if !ok {
		return t[0]
	}

	return v
}

func ParseType(value string) ComponentCategory {
	i, err := strconv.Atoi(value)
	if err == nil {
		return ComponentCategory(i)
	}

	t := map[string]ComponentCategory{"access": 0, "backhaul": 1, "power": 2, "switch": 3}

	v, ok := t[value]
	if !ok {
		return ComponentCategory(0)
	}

	return ComponentCategory(v)
}
