/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

type ContractRepo interface {
	GetContracts(company string, active bool) ([]*Contract, error)
}

type contractRepo struct {
	Db sql.Db
}

func NewContractRepo(db sql.Db) ContractRepo {
	return &contractRepo{
		Db: db,
	}
}

func (c *contractRepo) GetContracts(company string, active bool) ([]*Contract, error) {
	var contracts []*Contract
	if active {
		err := c.Db.GetGormDb().Where("company = ?", company).Order("effective_date desc").First(&contracts).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := c.Db.GetGormDb().Where("company = ?", company).Find(&contracts).Error
		if err != nil {
			return nil, err
		}
	}

	return contracts, nil
}
