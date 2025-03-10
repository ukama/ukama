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
)

type SimRepo interface {
	GetByIccid(iccid string) (*Sim, error)
	Add(sims []Sim) error
	GetSims() ([]Sim, error)
}

type simRepo struct {
	Db sql.Db
}

func NewSimRepo(db sql.Db) *simRepo {
	return &simRepo{
		Db: db,
	}
}

func (s *simRepo) GetSims() ([]Sim, error) {
	var sim []Sim
	result := s.Db.GetGormDb().Find(&sim)
	if result.Error != nil {
		return nil, result.Error
	}

	return sim, nil
}

func (s *simRepo) GetByIccid(iccid string) (*Sim, error) {
	var sim Sim
	result := s.Db.GetGormDb().Where("iccid = ?", iccid).First(&sim)

	if result.Error != nil {
		return nil, result.Error
	}

	return &sim, nil
}

func (s *simRepo) Add(sims []Sim) error {
	e := s.Db.GetGormDb().Create(&sims)
	if e != nil {
		return e.Error
	}

	return nil
}
