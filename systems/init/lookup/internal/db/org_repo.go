package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm/clause"
)

type OrgRepo interface {
	Upsert(org *Org) error
	GetByName(name string) (*Org, error)
}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) *orgRepo {
	return &orgRepo{
		Db: db,
	}
}

func (r *orgRepo) Upsert(org *Org) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"certificate", "ip"}),
	}).Create(org)

	return d.Error
}

func (r *orgRepo) GetByName(name string) (*Org, error) {
	org := &Org{}
	d := r.Db.GetGormDb().Where(&Org{Name: name}).First(&org)
	return org, d.Error
}
