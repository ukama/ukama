// This is an example of a repositoryasrRepo
package db

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const TaiNotUpdatedErr = "more recent tai for imsi exist"

// declare interface so that we can mock it
type AsrRecordRepo interface {
	Add(record *Asr) error
	Get(id int) (*Asr, error)
	GetByImsi(imsi string) (*Asr, error)
	GetByIccid(iccid string) (*Asr, error)
	Update(imsi string, record *Asr) error
	UpdatePackage(imsi string, packageId uuid.UUID, policy *Policy) error
	DeleteByIccid(iccid string, nestedFunc ...func(*gorm.DB) error) error
	Delete(imsi string, nestedFunc ...func(*gorm.DB) error) error
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

func (r *asrRecordRepo) Delete(imsi string, nestedFuncs ...func(*gorm.DB) error) error {
	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		asrRec := &Asr{}
		err := tx.Model(&Asr{}).Where("imsi=?", imsi).Find(&asrRec).Error
		if err != nil {
			return errors.Wrap(err, "unable to find record for subscriber "+imsi)
		}
		log.Debugf("Deleting ASR record %+v", asrRec)

		err = tx.Where("asr_id=?", asrRec.ID).Delete(&Policy{}).Error
		if err != nil {
			return errors.Wrap(err, "error deleting policy for subscriber "+imsi)
		}

		err = tx.Where("asr_id=?", asrRec.ID).Delete(&Asr{}).Error
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

func (r *asrRecordRepo) DeleteByIccid(iccid string, nestedFunc ...func(*gorm.DB) error) error {
	return r.db.ExecuteInTransaction2(func(tx *gorm.DB) *gorm.DB {
		return tx.Where(&Asr{Iccid: iccid}).Delete(&Asr{})
	}, nestedFunc...)

}

// ReplaceTai removes all TAI record for IMSI and adds new ones
func (r *asrRecordRepo) UpdateTai(imsi string, tai Tai) error {
	var imsiM Asr
	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Asr{}).Where("imsi=?", imsi).First(&imsiM).Error
		if err != nil {
			return errors.Wrap(err, "error getting imsi")
		}

		var count int64
		err = tx.Model(&tai).Where("asr_id = ? and device_updated_at > ?", imsiM.ID, tai.DeviceUpdatedAt).Count(&count).Error
		if err != nil {
			return errors.Wrap(err, "error getting tai count")
		}
		if count > 0 {
			return fmt.Errorf(TaiNotUpdatedErr)
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
