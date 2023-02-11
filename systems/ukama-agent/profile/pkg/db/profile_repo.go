// This is an example of a repositoryasrRepo
package db

import (
	"github.com/pkg/errors"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

// declare interface so that we can mock it
type ProfileRepo interface {
	Add(record *Profile) error
	GetByImsi(imsi string) (*Profile, error)
	GetByIccid(iccid string) (*Profile, error)
	UpdatePackage(imsi string, pkg PackageDetails) error
	UpdateUsage(imsi string, bytes uint64) error
	Delete(imsi string, reason StatusReason) error
	List() ([]Profile, error)
}

type profileRepo struct {
	db sql.Db
}

func NewProfileRepo(db sql.Db) *profileRepo {
	return &profileRepo{
		db: db,
	}
}

func (r *profileRepo) Add(rec *Profile) error {
	d := r.db.GetGormDb().Create(rec)
	return d.Error
}

func (r *profileRepo) UpdatePackage(imsiToUpdate string, p PackageDetails) error {
	rec := &Profile{PackageId: p.PackageId,
		AllowedTimeOfService:    p.AllowedTimeOfService,
		TotalDataBytes:          p.TotalDataBytes,
		ConsumedDataBytes:       0,
		LastStatusChangeReasons: PACKAGE_UPDATE,
	}
	d := r.db.GetGormDb().Where("imsi=?", imsiToUpdate).Updates(rec)
	return d.Error
}

func (r *profileRepo) UpdateUsage(imsi string, bytes uint64) error {
	rec := &Profile{
		ConsumedDataBytes: bytes,
	}
	d := r.db.GetGormDb().Where("imsi=?", imsi).Updates(rec)
	return d.Error
}

func (r *profileRepo) GetByImsi(imsi string) (*Profile, error) {
	var asr Profile
	result := r.db.GetGormDb().Where("imsi=?", imsi).First(&asr)
	if result.Error != nil {
		return nil, result.Error
	}

	return &asr, nil
}

func (r *profileRepo) GetByIccid(iccid string) (*Profile, error) {
	var pro Profile
	result := r.db.GetGormDb().Where("iccid=?", iccid).First(&pro)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pro, nil
}

func (r *profileRepo) Delete(imsi string, reason StatusReason) error {

	p := &Profile{
		LastStatusChangeReasons: reason,
	}

	return r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Profile{}).Where("imsi=?", imsi).Updates(p).Error
		if err != nil {
			return errors.Wrap(err, "error getting imsi")
		}

		err = tx.Where(&Profile{Imsi: imsi}).Delete(&Profile{}).Error
		if err != nil {
			return errors.Wrap(err, "error removing imsi")
		}
		return nil
	})
}

func (r *profileRepo) List() ([]Profile, error) {
	var p []Profile
	result := r.db.GetGormDb().Find(&p)
	if result.Error != nil {
		return nil, result.Error
	}

	return p, nil
}
