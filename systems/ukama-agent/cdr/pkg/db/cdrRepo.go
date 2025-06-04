/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"

	log "github.com/sirupsen/logrus"
)

// declare interface so that we can mock it
type CDRRepo interface {
	Add(cdr *CDR) error
	GetByImsi(imsi string) (*[]CDR, error)
	GetBySession(imsi string, session uint64) (*[]CDR, error)
	GetByFilters(imsi string, session uint64, policy string, startTime uint64, endTime uint64) (*[]CDR, error)
	GetByPolicy(imsi string, policy string) (*[]CDR, error)
	GetByTime(imsi string, startTime uint64, endTime uint64) (*[]CDR, error)
	GetByTimeAndNodeId(imsi string, startTime uint64, endTime uint64, nodeid string) (*[]CDR, error)

	QueryUsage(imsi, nodeId string, session, from, to uint64,
		policies []string, count uint32, sort bool) (uint64, error)
}

type cdrRepo struct {
	db sql.Db
}

func NewCDRRepo(db sql.Db) *cdrRepo {
	return &cdrRepo{
		db: db,
	}
}

func (p *cdrRepo) Add(cdr *CDR) error {

	r := p.db.GetGormDb().Create(cdr)
	if r.Error != nil {
		log.Errorf("error creating cdr %+v. Error: %v", cdr, r.Error)
		return r.Error
	}

	return nil
}

func (p *cdrRepo) GetByImsi(imsi string) (*[]CDR, error) {
	var cdr []CDR
	r := p.db.GetGormDb().Where("imsi = ?", imsi).Find(&cdr)
	if r.Error != nil {
		log.Errorf("error getting cdr for imsi %s.Error: %+v", imsi, r.Error)
		return nil, r.Error
	}
	return &cdr, nil
}

func (p *cdrRepo) GetBySession(imsi string, session uint64) (*[]CDR, error) {
	var cdr []CDR
	r := p.db.GetGormDb().Where("imsi = ? AND session = ?", imsi, session).Find(&cdr)
	if r.Error != nil {
		log.Errorf("error getting cdr for imsi %s with session %d.Error: %+v", imsi, session, r.Error)
		return nil, r.Error
	}
	return &cdr, nil
}

func (p *cdrRepo) GetByFilters(imsi string, session uint64, policy string, startTime uint64, endTime uint64) (*[]CDR, error) {
	var cdr []CDR
	var ret *gorm.DB
	var count int64
	var query, sessionq, policyq string

	if policy == "" {
		policyq = "policy = policy"
	} else {
		policyq = "policy = ?"
	}

	if session == 0 {
		sessionq = "session <> ?"
	} else {
		sessionq = "session = ?"
	}

	if endTime == 0 {
		endTime = uint64(time.Now().Unix())
	}

	query = fmt.Sprintf("imsi = ? AND %s AND %s AND start_time >= ? AND start_time < ?  AND end_time <= ?", sessionq, policyq)
	if policy == "" {
		ret = p.db.GetGormDb().Where(query, imsi, session, startTime, endTime, endTime).Find(&cdr)
		if ret.Error != nil {
			log.Errorf("error getting cdr for imsi %s with start time %d and end time %d.Error: %+v", imsi, startTime, endTime, ret.Error)
			return nil, ret.Error
		}
	} else {
		ret = p.db.GetGormDb().Where(query, imsi, session, policy, startTime, endTime, endTime).Find(&cdr)
		if ret.Error != nil {
			log.Errorf("error getting cdr for imsi %s with policy %s, start time %d and end time %d.Error: %+v", imsi, policy, startTime, endTime, ret.Error)
			return nil, ret.Error
		}
	}
	_ = ret.Count(&count)
	log.Infof("%d CDR record found: %+v", count, cdr)
	return &cdr, nil
}

func (p *cdrRepo) GetByTime(imsi string, startTime uint64, endTime uint64) (*[]CDR, error) {
	var cdr []CDR
	r := p.db.GetGormDb().Where("imsi = ? AND start_time >= ? AND start_time < ? AND end_time <= ?", imsi, startTime, endTime, endTime).Find(&cdr)
	if r.Error != nil {
		log.Errorf("error getting cdr for imsi %s with start time %d and end time %d.Error: %+v", imsi, startTime, endTime, r.Error)
		return nil, r.Error
	}
	return &cdr, nil
}

func (p *cdrRepo) GetByTimeAndNodeId(imsi string, startTime uint64, endTime uint64, nodeId string) (*[]CDR, error) {
	var cdr []CDR
	r := p.db.GetGormDb().Where("imsi = ? AND start_time >= ? AND start_time < ? AND end_time <= ? AND node_id = ?", imsi, startTime, endTime, endTime, nodeId).Find(&cdr)
	if r.Error != nil {
		log.Errorf("error getting cdr for imsi %s with start time %d and end time %d.Error: %+v", imsi, startTime, endTime, r.Error)
		return nil, r.Error
	}
	return &cdr, nil
}

func (p *cdrRepo) GetByPolicy(imsi string, policy string) (*[]CDR, error) {
	var cdr []CDR
	r := p.db.GetGormDb().Where("imsi = ? AND policy = ?", imsi, policy).Find(&cdr)
	if r.Error != nil {
		log.Errorf("error getting cdr for imsi %s with policy %s.Error: %+v", imsi, policy, r.Error)
		return nil, r.Error
	}
	return &cdr, nil
}

func (p *cdrRepo) QueryUsage(imsi, nodeId string, session, from, to uint64,
	policies []string, count uint32, sort bool) (uint64, error) {
	type CDR struct {
		Imsi  string
		Total uint64
	}

	var usage CDR

	tx := p.db.GetGormDb().Preload(clause.Associations).Select("imsi, sum(total_bytes) as  total").Group("imsi")

	if imsi != "" {
		tx = tx.Where("imsi = ?", imsi)
	}

	if session > 0 {
		tx = tx.Where("session = ?", session)
	}

	if from > 0 {
		tx = tx.Where("start_time >= ?", from)
	}

	if to > 0 {
		tx = tx.Where("end_time <= ?", to)
	}

	if len(policies) > 0 {
		tx = tx.Where("policy = ?", policies[0])
	}

	if sort {
		tx = tx.Order("created_at DESC")
	}

	if count > 0 {
		tx = tx.Limit(int(count))
	}

	result := tx.First(&usage)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return 0, result.Error
	}

	if usage.Total == 0 {
		return 0, gorm.ErrInvalidData
	}

	return usage.Total, nil
}