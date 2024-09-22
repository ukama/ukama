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

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)
 
 type UserNotificationRepo interface {
	 Add(un []*UserNotification) error
	 Update(id uuid.UUID, isRead bool) error
	 GetNotificationsByUserID(id string) ([]*Notifications, error)
 }
 
 type userNotificationRepo struct {
	 Db sql.Db
 }
 
 func NewUserNotificationRepo(db sql.Db) UserNotificationRepo {
	 return &userNotificationRepo{
		 Db: db,
	 }
 }
 
 func (r *userNotificationRepo) Add(un []*UserNotification) error {
	 d := r.Db.GetGormDb().Create(un)
	 return d.Error
 }
 
 func (r *userNotificationRepo) GetNotificationsByUserID(id string) ([]*Notifications, error) {
	 var notifications []*Notifications
	 q := fmt.Sprintf("SELECT user_notifications.is_read,notifications.title,notifications.description,notifications.scope,notifications.type,notifications.id,notifications.created_at,notifications.updated_at FROM user_notifications INNER JOIN notifications ON user_notifications.notification_id = notifications.id WHERE user_notifications.user_id = '%s';", id)
	 d := r.Db.GetGormDb().Exec(q).Find(&notifications)
	 if d.Error != nil {
		 return nil, d.Error
	 }
	 return notifications, nil
 }
 
 func (r *userNotificationRepo) Update(id uuid.UUID, isRead bool) error {
	 err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		 if err := tx.Model(&UserNotification{}).Where("notification_id = ?", id).Update("is_read", isRead).Error; err != nil {
			 return err
		 }
 
		 return nil
	 })
 
	 return err
 }