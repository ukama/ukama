// This is an example of a repository
//
package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/errors"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const TaiNotUpdatedErr = "more recent tai for imsi exist"

// declare interface so that we can mock it
type ImsiRepo interface {
	Add(orgName string, imsi *Imsi) error
	Get(id int) (*Imsi, error)
	GetByImsi(imsi string) (*Imsi, error)
	GetImsiByUserUuid(userUuid uuid.UUID) ([]*Imsi, error)
	Update(imsi string, subscriber *Imsi) error
	Delete(imsi string, nestedFunc ...func(*gorm.DB) error) error
	DeleteByUserId(user uuid.UUID, nestedFunc ...func(*gorm.DB) error) error
	UpdateTai(imis string, tai Tai) error
}

type imsiRepo struct {
	db sql.Db
}

func NewImsiRepo(db sql.Db) *imsiRepo {
	return &imsiRepo{
		db: db,
	}
}

func (r *imsiRepo) Add(orgName string, imsi *Imsi) error {
	org, err := makeUserOrgExist(r.db.GetGormDb(), orgName)
	if err != nil {
		return err
	}
	imsi.Org = org
	d := r.db.GetGormDb().Create(imsi)
	return d.Error
}

func (r *imsiRepo) Update(imsiToUpdate string, imsi *Imsi) error {
	d := r.db.GetGormDb().Where("imsi=?", imsiToUpdate).Updates(imsi)
	return d.Error
}

func (r *imsiRepo) Get(id int) (*Imsi, error) {
	var hss Imsi
	result := r.db.GetGormDb().Preload(clause.Associations).First(&hss, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &hss, nil
}

func (r *imsiRepo) GetByImsi(imsi string) (*Imsi, error) {
	var hss Imsi
	result := r.db.GetGormDb().Preload(clause.Associations).Where("imsi=?", imsi).First(&hss)
	if result.Error != nil {
		return nil, result.Error
	}

	return &hss, nil
}

func (r *imsiRepo) GetImsiByUserUuid(userUuid uuid.UUID) ([]*Imsi, error) {
	var imsis []*Imsi
	result := r.db.GetGormDb().Preload(clause.Associations).Where("user_uuid=?", userUuid).Find(&imsis)
	if result.Error != nil {
		return nil, result.Error
	}

	return imsis, nil
}

func (r *imsiRepo) Delete(imsi string, nestedFunc ...func(*gorm.DB) error) error {
	return r.db.ExecuteInTransaction2(func(tx *gorm.DB) *gorm.DB {
		return tx.Where(&Imsi{Imsi: imsi}).Delete(&Imsi{})
	}, nestedFunc...)
}

func (r *imsiRepo) DeleteByUserId(user uuid.UUID, nestedFunc ...func(*gorm.DB) error) error {
	return r.db.ExecuteInTransaction2(func(tx *gorm.DB) *gorm.DB {
		return tx.Where(&Imsi{UserUuid: user}).Delete(&Imsi{})
	}, nestedFunc...)

}

// ReplaceTai removes all TAI record for IMSI and adds new ones
func (r *imsiRepo) UpdateTai(imsi string, tai Tai) error {
	var imsiM Imsi
	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Imsi{}).Where("imsi=?", imsi).First(&imsiM).Error
		if err != nil {
			return errors.Wrap(err, "error getting imsi")
		}

		var count int64
		err = tx.Model(&tai).Where("imsi_id = ? and device_updated_at >= ?", imsiM.ID, tai.DeviceUpdatedAt).Count(&count).Error
		if err != nil {
			return errors.Wrap(err, "error getting tai count")
		}
		if count > 0 {
			return fmt.Errorf(TaiNotUpdatedErr)
		}

		err = tx.Where("imsi_id=?", imsiM.ID).Delete(&Tai{}).Error
		if err != nil {
			return errors.Wrap(err, "error deleting tai")
		}

		tai.ImsiID = imsiM.ID
		err = tx.Create(&tai).Error
		if err != nil {
			return errors.Wrap(err, "error adding tai")
		}
		return nil
	})
}
