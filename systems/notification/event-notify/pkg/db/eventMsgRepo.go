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

type EventMsgRepo interface {
	Add(event *EventMsg) (uint, error)
	Get(id uint) (*EventMsg, error)
}

type eventMsgRepo struct {
	Db sql.Db
}

func NewEventMsgRepo(db sql.Db) EventMsgRepo {
	return &eventMsgRepo{
		Db: db,
	}
}

func (r *eventMsgRepo) Add(event *EventMsg) (uint, error) {
	d := r.Db.GetGormDb().Create(event)
	return event.ID, d.Error
}

func (r *eventMsgRepo) Get(id uint) (*EventMsg, error) {
	event := &EventMsg{}
	d := r.Db.GetGormDb().Find(event, id)
	if d.Error != nil {
		return nil, d.Error
	}

	return event, nil
}
