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

type SiteRepo interface {
	Get(siteID string) (*Site, error)
	Ensure(siteID string) error
	List() ([]Site, error)
}

type siteRepo struct{ db sql.Db }

func NewSiteRepo(db sql.Db) SiteRepo { return &siteRepo{db: db} }

func (r *siteRepo) Get(siteID string) (*Site, error) {
	var s Site
	err := r.db.GetGormDb().First(&s, "site_id = ?", siteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *siteRepo) Ensure(siteID string) error {
	return ensureSite(r.db.GetGormDb(), siteID)
}

func (r *siteRepo) List() ([]Site, error) {
	var sites []Site
	err := r.db.GetGormDb().Find(&sites).Error
	return sites, err
}

// ensureSite ensures a registry row exists for siteID before inserting child rows.
func ensureSite(tx *gorm.DB, siteID string) error {
	if siteID == "" {
		return errors.New("site_id is required")
	}
	now := time.Now().UTC()
	s := Site{SiteID: siteID, CreatedAt: now, UpdatedAt: now}
	return tx.Where("site_id = ?", siteID).FirstOrCreate(&s).Error
}
