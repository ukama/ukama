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

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SnapshotRepo interface {
	UpsertNetwork(s *NetworkSnapshot) error
	UpsertSite(s *SiteSnapshot) error
	UpsertNode(s *NodeSnapshot) error
	UpdateNodeStatus(nodeId, status string, at time.Time) error
	UpsertCustomer(s *CustomerSnapshot) error
	DeleteCustomer(customerId string) error
	UpsertSim(s *SimSnapshot) error
	DeleteSim(simId string) error
	UpsertSimBatch(s *SimBatchSnapshot) error
	UpsertPackage(s *PackageSnapshot) error
	DeletePackage(packageId string) error
	UpsertInventory(s *InventorySnapshot) error
	UpsertBilling(s *BillingSnapshot) error
	UpsertHealthReport(s *HealthReportSnapshot) error
}

type snapshotRepo struct {
	Db gormHandle
}

func NewSnapshotRepo(db sql.Db) SnapshotRepo {
	return &snapshotRepo{
		Db: db,
	}
}

func NewSnapshotRepoWithGorm(db *gorm.DB) SnapshotRepo {
	return &snapshotRepo{
		Db: gormOnly{db: db},
	}
}

func (r *snapshotRepo) UpsertNetwork(s *NetworkSnapshot) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "network_id"}},
		UpdateAll: true,
	}).Create(s)

	return result.Error
}

func (r *snapshotRepo) UpsertSite(s *SiteSnapshot) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "site_id"}},
		UpdateAll: true,
	}).Create(s)

	return result.Error
}

func (r *snapshotRepo) UpsertNode(s *NodeSnapshot) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		UpdateAll: true,
	}).Create(s)

	return result.Error
}

func (r *snapshotRepo) UpdateNodeStatus(nodeId, status string, at time.Time) error {
	result := r.Db.GetGormDb().Model(&NodeSnapshot{}).Where("node_id = ?", nodeId).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": at,
		})

	return result.Error
}

func (r *snapshotRepo) UpsertCustomer(s *CustomerSnapshot) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "customer_id"}},
		UpdateAll: true,
	}).Create(s)

	return result.Error
}

func (r *snapshotRepo) DeleteCustomer(customerId string) error {
	result := r.Db.GetGormDb().Where("customer_id = ?", customerId).
		Delete(&CustomerSnapshot{})

	return result.Error
}

func (r *snapshotRepo) UpsertSim(s *SimSnapshot) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "sim_id"}},
		UpdateAll: true,
	}).Create(s)

	return result.Error
}

func (r *snapshotRepo) DeleteSim(simId string) error {
	result := r.Db.GetGormDb().Where("sim_id = ?", simId).Delete(&SimSnapshot{})

	return result.Error
}

func (r *snapshotRepo) UpsertSimBatch(s *SimBatchSnapshot) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "batch_id"}},
		UpdateAll: true,
	}).Create(s)

	return result.Error
}

func (r *snapshotRepo) UpsertPackage(s *PackageSnapshot) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "package_id"}},
		UpdateAll: true,
	}).Create(s)

	return result.Error
}

func (r *snapshotRepo) DeletePackage(packageId string) error {
	result := r.Db.GetGormDb().Where("package_id = ?", packageId).
		Delete(&PackageSnapshot{})

	return result.Error
}

func (r *snapshotRepo) UpsertInventory(s *InventorySnapshot) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "component_id"}},
		UpdateAll: true,
	}).Create(s)

	return result.Error
}

func (r *snapshotRepo) UpsertBilling(s *BillingSnapshot) error {
	/* Org-level singleton row. */
	if s.Id == 0 {
		s.Id = 1
	}

	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(s)

	return result.Error
}

func (r *snapshotRepo) UpsertHealthReport(s *HealthReportSnapshot) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		UpdateAll: true,
	}).Create(s)

	return result.Error
}
