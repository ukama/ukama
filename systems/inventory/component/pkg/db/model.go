/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/uuid"
)

type Component struct {
	Id            uuid.UUID `gorm:"primaryKey;type:uuid"`
	Company       string
	InventoryId   string
	Category      string
	Type          ComponentType
	Description   string
	DatasheetURL  string
	ImagesURL     string
	PartNumber    string
	Manufacturer  string
	Managed       string
	Warranty      uint32
	Specification string
}

type ComponentType uint8

const (
	Power    ComponentType = 0
	Backhaul ComponentType = 1
	Switch   ComponentType = 2
	Access   ComponentType = 3
)

func (e *ComponentType) Scan(value interface{}) error {
	*e = ComponentType(uint8(value.(int64)))

	return nil
}

func (e ComponentType) Value() (uint8, error) {
	return uint8(e), nil
}
