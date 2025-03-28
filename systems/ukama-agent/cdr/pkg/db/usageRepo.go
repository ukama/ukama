/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
 
package db

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm/clause"
)

// declare interface so that we can mock it
type UsageRepo interface {
	Add(usage *Usage) error
	Get(imsi string) (*Usage, error)
}

type usageRepo struct {
	db sql.Db
}

func NewUsageRepo(db sql.Db) *usageRepo {
	return &usageRepo{
		db: db,
	}
}

func (p *usageRepo) Add(usage *Usage) error {

	r := p.db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "imsi"}},                                                                                                                     // key colume
		DoUpdates: clause.AssignmentColumns([]string{"usage", "historical", "last_session_id", "last_session_usage", "last_node_id", "last_cdr_updated_at", "policy"}), // column needed to be updated
	}).Create(&usage)
	if r.Error != nil {
		log.Errorf("error creating usage %+v. Error: %v", usage, r.Error)
		return r.Error
	}

	return nil
}

func (p *usageRepo) Get(imsi string) (*Usage, error) {
	var usage Usage
	r := p.db.GetGormDb().Where("imsi = ?", imsi).Find(&usage)
	if r.Error != nil {
		log.Errorf("error getting usage for imsi %s.Error: %+v", imsi, r.Error)
		return nil, r.Error
	}
	return &usage, nil
}
