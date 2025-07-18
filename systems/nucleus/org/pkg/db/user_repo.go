/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
)

type UserRepo interface {
	Add(user *User, nestedFunc func(*User, *gorm.DB) error) error
	Get(uuid uuid.UUID) (*User, error)
	Update(*User) (*User, error)
	Delete(uuid uuid.UUID) error
	GetUserCount() (int64, int64, error)
	AddOrgToUser(user *User, org *Org) error
	RemoveOrgFromUser(user *User, org *Org) error
}

type userRepo struct {
	Db sql.Db
}

func NewUserRepo(db sql.Db) UserRepo {
	return &userRepo{
		Db: db,
	}
}

func (u *userRepo) Add(user *User, nestedFunc func(user *User, tx *gorm.DB) error) error {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		d := tx.Create(user)

		if d.Error != nil {
			return d.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(user, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}

func (u *userRepo) Get(uuid uuid.UUID) (*User, error) {
	var user User

	result := u.Db.GetGormDb().Preload(clause.Associations).Where("uuid = ?", uuid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *userRepo) Update(user *User) (*User, error) {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Model(User{}).Where("uuid = ?", user.Uuid).Updates(user)

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if result.Error != nil {
			return result.Error
		}

		// member := &OrgUser{
		// 	Deactivated: user.Deactivated,
		// }

		// result = tx.Model(OrgUser{}).Where("uuid = ?", user.Uuid).Updates(member)

		// if result.Error != nil {
		// 	return result.Error
		// }

		return nil
	})

	return user, err
}

func (u *userRepo) AddOrgToUser(user *User, org *Org) error {
	err := u.Db.GetGormDb().Model(user).Association("orgs").Append(org)
	return err
}

func (u *userRepo) RemoveOrgFromUser(user *User, org *Org) error {
	err := u.Db.GetGormDb().Model(user).Association("orgs").Delete(org)
	return err
}

func (u *userRepo) Delete(userUUID uuid.UUID) error {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where(&User{Uuid: userUUID}).Delete(&User{})

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (u *userRepo) GetUserCount() (int64, int64, error) {
	var activeUserCount int64
	var deactiveUserCount int64

	result := u.Db.GetGormDb().Model(&User{}).Where("deactivated = ?", false).Count(&activeUserCount)
	if result.Error != nil {
		return 0, 0, result.Error
	}

	result = u.Db.GetGormDb().Model(&User{}).Where("deactivated = ?", true).Count(&deactiveUserCount)
	if result.Error != nil {
		return 0, 0, result.Error
	}

	return activeUserCount, deactiveUserCount, nil
}
