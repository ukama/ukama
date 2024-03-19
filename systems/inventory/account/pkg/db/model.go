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

	"github.com/ukama/ukama/systems/common/uuid"
)

type Account struct {
	Id            uuid.UUID `gorm:"primaryKey;type:uuid"`
	Item          string
	Company       string
	Quantity      uint32 `gorm:"type:integer"`
	Category      string
	PricePerUnit  float64 `gorm:"type:float"`
	TotalPrice    float64 `gorm:"type:float"`
	Description   string
	Specification string
	PaymentType   PaymentType
}

type Category uint8

const (
	SASS               Category = 0
	NODE               Category = 1
	SIM                Category = 2
	BACKHAUL_FIBER     Category = 3
	BACKHAUL_SATELLITE Category = 4
)

func (s *Category) Scan(value interface{}) error {
	*s = Category(uint8(value.(int64)))
	return nil
}

func (s Category) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s Category) String() string {
	t := map[Category]string{0: "SASS", 1: "NODE", 2: "SIM", 3: "BACKHAUL_FIBER", 4: "BACKHAUL_SATELLITE"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseCategoryType(value string) Category {
	i, err := strconv.Atoi(value)
	if err == nil {
		return Category(i)
	}

	t := map[string]Category{"SASS": 0, "NODE": 1, "SIM": 2, "BACKHAUL_FIBER": 3, "BACKHAUL_SATELLITE": 4}

	v, ok := t[value]
	if !ok {
		return Category(0)
	}

	return Category(v)
}

type PaymentType uint8

const (
	ONE_TIME PaymentType = 0
	MONTHLY  PaymentType = 1
	YEARLY   PaymentType = 2
)

func (s *PaymentType) Scan(value interface{}) error {
	*s = PaymentType(uint8(value.(int64)))
	return nil
}

func (s PaymentType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s PaymentType) String() string {
	t := map[PaymentType]string{0: "ONE_TIME", 1: "MONTHLY", 2: "YEARLY"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseType(value string) PaymentType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return PaymentType(i)
	}

	t := map[string]PaymentType{"ONE_TIME": 0, "MONTHLY": 1, "YEARLY": 2}

	v, ok := t[value]
	if !ok {
		return PaymentType(0)
	}

	return PaymentType(v)
}
