// This is an example of a repositoryasrRepo
package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

// declare interface so that we can mock it
type ProfileRepo interface {
	Add(record *Profile) error
	GetByImsi(imsi string) (*Profile, error)
	GetByIccid(iccid string) (*Profile, error)
	UpdatePackage(imsi string, pkg PackageDetails) error
	Delete(imsi string) error
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
		AllowedTimeOfService: p.AllowedTimeOfService,
		TotalDataBytes:       p.TotalDataBytes,
		ConsumedDataBytes:    0,
	}
	d := r.db.GetGormDb().Where("imsi=?", imsiToUpdate).Updates(rec)
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

func (r *profileRepo) Delete(imsi string) error {

	res := r.db.GetGormDb().Where(&Profile{Imsi: imsi}).Delete(&Profile{})
	return res.Error
}
