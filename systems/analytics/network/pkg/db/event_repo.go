/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"time"

	"github.com/ukama/ukama/systems/common/sql"
)

type EventRepo interface {
	Recent(networkId, siteId, nodeId string, from, to time.Time, page, pageSize uint32) ([]EventLog, int64, error)
}

type eventRepo struct {
	Db sql.Db
}

func NewEventRepo(db sql.Db) EventRepo {
	return &eventRepo{
		Db: db,
	}
}

// Recent returns event logs in the window, newest first. Events are filtered
// by resource by matching the routing key or payload against the resource id
// (events carry routing keys like event.cloud.local.{org}.{system}.{service}.{action}
// and resource ids in the payload).
func (r *eventRepo) Recent(networkId, siteId, nodeId string, from, to time.Time, page, pageSize uint32) ([]EventLog, int64, error) {
	var events []EventLog
	var count int64

	q := r.Db.GetGormDb().Model(&EventLog{})

	if !from.IsZero() {
		q = q.Where("occurred_at >= ?", from)
	}
	if !to.IsZero() {
		q = q.Where("occurred_at < ?", to)
	}

	if nodeId != "" {
		q = q.Where("payload::text ILIKE ?", "%"+nodeId+"%")
	} else if siteId != "" {
		q = q.Where("payload::text ILIKE ?", "%"+siteId+"%")
	} else if networkId != "" {
		q = q.Where("payload::text ILIKE ?", "%"+networkId+"%")
	}

	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if pageSize > 0 {
		if page < 1 {
			page = 1
		}
		q = q.Offset(int((page - 1) * pageSize)).Limit(int(pageSize))
	}

	if err := q.Order("occurred_at desc").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, count, nil
}
