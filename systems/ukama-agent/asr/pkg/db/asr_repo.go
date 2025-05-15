/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
)

const TaiNotUpdatedErr = "more recent tai for imsi exist"

// declare interface so that we can mock it
type AsrRecordRepo interface {
	Add(record *Asr) error
	Get(id int) (*Asr, error)
	List() ([]Asr, error)
	GetByImsi(imsi string) (*Asr, error)
	GetByIccid(iccid string) (*Asr, error)
	Update(imsi string, record *Asr) error
	UpdatePackage(imsi string, packageId uuid.UUID, policy *Policy) error
	DeleteByIccid(iccid string, reason StatusReason, nestedFunc ...func(*gorm.DB) error) error
	Delete(imsi string, reason StatusReason, nestedFunc ...func(*gorm.DB) error) error
	UpdateTai(imis string, tai Tai) error
}

type asrRecordRepo struct {
	db sql.Db
}

func NewAsrRecordRepo(db sql.Db) *asrRecordRepo {
	return &asrRecordRepo{
		db: db,
	}
}

func (r *asrRecordRepo) Add(rec *Asr) error {
	d := r.db.GetGormDb().Create(rec)
	return d.Error
}

func (r *asrRecordRepo) Update(imsiToUpdate string, rec *Asr) error {
	d := r.db.GetGormDb().Where("imsi=?", imsiToUpdate).Updates(rec)
	return d.Error
}

func (r *asrRecordRepo) UpdatePackage(imsiToUpdate string, packageId uuid.UUID, policy *Policy) error {
	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		asrRec := &Asr{}
		err := tx.Model(&Asr{}).Where("imsi=?", imsiToUpdate).Find(&asrRec).Error
		if err != nil {
			return errors.Wrap(err, "unable to find record for subscriber "+imsiToUpdate)
		}
		log.Debugf("Updating ASR record %+v", asrRec)

		err = tx.Where("asr_id=?", asrRec.ID).Delete(&Policy{}).Error
		if err != nil {
			return errors.Wrap(err, "error deleting policy for subscriber "+imsiToUpdate)
		}

		policy.AsrID = asrRec.ID
		err = tx.Create(&policy).Error
		if err != nil {
			return errors.Wrap(err, "error adding policy")
		}
		asrRec.PackageId = packageId
		asrRec.LastStatusChangeAt = time.Now()
		asrRec.LastStatusChangeReasons = PACKAGE_UPDATE

		err = tx.Model(&Asr{}).Where("imsi=?", imsiToUpdate).Updates(asrRec).Error
		if err != nil {
			return errors.Wrap(err, "error updating reason")
		}

		return nil
	})

}

func (r *asrRecordRepo) Get(id int) (*Asr, error) {
	var hss Asr
	result := r.db.GetGormDb().Preload(clause.Associations).First(&hss, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &hss, nil
}

func (r *asrRecordRepo) List() ([]Asr, error) {
	var p []Asr
	result := r.db.GetGormDb().Preload(clause.Associations).Find(&p)
	if result.Error != nil {
		return nil, result.Error
	}

	return p, nil
}

func (r *asrRecordRepo) GetByImsi(imsi string) (*Asr, error) {
	var asr Asr
	result := r.db.GetGormDb().Preload(clause.Associations).Where("imsi=?", imsi).First(&asr)
	if result.Error != nil {
		return nil, result.Error
	}

	return &asr, nil
}

func (r *asrRecordRepo) GetByIccid(iccid string) (*Asr, error) {
	var asr Asr
	result := r.db.GetGormDb().Preload(clause.Associations).Where("iccid=?", iccid).First(&asr)
	if result.Error != nil {
		return nil, result.Error
	}

	return &asr, nil
}

func (r *asrRecordRepo) Delete(imsi string, reason StatusReason, nestedFuncs ...func(*gorm.DB) error) error {
	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		asrRec := &Asr{}
		err := tx.Model(&Asr{}).Where("imsi=?", imsi).Find(&asrRec).Error
		if err != nil {
			return errors.Wrap(err, "unable to find record for subscriber "+imsi)
		}
		log.Debugf("Deleting ASR record %+v", asrRec)

		asrRec.LastStatusChangeAt = time.Now()
		asrRec.LastStatusChangeReasons = reason

		err = tx.Model(&Asr{}).Where("imsi=?", imsi).Updates(&asrRec).Error
		if err != nil {
			return errors.Wrap(err, "unable to update record for subscriber "+imsi)
		}

		err = tx.Where("asr_id=?", asrRec.ID).Delete(&Policy{}).Error
		if err != nil {
			return errors.Wrap(err, "error deleting policy for subscriber "+imsi)
		}

		err = tx.Where("id=?", asrRec.ID).Delete(&Asr{}).Error
		if err != nil {
			return errors.Wrap(err, "error deleting ASR for subscriber "+imsi)
		}

		if len(nestedFuncs) > 0 {
			for _, n := range nestedFuncs {
				if n != nil {
					nestErr := n(tx)
					if nestErr != nil {
						return nestErr
					}
				}
			}
		}
		return nil
	})
}

func (r *asrRecordRepo) DeleteByIccid(iccid string, reason StatusReason, nestedFunc ...func(*gorm.DB) error) error {
	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		s := &Asr{
			LastStatusChangeAt:      time.Now(),
			LastStatusChangeReasons: reason,
		}

		err := tx.Model(&Asr{}).Where("iccid=?", iccid).Updates(s).Error
		if err != nil {
			return errors.Wrap(err, "error updating reason")
		}

		err = tx.Where(&Asr{Iccid: iccid}).Delete(&Asr{}).Error
		if err != nil {
			return errors.Wrap(err, "error deleting subscriber")
		}

		if len(nestedFunc) > 0 {
			for _, n := range nestedFunc {
				if n != nil {
					nestErr := n(tx)
					if nestErr != nil {
						return nestErr
					}
				}
			}
		}

		return nil
	})

}

// ReplaceTai removes all TAI record for IMSI and adds new onesactionName
func (r *asrRecordRepo) UpdateTai(imsi string, tai Tai) error {
	var imsiM Asr
	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Asr{}).Where("imsi=?", imsi).First(&imsiM).Error
		if err != nil {
			return errors.Wrap(err, "error getting ASR record for given imsi")
		}

		var count int64
		err = tx.Model(&tai).Where("asr_id = ? and device_updated_at > ?", imsiM.ID, tai.DeviceUpdatedAt).Count(&count).Error
		if err != nil {
			return errors.Wrap(err, "error getting tai count")
		}
		if count > 0 {
			return errors.New(TaiNotUpdatedErr)
		}

		err = tx.Where("asr_id=?", imsiM.ID).Delete(&Tai{}).Error
		if err != nil {
			return errors.Wrap(err, "error deleting tai")
		}

		tai.AsrID = imsiM.ID
		err = tx.Create(&tai).Error
		if err != nil {
			return errors.Wrap(err, "error adding tai")
		}
		return nil
	})
}
