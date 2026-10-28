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
	Delete(siteID string) error
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

func (r *siteRepo) Delete(siteID string) error {
	if siteID == "" {
		return errors.New("site_id is required")
	}
	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		intentIDs := tx.Model(&SiteIntent{}).Select("id").Where("site_id = ?", siteID)
		if err := tx.Where("site_intent_id IN (?)", intentIDs).Delete(&SiteIntentFlight{}).Error; err != nil {
			return err
		}
		for _, m := range []interface{}{&SiteIntent{}, &SiteState{}, &SiteComponent{}, &SitePortMap{}, &Site{}} {
			if err := tx.Where("site_id = ?", siteID).Delete(m).Error; err != nil {
				return err
			}
		}
		return nil
	})
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
