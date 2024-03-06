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
	Get() (*Contract, error)
}

type contractRepo struct {
	Db sql.Db
}

func NewContractRepo(db sql.Db) ContractRepo {
	return &contractRepo{
		Db: db,
	}
}

func (s contractRepo) Get() (*Contract, error) {
	var contract Contract

	result := s.Db.GetGormDb().First(&contract)
	if result.Error != nil {
		return nil, result.Error
	}

	return &contract, nil
}
