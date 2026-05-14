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

type StateRepo interface {
	Get(siteID string) (*SiteState, error)
	Upsert(state *SiteState) error
}
type stateRepo struct{ db sql.Db }

func NewStateRepo(db sql.Db) StateRepo { return &stateRepo{db: db} }
func (r *stateRepo) Get(siteID string) (*SiteState, error) {
	var m SiteState
	err := r.db.GetGormDb().First(&m, "site_id = ?", siteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func (r *stateRepo) Upsert(m *SiteState) error {
	db := r.db.GetGormDb()
	if err := ensureSite(db, m.SiteID); err != nil {
		return err
	}
	m.UpdatedAt = time.Now().UTC()
	return db.Save(m).Error
}
