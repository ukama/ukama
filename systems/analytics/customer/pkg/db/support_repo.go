/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
)

// SupportRepo is a read-only repository over the analytics database
// providing the inputs needed for support diagnosis.
type SupportRepo interface {
	RecentActivityFor(customerId uuid.UUID, limit int) ([]EventLog, error)
	SiteHealth(siteId uuid.UUID) (*SiteHealthRollupHourly, error)
}

type supportRepo struct {
	Db sql.Db
}

func NewSupportRepo(db sql.Db) SupportRepo {
	return &supportRepo{
		Db: db,
	}
}

// RecentActivityFor returns the most recent event log entries whose payload
// references the given customer id.
func (r supportRepo) RecentActivityFor(customerId uuid.UUID, limit int) ([]EventLog, error) {
	var logs []EventLog

	if limit <= 0 {
		limit = 10
	}

	result := r.Db.GetGormDb().Model(&EventLog{}).
		Where("payload::text ILIKE ?", "%"+customerId.String()+"%").
		Order("occurred_at desc").
		Limit(limit).
		Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}

	return logs, nil
}

// SiteHealth returns the latest hourly health rollup for the given site.
func (r supportRepo) SiteHealth(siteId uuid.UUID) (*SiteHealthRollupHourly, error) {
	var health SiteHealthRollupHourly

	result := r.Db.GetGormDb().
		Where("site_id = ?", siteId).
		Order("hour desc").
		First(&health)
	if result.Error != nil {
		return nil, result.Error
	}

	return &health, nil
}
