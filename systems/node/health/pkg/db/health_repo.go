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
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/gorm"
)

type HealthRepo interface {
	StoreRunningAppsInfo(health *Health, nestedFunc func(string, string) error) error
	List(id string, nodeId string, timestamp string, filter ukama.FilterTimeframesType) ([]*Health, error)
}
type healthRepo struct {
	Db sql.Db
}

func NewHealthRepo(db sql.Db) HealthRepo {
	return &healthRepo{
		Db: db,
	}
}

func (r *healthRepo) List(id string, nodeId string, timestamp string, filter ukama.FilterTimeframesType) ([]*Health, error) {
	query := r.Db.GetGormDb().
		Preload("System").
		Preload("Capps.Resources").
		Order("created_at DESC")

	if id != "" {
		query = query.Where("id = ?", id)
	}

	if timestamp != "" {
		query = query.Where("time_stamp = ?", timestamp)
	}

	if nodeId != "" {
		query = query.Where("node_id = ?", nodeId)
	}

	if filter == ukama.FilterTimeframesTypeLatest {
		var health Health
		result := query.Limit(1).First(&health)
		if result.Error != nil {
			return nil, result.Error
		}
		return []*Health{&health}, nil
	}

	var healths []*Health
	result := query.Find(&healths)
	return healths, result.Error
}

func (r *healthRepo) StoreRunningAppsInfo(health *Health, nestedFunc func(string, string) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc("", "")
			if nestErr != nil {
				return nestErr
			}
		}
		if err := tx.Create(health).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
