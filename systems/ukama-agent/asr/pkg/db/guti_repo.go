/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	dbsql "database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

const GutiNotUpdatedErr = "more recent guti for imsi exist"

type GutiRepo interface {
	Update(guti *Guti) error
	GetImsi(plmnId string, mmegi, mmec, mTmsi uint32) (string, error)
}

type gutiRepo struct {
	db sql.Db
}

func NewGutiRepo(db sql.Db) *gutiRepo {
	return &gutiRepo{db: db}
}

// Only one guti per IMSI
func (g gutiRepo) Update(guti *Guti) error {
	var count int64

	// Two concurrent Update() calls for the same imsi can both
	// pass the "no newer guti" count check before either commits, breaking
	// the "only one guti per imsi" invariant this function is meant to
	// enforce. Serializable makes the DB abort the losing transaction
	// instead.
	err := g.db.GetGormDb().Transaction(
		func(tx *gorm.DB) error {
			err := tx.Model(&Guti{}).Where("imsi = ? and device_updated_at > ?", guti.Imsi, guti.DeviceUpdatedAt).Count(&count).Error
			if err != nil {
				return errors.Wrap(err, "failed get guti count")
			}
			if count > 0 {
				return errors.New(GutiNotUpdatedErr)
			}

			err = tx.Delete(&Guti{}, "imsi = ? and device_updated_at <= ?  ", guti.Imsi, guti.DeviceUpdatedAt).Error
			if err != nil {
				return errors.Wrap(err, "failed delete guti")
			}

			guti.CreatedAt = time.Now().UTC()
			return tx.Create(guti).Error
		}, &dbsql.TxOptions{Isolation: dbsql.LevelSerializable})
	return err
}

func (g gutiRepo) GetImsi(plmnId string, mmegi, mmec, mTmsi uint32) (string, error) {
	res := Guti{}

	err := g.db.GetGormDb().
		Where("plmn_id = ? AND mmegi = ? AND mmec = ? AND m_tmsi = ?", plmnId, mmegi, mmec, mTmsi).
		First(&res).Error

	return res.Imsi, err
}
