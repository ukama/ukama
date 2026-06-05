/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"errors"
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type IntentRepo interface {
	Add(intent *MonitoredIntent) error
	Get(operationId uuid.UUID) (*MonitoredIntent, error)
	FindWatchingByResource(resourceKey string) ([]MonitoredIntent, error)
	MarkTerminal(operationId uuid.UUID, status IntentStatus) (*MonitoredIntent, error)
	FindExpired(now time.Time, limit int) ([]MonitoredIntent, error)
}

type intentRepo struct {
	db sql.Db
}

func NewIntentRepo(db sql.Db) IntentRepo {
	return &intentRepo{db: db}
}

func (r *intentRepo) Add(intent *MonitoredIntent) error {
	return r.db.GetGormDb().Create(intent).Error
}

func (r *intentRepo) Get(operationId uuid.UUID) (*MonitoredIntent, error) {
	var intent MonitoredIntent
	if err := r.db.GetGormDb().Where("operation_id = ?", operationId).First(&intent).Error; err != nil {
		return nil, err
	}
	return &intent, nil
}

func (r *intentRepo) FindWatchingByResource(resourceKey string) ([]MonitoredIntent, error) {
	var intents []MonitoredIntent
	err := r.db.GetGormDb().
		Where("resource_key = ? AND status = ?", resourceKey, IntentWatching).
		Find(&intents).Error
	return intents, err
}

func (r *intentRepo) MarkTerminal(operationId uuid.UUID, status IntentStatus) (*MonitoredIntent, error) {
	if !status.IsTerminal() {
		return nil, errors.New("MarkTerminal: status must be terminal")
	}
	var intent MonitoredIntent
	err := r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("operation_id = ?", operationId).First(&intent).Error; err != nil {
			return err
		}
		if intent.Status.IsTerminal() {
			return nil
		}
		intent.Status = status
		return tx.Save(&intent).Error
	})
	if err != nil {
		return nil, err
	}
	return &intent, nil
}

func (r *intentRepo) FindExpired(now time.Time, limit int) ([]MonitoredIntent, error) {
	var intents []MonitoredIntent
	err := r.db.GetGormDb().
		Where("deadline < ? AND status = ?", now, IntentWatching).
		Limit(limit).
		Find(&intents).Error
	return intents, err
}
