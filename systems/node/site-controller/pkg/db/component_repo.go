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
	db := r.db.GetGormDb()
	if err := ensureSite(db, m.SiteID); err != nil {
		return err
	}
	m.UpdatedAt = time.Now().UTC()
	return db.Save(m).Error
}
