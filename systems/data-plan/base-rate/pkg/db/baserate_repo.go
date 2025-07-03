/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"time"

	log "github.com/sirupsen/logrus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
	ukama "github.com/ukama/ukama/systems/common/ukama"
)

type BaseRateRepo interface {
	GetBaseRateById(uuid uuid.UUID) (*BaseRate, error)
	GetBaseRatesHistoryByCountry(country, provider string, sType ukama.SimType) ([]BaseRate, error)
	GetBaseRatesByCountry(country, provider string, simType ukama.SimType) ([]BaseRate, error)
	GetBaseRatesForPeriod(country, provider string, from, to time.Time, simType ukama.SimType) ([]BaseRate, error)
	GetBaseRatesForPackage(country, provider string, from, to time.Time, simType ukama.SimType) ([]BaseRate, error)
	UploadBaseRates(rateList []BaseRate) error
}

type baseRateRepo struct {
	Db sql.Db
}

func NewBaseRateRepo(db sql.Db) *baseRateRepo {
	return &baseRateRepo{
		Db: db,
	}
}

func (u *baseRateRepo) GetBaseRateById(uuid uuid.UUID) (*BaseRate, error) {
	rate := &BaseRate{}
	result := u.Db.GetGormDb().First(rate, "uuid=?", uuid)
	if result.Error != nil {
		return nil, result.Error
	}
	return rate, nil
}

func (b *baseRateRepo) GetBaseRatesHistoryByCountry(country, provider string, sType ukama.SimType) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Unscoped().Where(&BaseRate{Country: country, Provider: provider, SimType: sType}).Find(&rates)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) GetBaseRatesByCountry(country, provider string, simType ukama.SimType) ([]BaseRate, error) {
	var rates []BaseRate
	t := time.Now().Add(time.Second * 1).Format(time.RFC3339)
	result := b.Db.GetGormDb().Model(BaseRate{}).Where("country = ?", country).Where("provider = ?", provider).
		Where("sim_type = ?", simType).Where("effective_at <= ?", t).Order("effective_at desc").Limit(1).Find(&rates)
	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) GetBaseRatesForPeriod(country, provider string, from, to time.Time, simType ukama.SimType) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Model(BaseRate{}).Unscoped().Where("country = ?", country).Where("provider = ?", provider).
		Where("sim_type = ?", simType).Where("effective_at <= ?", from).Where("end_at >= ?", to).Find(&rates)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) GetBaseRatesForPackage(country, provider string, from, to time.Time, simType ukama.SimType) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Model(BaseRate{}).Unscoped().Where("country = ?", country).Where("provider = ?", provider).
		Where("sim_type = ?", simType).Where("effective_at <= ?", from).Where("end_at >= ?", to).Order("created_at desc").Limit(1).Find(&rates)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) UploadBaseRates(rateList []BaseRate) error {
	err := b.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		for _, r := range rateList {
			o := &BaseRate{
				Country:     r.Country,
				Provider:    r.Provider,
				EffectiveAt: r.EffectiveAt,
				SimType:     r.SimType,
			}
			result := tx.Model(BaseRate{}).Where("country = ?", r.Country).Where("provider = ?", r.Provider).
				Where("sim_type = ?", r.SimType).Where("effective_at = ?", r.EffectiveAt).Delete(o)
			if result.Error != nil {
				if !sql.IsNotFoundError(result.Error) {
					log.Errorf("Error deleting rate %+v . Error %s", o, result.Error.Error())
					return result.Error
				}
			}

			result = tx.Model(BaseRate{}).Create(&r)
			if result.Error != nil {
				log.Errorf("Error creating rate %+v . Error %s", r, result.Error.Error())
				return result.Error
			}
		}

		return nil
	})

	return err
}
