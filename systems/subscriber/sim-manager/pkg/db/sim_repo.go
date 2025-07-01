/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
)

type SimRepo interface {
	Add(sim *Sim, nestedFunc func(*Sim, *gorm.DB) error) error
	Get(simID uuid.UUID) (*Sim, error)

	// Deprecated: Use db.SimRepo.List with iccid as filtering param instead.
	GetByIccid(iccid string) (*Sim, error)

	// Deprecated: Use db.SimRepo.List with subscriberId as filtering param instead.
	GetBySubscriber(subscriberID uuid.UUID) ([]Sim, error)

	// Deprecated: Use db.SimRepo.List with networkId as filtering param instead.
	GetByNetwork(networkID uuid.UUID) ([]Sim, error)

	List(iccid, imsi, SubscriberId, networkId string, simType ukama.SimType, status ukama.SimStatus,
		TrafficPolicy uint32, IsPhysical bool, count uint32, sort bool) ([]Sim, error)

	Update(sim *Sim, nestedFunc func(*Sim, *gorm.DB) error) error
	Delete(simID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error
}

type simRepo struct {
	Db sql.Db
}

func NewSimRepo(db sql.Db) SimRepo {
	return &simRepo{
		Db: db,
	}
}

func (s *simRepo) Add(sim *Sim, nestedFunc func(sim *Sim, tx *gorm.DB) error) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(sim, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		log.Info("Adding sim", sim)
		result := tx.Create(sim)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (s *simRepo) Get(simID uuid.UUID) (*Sim, error) {
	var sim Sim

	result := s.Db.GetGormDb().Model(&Sim{}).
		Preload("Package", "is_active is true").First(&sim, simID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &sim, nil
}

// Deprecated: Use db.SimRepo.List with iccid as filtering param instead.
func (s *simRepo) GetByIccid(iccid string) (*Sim, error) {
	var sim Sim

	result := s.Db.GetGormDb().Model(&Sim{}).Where(&Sim{Iccid: iccid}).
		Preload("Package", "is_active is true").First(&sim)

	if result.Error != nil {
		return nil, result.Error
	}

	return &sim, nil
}

// Deprecated: Use db.SimRepo.List with subscriberId as filtering param instead.
func (s *simRepo) GetBySubscriber(subscriberID uuid.UUID) ([]Sim, error) {
	var sims []Sim

	result := s.Db.GetGormDb().Model(&Sim{}).Where(&Sim{SubscriberId: subscriberID}).
		Preload("Package", "is_active is true").Find(&sims)

	if result.Error != nil {
		return nil, result.Error
	}

	return sims, nil
}

// Deprecated: Use db.SimRepo.List with networkId as filtering param instead.
func (s *simRepo) GetByNetwork(networkID uuid.UUID) ([]Sim, error) {
	var sims []Sim

	result := s.Db.GetGormDb().Model(&Sim{}).Where(&Sim{NetworkId: networkID}).
		Preload("Package", "is_active is true").Find(&sims)

	if result.Error != nil {
		return nil, result.Error
	}

	return sims, nil
}

func (r *simRepo) List(iccid, imsi, subscriberId, networkId string, simType ukama.SimType, simStatus ukama.SimStatus,
	trafficPolicy uint32, IsPhysical bool, count uint32, sort bool) ([]Sim, error) {

	sims := []Sim{}

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	if iccid != "" {
		tx = tx.Where("iccid = ?", iccid)
	}

	if imsi != "" {
		tx = tx.Where("imsi = ?", imsi)
	}

	if subscriberId != "" {
		tx = tx.Where("subscriber_id = ?", subscriberId)
	}

	if networkId != "" {
		tx = tx.Where("network_id = ?", networkId)
	}

	if simType != ukama.SimTypeUnknown {
		tx = tx.Where("type = ?", simType)
	}

	if simStatus != ukama.SimStatusUnknown {
		tx = tx.Where("status = ?", simStatus)
	}

	if trafficPolicy > 0 {
		tx = tx.Where("traffic_policy = ?", trafficPolicy)
	}

	if IsPhysical {
		tx = tx.Where("is_physical = ?", true)
	}

	if sort {
		tx = tx.Order("allocated_at DESC")
	}

	if count > 0 {
		tx = tx.Limit(int(count))
	}

	result := tx.Preload("Package", "is_active is true").Find(&sims)
	if result.Error != nil {
		return nil, result.Error
	}

	return sims, nil
}

// Update package modified non-empty fields provided by Package struct
func (s *simRepo) Update(sim *Sim, nestedFunc func(*Sim, *gorm.DB) error) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(sim, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Clauses(clause.Returning{}).Updates(sim)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})

	return err
}

func (s *simRepo) Delete(simID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Sim{}, simID)
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(simID, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}
