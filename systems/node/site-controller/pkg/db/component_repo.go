/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"errors"
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type ComponentRepo interface {
	Get(siteID string) (*SiteComponent, error)
	Upsert(component *SiteComponent) error
}

type componentRepo struct{ db sql.Db }

func NewComponentRepo(db sql.Db) ComponentRepo { return &componentRepo{db: db} }

func (r *componentRepo) Get(siteID string) (*SiteComponent, error) {
	var m SiteComponent
	err := r.db.GetGormDb().First(&m, "site_id = ?", siteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *componentRepo) Upsert(m *SiteComponent) error {
	if m == nil {
		return errors.New("component is required")
	}
	if m.SiteID == "" {
		return errors.New("site_id is required")
	}

	gdb := r.db.GetGormDb()
	if err := ensureSite(gdb, m.SiteID); err != nil {
		return err
	}
	now := time.Now().UTC()

	var existing SiteComponent
	err := gdb.First(&existing, "site_id = ?", m.SiteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		row := *m
		if row.ID == uuid.Nil {
			row.ID = uuid.NewV4()
		}
		row.UpdatedAt = now
		return gdb.Create(&row).Error
	}
	if err != nil {
		return err
	}

	existing.Components = m.Components
	existing.UpdatedAt = now
	return gdb.Save(&existing).Error
}
