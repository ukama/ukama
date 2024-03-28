/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"database/sql/driver"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
)

type Invoice struct {
	Id         uuid.UUID      `gorm:"primaryKey;type:uuid"`
	InvoiceeId uuid.UUID      `gorm:"uniqueIndex:invoicee_id_period,where:deleted_at is null;not null;type:uuid"`
	NetworkId  uuid.UUID      `gorm:"not null;type:uuid"`
	Period     time.Time      `gorm:"uniqueIndex:invoicee_id_period,where:deleted_at is null;not null"`
	RawInvoice datatypes.JSON `gorm:"not null"`
	IsPaid     bool           `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type InvoiceeType uint8

const (
	InvoiceeTypeUnknown InvoiceeType = iota
	InvoiceeTypeOrg
	InvoiceeTypeSubscriber
)

func (s *InvoiceeType) Scan(value interface{}) error {
	*s = InvoiceeType(uint8(value.(int64)))

	return nil
}

func (s InvoiceeType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s InvoiceeType) String() string {
	t := map[InvoiceeType]string{0: "unknown", 1: "org", 2: "subscriber"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseInvoiceeType(value string) InvoiceeType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return InvoiceeType(i)
	}

	t := map[string]InvoiceeType{"unknown": 0, "org": 1, "subscriber": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return InvoiceeType(0)
	}

	return InvoiceeType(v)
}
