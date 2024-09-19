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
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type NotificationResult struct {
	Notifications
	NodeStateID   *uuid.UUID `gorm:"column:node_state_id"`
	NodeID        *string    `gorm:"column:node_id"`
	Latitude      *float64   `gorm:"column:latitude"`
	Longitude     *float64   `gorm:"column:longitude"`
	CurrentState  *string    `gorm:"column:current_state"`
	NodeStateName *string    `gorm:"column:name"`
}

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
	var results []NotificationResult
	query := `
		SELECT
			UN.is_read,
			N.title,
			N.description,
			N.scope,
			N.type,
			N.id,
			N.created_at,
			N.updated_at,
			NS.id as node_state_id,
			NS.node_id,
			NS.latitude,
			NS.longitude,
			NS.current_state,
			NS.name
		FROM
			user_notifications AS UN
		INNER JOIN
			notifications AS N ON UN.notification_id = N.id
		LEFT JOIN
			node_states AS NS ON UN.node_state_id = NS.id
		WHERE
			UN.user_id = ?;
	`

	result := r.Db.GetGormDb().Raw(query, id).Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	notifications := make([]*Notifications, len(results))

	for i, res := range results {
		notifications[i] = &res.Notifications
		if res.NodeStateID != nil {
			notifications[i].NodeState = &NodeState{
				Id:           *res.NodeStateID,
				NodeId:       *res.NodeID,
				Latitude:     *res.Latitude,
				Longitude:    *res.Longitude,
				CurrentState: *res.CurrentState,
				Name:         *res.NodeStateName,
			}
		}
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
