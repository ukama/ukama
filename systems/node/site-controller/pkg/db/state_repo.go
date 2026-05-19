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
	Upsert(patch *SiteState) error
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

func (r *stateRepo) Upsert(patch *SiteState) error {
	if patch == nil {
		return errors.New("state is required")
	}
	if patch.SiteID == "" {
		return errors.New("site_id is required")
	}

	gdb := r.db.GetGormDb()
	if err := ensureSite(gdb, patch.SiteID); err != nil {
		return err
	}
	now := time.Now().UTC()

	var existing SiteState
	err := gdb.First(&existing, "site_id = ?", patch.SiteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		row := *patch
		if row.ID == uuid.Nil {
			row.ID = uuid.NewV4()
		}
		row.UpdatedAt = now
		cols := siteStatePatchColumns(patch)
		cols = append([]string{"id", "site_id", "updated_at"}, cols...)
		return gdb.Select(cols).Create(&row).Error
	}
	if err != nil {
		return err
	}

	updates := siteStatePatchValues(patch, now)
	return gdb.Model(&existing).Updates(updates).Error
}

func siteStatePatchColumns(patch *SiteState) []string {
	var cols []string
	if patch.PowerState != "" {
		cols = append(cols, "power_state")
	}
	if patch.ServiceState != "" {
		cols = append(cols, "service_state")
	}
	if patch.RadioState != "" {
		cols = append(cols, "radio_state")
	}
	if patch.AccessState != "" {
		cols = append(cols, "access_state")
	}
	if patch.Reason != "" {
		cols = append(cols, "reason")
	}
	return cols
}

func siteStatePatchValues(patch *SiteState, now time.Time) map[string]interface{} {
	updates := map[string]interface{}{"updated_at": now}
	if patch.PowerState != "" {
		updates["power_state"] = patch.PowerState
	}
	if patch.ServiceState != "" {
		updates["service_state"] = patch.ServiceState
	}
	if patch.RadioState != "" {
		updates["radio_state"] = patch.RadioState
	}
	if patch.AccessState != "" {
		updates["access_state"] = patch.AccessState
	}
	if patch.Reason != "" {
		updates["reason"] = patch.Reason
	}
	return updates
}
