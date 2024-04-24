/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
)

type NotificationRepo interface {
	Add(org *Notification) error
	Get(id uuid.UUID) (*Notification, error)
	Update(id uuid.UUID, isRead bool) error
	GetAll(orgId string, networkId string, subscriberId string, userId string) ([]Notification, error)
}

type notificationRepo struct {
	Db sql.Db
}

func NewNotificationRepo(db sql.Db) NotificationRepo {
	return &notificationRepo{
		Db: db,
	}
}

func (r *notificationRepo) Add(notification *Notification) (err error) {
	d := r.Db.GetGormDb().Create(notification)
	return d.Error
}

func (r *notificationRepo) Get(id uuid.UUID) (*Notification, error) {
	var notification Notification
	result := r.Db.GetGormDb().First(&notification, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &notification, nil
}

func (r *notificationRepo) GetAll(orgId string, networkId string, subscriberId string, userId string) ([]Notification, error) {
	var notifications []Notification

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	if orgId != "" {
		tx = tx.Where("org_id = ?", orgId)
	}

	if networkId != "" {
		tx = tx.Where("network_id = ?", networkId)
	}

	if subscriberId != "" {
		tx = tx.Where("subscriber_id = ?", subscriberId)
	}

	if userId != "" {
		tx = tx.Where("user_id = ?", userId)
	}

	result := tx.Find(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return notifications, nil
}

func (r *notificationRepo) Update(id uuid.UUID, isRead bool) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Notification{}).Where("id = ?", id).Update("is_read", isRead).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
