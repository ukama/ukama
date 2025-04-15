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
	"gorm.io/gorm/clause"
)

type SiteRepo interface {
	Add(site *Site, nestedFunc func(*Site, *gorm.DB) error) error
	Get(siteId uuid.UUID) (*Site, error)
	GetSites(networkId uuid.UUID) ([]Site, error)
	Update(site *Site) error
	GetSiteCount(networkId uuid.UUID) (int64, error)
	List(networkId string, isDeactivated bool) ([]Site, error)
}

type siteRepo struct {
	Db sql.Db
}

func NewSiteRepo(db sql.Db) SiteRepo {
	return &siteRepo{
		Db: db,
	}
}

func (s *siteRepo) Get(siteId uuid.UUID) (*Site, error) {
	var site Site
	err := s.Db.GetGormDb().First(&site, siteId).Error
	if err != nil {
		return nil, err
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

func (s siteRepo) GetSites(networkId uuid.UUID) ([]Site, error) {
	var sites []Site

	err := s.Db.GetGormDb().Where("network_id = ?", networkId).Find(&sites).Error
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func (s siteRepo) List(networkId string, isDeactivated bool) ([]Site, error) {
	sites := []Site{}

	tx := s.Db.GetGormDb().Preload(clause.Associations)

	if networkId != "" {
		tx = tx.Where("network_id = ?", networkId)
	}

	result := tx.Where("is_deactivated = ?", isDeactivated).Find(&sites)
	if result.Error != nil {
		return nil, result.Error
	}

	return sites, nil
}

func (s *siteRepo) Update(site *Site) error {
	result := s.Db.GetGormDb().Model(site).Updates(site)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s siteRepo) GetSiteCount(networkId uuid.UUID) (int64, error) {
	var count int64
	result := s.Db.GetGormDb().Model(&Site{}).Where("network_id = ?", networkId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
