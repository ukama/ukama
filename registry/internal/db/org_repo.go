package db

import (
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm/clause"
)

type OrgRepo interface {
	Add(org *Org) error
	Get(id int) (*Org, error)
}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) OrgRepo {
	return &orgRepo{
		Db: db,
	}
}

func (r *orgRepo) Add(org *Org) error {
	d := r.Db.GetGormDb().Create(org)
	return d.Error
}

func (r *orgRepo) Get(id int) (*Org, error) {
	var org Org
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&org, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &org, nil
}
