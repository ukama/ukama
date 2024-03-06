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

type AccountRepo interface {
	Get() (*Account, error)
}

type accountRepo struct {
	Db sql.Db
}

func NewAccountRepo(db sql.Db) AccountRepo {
	return &accountRepo{
		Db: db,
	}
}

func (s accountRepo) Get() (*Account, error) {
	var account Account

	result := s.Db.GetGormDb().First(&account)
	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}
