/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const EmptyUUID = "00000000-0000-0000-0000-000000000000"

type UserRepo interface {
	Add(user *Users) error
	GetUsers(orgId string, networkId string, subscriberId string, userId string) ([]*Users, error)
	GetAllUsers(orgId string) ([]*Users, error)
	GetUser(userId string) (*Users, error)
	GetSubscriber(subscriberId string) (*Users, error)
	GetUserWithRoles(orgId string, roles []roles.RoleType) ([]*Users, error)
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

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	if orgId != "" && orgId != EmptyUUID {
		tx = tx.Where("org_id = ?", orgId)
	}

	if networkId != "" && networkId != EmptyUUID {
		tx = tx.Where("network_id = ?", networkId)
	}

	if subscriberId != "" && subscriberId != EmptyUUID {
		tx = tx.Where("subscriber_id = ?", subscriberId)
	}

	if userId != "" && userId != EmptyUUID {
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

func (r *userRepo) GetAllUsers(orgId string) ([]*Users, error) {
	var users []*Users

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	if orgId != "" && orgId != EmptyUUID {
		tx = tx.Where("org_id = ?", orgId)
	} else {
		return nil, fmt.Errorf("invalid uuid %s", orgId)
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

func (r *userRepo) GetUser(userId string) (*Users, error) {
	var user *Users

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	if userId != "" && userId != EmptyUUID {
		tx = tx.Where("user_id = ?", userId)
	} else {
		return nil, fmt.Errorf("invalid uuid %s", userId)
	}

	result := tx.Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (r *userRepo) GetSubscriber(subscriberId string) (*Users, error) {
	var user *Users

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	if subscriberId != "" && subscriberId != EmptyUUID {
		tx = tx.Where("subscriber_id = ?", subscriberId)
	} else {
		return nil, fmt.Errorf("invalid uuid %s", subscriberId)
	}

	result := tx.Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (r *userRepo) GetUserWithRoles(orgId string, roleTypes []roles.RoleType) ([]*Users, error) {
	var users []*Users

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	if orgId != "" && orgId != EmptyUUID {
		tx = tx.Where("org_id = ?", orgId)
	}

	tx = tx.Where("role IN ?", roleTypes)

	result := tx.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return users, nil
}
