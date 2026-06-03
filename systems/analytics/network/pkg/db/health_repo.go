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

type HealthRepo interface {
	NetworkHealthLatest(networkId string) (*NetworkHealthRollupHourly, error)
	NetworkHealthSeries(networkId string, from, to time.Time) ([]NetworkHealthRollupHourly, error)
}

type healthRepo struct {
	Db sql.Db
}

func NewHealthRepo(db sql.Db) HealthRepo {
	return &healthRepo{
		Db: db,
	}
}

func (r *healthRepo) NetworkHealthLatest(networkId string) (*NetworkHealthRollupHourly, error) {
	var row NetworkHealthRollupHourly

	q := r.Db.GetGormDb().Order("hour desc")
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}

	result := q.First(&row)
	if result.Error != nil {
		return nil, result.Error
	}

	return &row, nil
}

func (r *healthRepo) NetworkHealthSeries(networkId string, from, to time.Time) ([]NetworkHealthRollupHourly, error) {
	var rows []NetworkHealthRollupHourly

	q := r.Db.GetGormDb().Model(&NetworkHealthRollupHourly{}).
		Where("hour >= ? AND hour < ?", from, to)
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}

	if err := q.Order("hour asc").Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}
