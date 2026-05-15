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
	if m == nil {
		return errors.New("state is required")
	}
	gdb := r.db.GetGormDb()
	if err := ensureSite(gdb, m.SiteID); err != nil {
		return err
	}
	now := time.Now().UTC()
	m.UpdatedAt = now

	var existing SiteState
	err := gdb.First(&existing, "site_id = ?", m.SiteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if m.ID == uuid.Nil {
			m.ID = uuid.NewV4()
		}
		return gdb.Create(m).Error
	}
	if err != nil {
		return err
	}

	existing.PowerState = m.PowerState
	existing.ServiceState = m.ServiceState
	existing.RadioState = m.RadioState
	existing.AccessState = m.AccessState
	existing.Reason = m.Reason
	existing.UpdatedAt = now
	return gdb.Save(&existing).Error
}
