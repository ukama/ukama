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

type IntentFlightRepo interface {
	GetBySiteIntentID(siteIntentID uuid.UUID) (*SiteIntentFlight, error)
	Upsert(flight *SiteIntentFlight) error
}

type intentFlightRepo struct{ db sql.Db }

func NewIntentFlightRepo(db sql.Db) IntentFlightRepo { return &intentFlightRepo{db: db} }

func (r *intentFlightRepo) GetBySiteIntentID(siteIntentID uuid.UUID) (*SiteIntentFlight, error) {
	var m SiteIntentFlight
	err := r.db.GetGormDb().Where("site_intent_id = ?", siteIntentID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *intentFlightRepo) Upsert(m *SiteIntentFlight) error {
	if m == nil {
		return errors.New("flight is required")
	}
	if m.SiteIntentID == uuid.Nil {
		return errors.New("site_intent_id is required")
	}

	gdb := r.db.GetGormDb()
	now := time.Now().UTC()

	var existing SiteIntentFlight
	err := gdb.Where("site_intent_id = ?", m.SiteIntentID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if m.ID == uuid.Nil {
			m.ID = uuid.NewV4()
		}
		if m.CreatedAt.IsZero() {
			m.CreatedAt = now
		}
		m.UpdatedAt = now
		return gdb.Create(m).Error
	}
	if err != nil {
		return err
	}

	existing.IntentFlight = m.IntentFlight
	existing.ExpiresAt = m.ExpiresAt
	existing.UpdatedAt = now
	return gdb.Save(&existing).Error
}
