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

type NotificationRepo interface {
	Add(org *Notification) error
	AddWithEvent(event *EventMsg, notification *Notification, un []*UserNotification) error
	Get(id uuid.UUID) (*Notification, error)
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

func (r *notificationRepo) AddWithEvent(event *EventMsg, notification *Notification, un []*UserNotification) error {
	return r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		notification.EventMsgID = event.ID
		if err := tx.Create(notification).Error; err != nil {
			return err
		}

		if len(un) == 0 {
			return nil
		}

		return tx.Create(un).Error
	})
}

func (r *notificationRepo) Get(id uuid.UUID) (*Notification, error) {
	var notification Notification
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&notification, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &notification, nil
}
