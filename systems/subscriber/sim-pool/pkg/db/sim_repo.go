/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/gorm"
)

type SimRepo interface {
	Get(isPhysicalSim bool, simType ukama.SimType) (*Sim, error)
	GetByIccid(iccid string) (*Sim, error)
	GetSimsByType(simType ukama.SimType) ([]Sim, error)
	Add(sims []Sim) error
	Delete(id []uint64) error
	UpdateStatus(iccid string, isAllocated, IsFailed bool) error
	GetSims(simType ukama.SimType) ([]Sim, error)
}

type simRepo struct {
	Db sql.Db
}

func NewSimRepo(db sql.Db) *simRepo {
	return &simRepo{
		Db: db,
	}
}

func (s *simRepo) GetSims(simType ukama.SimType) ([]Sim, error) {
	var sim []Sim
	if simType != ukama.SimTypeUnknown {
		result := s.Db.GetGormDb().Where("sim_type = ?", simType).Find(&sim)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		result := s.Db.GetGormDb().Find(&sim)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return sim, nil

}

func (s *simRepo) Get(isPhysicalSim bool, simType ukama.SimType) (*Sim, error) {
	var sim Sim
	result := s.Db.GetGormDb().Where("is_allocated = ?", false).Where("is_physical = ?", isPhysicalSim).Where("sim_type = ?", simType).First(&sim)

	if result.Error != nil {
		return nil, result.Error
	}

	return &sim, nil
}

func (s *simRepo) UpdateStatus(iccid string, isAllocated, IsFailed bool) error {
	err := s.Db.GetGormDb().Model(&Sim{}).Where("iccid = ?", iccid).Updates(map[string]interface{}{
		"is_allocated": isAllocated,
		"is_failed":    IsFailed,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func (s *simRepo) GetByIccid(iccid string) (*Sim, error) {
	var sim Sim
	result := s.Db.GetGormDb().Where("is_allocated = ? AND iccid = ?", false, iccid).First(&sim)

	if result.Error != nil {
		return nil, result.Error
	}

	return &sim, nil
}

func (s *simRepo) GetSimsByType(simType ukama.SimType) ([]Sim, error) {
	var sim []Sim
	result := s.Db.GetGormDb().Where("sim_type = ?", simType).Find(&sim)

	if result.Error != nil {
		return nil, result.Error
	}

	return sim, nil
}

func (s *simRepo) Add(sims []Sim) error {
	e := s.Db.GetGormDb().Create(&sims)
	if e != nil {
		return e.Error
	}

	return nil
}

func (s *simRepo) Delete(Id []uint64) error {
	result := s.Db.GetGormDb().Where("id IN (?)", Id).Delete(&Sim{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
