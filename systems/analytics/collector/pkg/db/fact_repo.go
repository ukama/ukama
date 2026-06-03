/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FactRepo interface {
	AddPaymentEvent(e *PaymentEvent) error
	AddUsageEvent(e *UsageEvent) error
	AddMetricSample(e *MetricSample) error
	AddAlarmEvent(e *AlarmEvent) error
	AddNodeStateEvent(e *NodeStateEvent) error
	AddSiteStateEvent(e *SiteStateEvent) error
	AddCustomerEvent(e *CustomerEvent) error
	AddSimEvent(e *SimEvent) error
	AddPackageEvent(e *PackageEvent) error
	AddInventoryEvent(e *InventoryEvent) error

	/* Interval helpers, derived from state events. */
	TransitionNodeState(nodeId, state string, at time.Time) error
	TransitionSiteState(siteId uuid.UUID, state string, at time.Time) error
	TransitionSimState(simId, state string, at time.Time) error
	OpenCustomerPackageInterval(customerId, packageId uuid.UUID, state string, at time.Time) error
	CloseCustomerPackageInterval(customerId uuid.UUID, at time.Time) error
}

type factRepo struct {
	Db sql.Db
}

func NewFactRepo(db sql.Db) FactRepo {
	return &factRepo{
		Db: db,
	}
}

func (r *factRepo) AddPaymentEvent(e *PaymentEvent) error {
	/* ExternalId is the payment processor id: dedupe on redelivery. */
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoNothing: true,
	}).Create(e)

	return result.Error
}

func (r *factRepo) AddUsageEvent(e *UsageEvent) error {
	return r.Db.GetGormDb().Create(e).Error
}

func (r *factRepo) AddMetricSample(e *MetricSample) error {
	return r.Db.GetGormDb().Create(e).Error
}

func (r *factRepo) AddAlarmEvent(e *AlarmEvent) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "alarm_id"}},
		UpdateAll: true,
	}).Create(e)

	return result.Error
}

func (r *factRepo) AddNodeStateEvent(e *NodeStateEvent) error {
	return r.Db.GetGormDb().Create(e).Error
}

func (r *factRepo) AddSiteStateEvent(e *SiteStateEvent) error {
	return r.Db.GetGormDb().Create(e).Error
}

func (r *factRepo) AddCustomerEvent(e *CustomerEvent) error {
	return r.Db.GetGormDb().Create(e).Error
}

func (r *factRepo) AddSimEvent(e *SimEvent) error {
	return r.Db.GetGormDb().Create(e).Error
}

func (r *factRepo) AddPackageEvent(e *PackageEvent) error {
	return r.Db.GetGormDb().Create(e).Error
}

func (r *factRepo) AddInventoryEvent(e *InventoryEvent) error {
	return r.Db.GetGormDb().Create(e).Error
}

// TransitionNodeState closes the currently open node state interval (if any)
// and opens a new one with the given state.
func (r *factRepo) TransitionNodeState(nodeId, state string, at time.Time) error {
	return r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&NodeStateInterval{}).
			Where("node_id = ? AND end_at IS NULL", nodeId).
			Updates(map[string]interface{}{
				"end_at":           at,
				"duration_seconds": gorm.Expr("EXTRACT(EPOCH FROM (?::timestamptz - start_at))", at),
			}).Error; err != nil {
			return err
		}

		return tx.Create(&NodeStateInterval{
			NodeId:  nodeId,
			State:   state,
			StartAt: at,
		}).Error
	})
}

// TransitionSiteState closes the currently open site state interval (if any)
// and opens a new one with the given state.
func (r *factRepo) TransitionSiteState(siteId uuid.UUID, state string, at time.Time) error {
	return r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&SiteStateInterval{}).
			Where("site_id = ? AND end_at IS NULL", siteId).
			Updates(map[string]interface{}{
				"end_at":           at,
				"duration_seconds": gorm.Expr("EXTRACT(EPOCH FROM (?::timestamptz - start_at))", at),
			}).Error; err != nil {
			return err
		}

		return tx.Create(&SiteStateInterval{
			SiteId:  siteId,
			State:   state,
			StartAt: at,
		}).Error
	})
}

// TransitionSimState closes the currently open sim state interval (if any)
// and opens a new one with the given state.
func (r *factRepo) TransitionSimState(simId, state string, at time.Time) error {
	return r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&SimStateInterval{}).
			Where("sim_id = ? AND end_at IS NULL", simId).
			Update("end_at", at).Error; err != nil {
			return err
		}

		return tx.Create(&SimStateInterval{
			SimId:   simId,
			State:   state,
			StartAt: at,
		}).Error
	})
}

// OpenCustomerPackageInterval closes any open package interval for the
// customer, then opens a new one for the given package.
func (r *factRepo) OpenCustomerPackageInterval(customerId, packageId uuid.UUID, state string, at time.Time) error {
	return r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&CustomerPackageInterval{}).
			Where("customer_id = ? AND end_at IS NULL", customerId).
			Update("end_at", at).Error; err != nil {
			return err
		}

		return tx.Create(&CustomerPackageInterval{
			CustomerId: customerId,
			PackageId:  packageId,
			State:      state,
			StartAt:    at,
		}).Error
	})
}

// CloseCustomerPackageInterval closes any open package interval for the
// customer without opening a new one.
func (r *factRepo) CloseCustomerPackageInterval(customerId uuid.UUID, at time.Time) error {
	return r.Db.GetGormDb().Model(&CustomerPackageInterval{}).
		Where("customer_id = ? AND end_at IS NULL", customerId).
		Update("end_at", at).Error
}
