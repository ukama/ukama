/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm/clause"
)

type SystemRepo interface {
	Add(ctx context.Context, sys *System) error
	Update(ctx context.Context, sys *System, org uint) error
	Delete(ctx context.Context, sys string, org uint) error
	GetByName(ctx context.Context, sys string, org uint) (*System, error)
}

type systemRepo struct {
	Db sql.Db
}

func NewSystemRepo(db sql.Db) *systemRepo {
	return &systemRepo{
		Db: db,
	}
}

// Atomic upsert on (lower(name), org_id); Uuid is excluded from DoUpdates so a
// system's original id survives repeated/concurrent registrations.
func (s *systemRepo) Add(ctx context.Context, sys *System) error {
	result := s.Db.GetGormDb().WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "lower(name)", Raw: true},
			{Name: "org_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"certificate", "api_gw_ip", "api_gw_url", "api_gw_port",
			"node_gw_ip", "node_gw_port", "node_gw_health", "api_gw_health",
			"updated_at",
		}),
	}).Create(sys)

	return result.Error
}

func (s *systemRepo) Update(ctx context.Context, sys *System, id uint) error {
	d := s.Db.GetGormDb().WithContext(ctx).Where(&System{Name: sys.Name, OrgID: id}).Updates(sys)
	return d.Error
}

func (s *systemRepo) Delete(ctx context.Context, sys string, id uint) error {
	var system System
	result := s.Db.GetGormDb().WithContext(ctx).Unscoped().Delete(&system, "name = ? and org_id = ?", strings.ToLower(sys), id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	}

	return fmt.Errorf("%s system missing", sys)

}

func (s *systemRepo) GetByName(ctx context.Context, sys string, id uint) (*System, error) {
	system := &System{}
	result := s.Db.GetGormDb().WithContext(ctx).First(system, "name = ? and org_id = ?", strings.ToLower(sys), id)
	if result.Error != nil {
		return nil, result.Error
	}
	return system, nil
}
