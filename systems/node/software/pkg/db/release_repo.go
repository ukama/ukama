/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReleaseRepo interface {
	// Release catalog (availability)
	Upsert(r *ReleaseCatalog) error
	SetChunked(name, rtype, version string) error
	Get(name, rtype, version string) (*ReleaseCatalog, error)
	Exists(name, rtype, version string) (bool, error)
	List(name, rtype string) ([]ReleaseCatalog, error)

	// Desired release (promotion)
	GetDesired(name, rtype string) (*AppDesiredRelease, error)
	SetDesired(d *AppDesiredRelease) error
	ListDesired() ([]AppDesiredRelease, error)
}

type releaseRepo struct {
	Db sql.Db
}

func NewReleaseRepo(db sql.Db) ReleaseRepo {
	return &releaseRepo{Db: db}
}

func defType(t string) string {
	if t == "" {
		return "app"
	}
	return t
}

// Upsert records availability. It deliberately does NOT touch `chunked`
// (owned by SetChunked) so a re-announced upload can't reset the chunk flag.
func (r *releaseRepo) Upsert(rel *ReleaseCatalog) error {
	if rel.Id == uuid.Nil {
		rel.Id = uuid.NewV4()
	}
	rel.Type = defType(rel.Type)
	return r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "type"}, {Name: "version"}},
		DoUpdates: clause.AssignmentColumns([]string{"digest", "size_bytes", "available", "uploaded_at", "updated_at"}),
	}).Create(rel).Error
}

func (r *releaseRepo) SetChunked(name, rtype, version string) error {
	return r.Db.GetGormDb().Model(&ReleaseCatalog{}).
		Where("name = ? AND type = ? AND version = ?", name, defType(rtype), version).
		Update("chunked", true).Error
}

func (r *releaseRepo) Get(name, rtype, version string) (*ReleaseCatalog, error) {
	var rel ReleaseCatalog
	err := r.Db.GetGormDb().
		Where("name = ? AND type = ? AND version = ?", name, defType(rtype), version).
		First(&rel).Error
	if err != nil {
		return nil, err
	}
	return &rel, nil
}

func (r *releaseRepo) Exists(name, rtype, version string) (bool, error) {
	var count int64
	err := r.Db.GetGormDb().Model(&ReleaseCatalog{}).
		Where("name = ? AND type = ? AND version = ?", name, defType(rtype), version).
		Count(&count).Error
	return count > 0, err
}

func (r *releaseRepo) List(name, rtype string) ([]ReleaseCatalog, error) {
	var out []ReleaseCatalog
	tx := r.Db.GetGormDb().Model(&ReleaseCatalog{})
	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if rtype != "" {
		tx = tx.Where("type = ?", rtype)
	}
	err := tx.Order("name, version").Find(&out).Error
	return out, err
}

func (r *releaseRepo) GetDesired(name, rtype string) (*AppDesiredRelease, error) {
	var d AppDesiredRelease
	err := r.Db.GetGormDb().
		Where("name = ? AND type = ?", name, defType(rtype)).
		First(&d).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // no desired set yet — not an error
		}
		return nil, err
	}
	return &d, nil
}

func (r *releaseRepo) SetDesired(d *AppDesiredRelease) error {
	if d.Id == uuid.Nil {
		d.Id = uuid.NewV4()
	}
	d.Type = defType(d.Type)
	return r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "type"}},
		DoUpdates: clause.AssignmentColumns([]string{"desired_version", "promoted_at", "promoted_by", "updated_at"}),
	}).Create(d).Error
}

func (r *releaseRepo) ListDesired() ([]AppDesiredRelease, error) {
	var out []AppDesiredRelease
	err := r.Db.GetGormDb().Find(&out).Error
	return out, err
}
