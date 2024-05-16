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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo interface {
	Add(user *Users) error
	GetUsers(orgId string, networkId string, subscriberId string, userId string) ([]*Users, error)
}

type userRepo struct {
	Db sql.Db
}

func NewUserRepo(db sql.Db) UserRepo {
	return &userRepo{
		Db: db,
	}
}

func (r *userRepo) Add(user *Users) (err error) {
	d := r.Db.GetGormDb().Create(user)
	return d.Error
}

func (r *userRepo) GetUsers(orgId string, networkId string, subscriberId string, userId string) ([]*Users, error) {
	var users []*Users

	const emptyUUID = "00000000-0000-0000-0000-000000000000"

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	if orgId != "" && orgId != emptyUUID {
		tx = tx.Where("org_id = ?", orgId)
	}

	if networkId != "" && networkId != emptyUUID {
		tx = tx.Where("network_id = ?", networkId)
	}

	if subscriberId != "" && subscriberId != emptyUUID {
		tx = tx.Where("subscriber_id = ?", subscriberId)
	}

	if userId != "" && userId != emptyUUID {
		tx = tx.Where("user_id = ?", userId)
	}

	result := tx.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return users, nil
}
