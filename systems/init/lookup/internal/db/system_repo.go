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

func (s *systemRepo) Add(ctx context.Context, sys *System) error {
	// result := s.Db.GetGormDb().Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "name"}, {Name: "org_id"}},
	// 	UpdateAll: true,
	// }).Create(sys)
	t := s.Db.GetGormDb().WithContext(ctx).Model(&System{}).Where("org_id = ? AND name = ?", sys.OrgID, sys.Name).Updates(sys)
	if t.Error != nil {
		return t.Error
	} else {
		if t.RowsAffected == 0 {
			err := s.Db.GetGormDb().WithContext(ctx).Create(sys).Error // create new record from newUser
			return err
		}
	}

	return nil
}

func (s *systemRepo) Update(ctx context.Context, sys *System, id uint) error {
	d := s.Db.GetGormDb().WithContext(ctx).Preload("Org").Where(&System{Name: sys.Name, OrgID: id}).Updates(sys)
	return d.Error
}

func (s *systemRepo) Delete(ctx context.Context, sys string, id uint) error {
	var system System
	result := s.Db.GetGormDb().WithContext(ctx).Preload("Org").Unscoped().Delete(&system, "name = ? and org_id = ?", strings.ToLower(sys), id)
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
	result := s.Db.GetGormDb().WithContext(ctx).Preload("Org").First(system, "name = ? and org_id = ?", strings.ToLower(sys), id)
	if result.Error != nil {
		return nil, result.Error
	}
	return system, nil
}
