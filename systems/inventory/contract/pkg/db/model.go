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

type Contract struct {
	Id            uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name          string    `gorm:"uniqueIndex:contract_name_idx"`
	VAT           string
	Company       string
	Type          ContractType
	Description   string
	OpexFee       string `gorm:"column:opex_fee"`
	EffectiveDate string `gorm:"column:effective_date"`
}

type ContractType uint8

const (
	UKAMA_PRODUCT ContractType = 0
	BACKHAUL      ContractType = 1
)

func (e *ContractType) Scan(value interface{}) error {
	*e = ContractType(uint8(value.(int64)))
	return nil
}

func (e ContractType) Value() (uint8, error) {
	return uint8(e), nil
}
