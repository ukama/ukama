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

type UserNotificationRepo interface {
	Add(un []*UserNotification) error
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
