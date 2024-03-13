/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	"gorm.io/gorm"
)

type SiteRepo interface {
	Add(site *Site, nestedFunc func(*Site, *gorm.DB) error) error
	Get(netID,siteID uuid.UUID) (*Site, error)
	GetSites(netID uuid.UUID) ([]Site, error) 

}

type siteRepo struct {
	Db sql.Db
}

func NewSiteRepo(db sql.Db) SiteRepo {
	return &siteRepo{
		Db: db,
	}
}



func (s *siteRepo) Get(netID, siteID uuid.UUID) (*Site, error) {
    var site Site

    result := s.Db.GetGormDb().Where("network_id = ? AND id = ?", netID, siteID).First(&site)
    if result.Error != nil {
        return nil, result.Error
    }

    return &site, nil
}



func (s siteRepo) Add(site *Site, nestedFunc func(site *Site, tx *gorm.DB) error) error {
	if !validation.IsValidDnsLabelName(site.Name) {
		return fmt.Errorf("invalid name. must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(site, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(site)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (s siteRepo) GetSites(netID uuid.UUID) ([]Site, error) {
    var sites []Site
    db := s.Db.GetGormDb()

    result := db.Where("network_id = ?", netID).Find(&sites)
    if result.Error != nil {
        return nil, result.Error
    }

    return sites, nil
}
