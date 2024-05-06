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
	initTrigger(db)
	return &userNotificationRepo{
		Db: db,
	}
}

func initTrigger(db sql.Db) {
	db.GetGormDb().Exec("CREATE FUNCTION public.user_notifications_trigger() RETURNS TRIGGER AS $$ DECLARE notification_data text; BEGIN notification_data := NEW.id::text || ',' || NEW.notification_id::text; PERFORM pg_notify('user_notifications_channel', notification_data); RETURN NEW; END; $$ LANGUAGE plpgsql;")
	db.GetGormDb().Exec("CREATE TRIGGER notify_trigger AFTER INSERT OR UPDATE ON user_notifications FOR EACH ROW EXECUTE FUNCTION public.user_notifications_trigger();")
}

func (r *userNotificationRepo) Add(un []*UserNotification) error {
	d := r.Db.GetGormDb().Create(un)
	return d.Error
}
