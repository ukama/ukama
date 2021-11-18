// This is an example of a repository
//
package db

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm/clause"
)

// declare interface so that we can mock it
type ImsiRepo interface {
	Add(orgName string, imsi *Imsi) error
	Get(id int) (*Imsi, error)
	GetByImsi(imsi string) (*Imsi, error)
	GetImsiByUserUuid(userUuid uuid.UUID) ([]*Imsi, error)
	Update(imsi string, subscriber *Imsi) error
	Delete(imsi string) error
	DeleteByUserId(user uuid.UUID) error
}

type imsiRepo struct {
	Db sql.Db
}

func NewImsiRepo(db sql.Db) *imsiRepo {
	return &imsiRepo{
		Db: db,
	}
}

func (r *imsiRepo) Add(orgName string, imsi *Imsi) error {
	org, err := makeUserOrgExist(r.Db.GetGormDb(), orgName)
	if err != nil {
		return err
	}
	imsi.Org = org
	d := r.Db.GetGormDb().Create(imsi)
	return d.Error
}

func (r *imsiRepo) Update(imsiToUpdate string, imsi *Imsi) error {
	d := r.Db.GetGormDb().Where("imsi=?", imsiToUpdate).Updates(imsi)
	return d.Error
}

func (r *imsiRepo) Get(id int) (*Imsi, error) {
	var hss Imsi
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&hss, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &hss, nil
}

func (r *imsiRepo) GetByImsi(imsi string) (*Imsi, error) {
	var hss Imsi
	result := r.Db.GetGormDb().Preload(clause.Associations).Where("imsi=?", imsi).First(&hss)
	if result.Error != nil {
		return nil, result.Error
	}

	return &hss, nil
}

func (r *imsiRepo) GetImsiByUserUuid(userUuid uuid.UUID) ([]*Imsi, error) {
	var imsis []*Imsi
	result := r.Db.GetGormDb().Where("user_uuid=?", userUuid).Find(&imsis)
	if result.Error != nil {
		return nil, result.Error
	}

	return imsis, nil
}

func (r *imsiRepo) Delete(imsi string) error {
	result := r.Db.GetGormDb().Where(&Imsi{Imsi: imsi}).Delete(&Imsi{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *imsiRepo) DeleteByUserId(user uuid.UUID) error {
	result := r.Db.GetGormDb().Where(&Imsi{UserUuid: user}).Delete(&Imsi{})
	if result.Error != nil {
		return result.Error
	}
	logrus.Infof("Deleted %d imsis", result.RowsAffected)
	return nil
}
