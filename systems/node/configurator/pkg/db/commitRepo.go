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

	"gorm.io/gorm/clause"
)

type CommitRepo interface {
	Add(hash string) error
	Get(hash string) (*Commit, error)
	GetAll() ([]Commit, error)
	GetLatest() (*Commit, error)
}

type commitRepo struct {
	Db sql.Db
}

func NewCommitRepo(db sql.Db) CommitRepo {
	return &commitRepo{
		Db: db,
	}
}

func (n *commitRepo) Add(hash string) error {
	commit := Commit{
		Hash: hash,
	}

	r := n.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "hash"}},
		DoNothing: true,
	}).Create(&commit)

	return r.Error
}

func (n *commitRepo) Get(hash string) (*Commit, error) {
	var commit Commit

	result := n.Db.GetGormDb().First(&commit, "hash=?", hash)
	if result.Error != nil {
		return nil, result.Error
	}

	return &commit, nil
}

func (n *commitRepo) GetAll() ([]Commit, error) {
	var commit []Commit

	result := n.Db.GetGormDb().Find(&commit)

	if result.Error != nil {
		return nil, result.Error
	}

	return commit, nil
}

func (n *commitRepo) GetLatest() (*Commit, error) {
	var commit Commit

	result := n.Db.GetGormDb().Last(&commit)
	if result.Error != nil {
		return nil, result.Error
	}

	return &commit, nil
}
