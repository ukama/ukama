package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
)

type OrgRepo interface {
	Add(org *Org) error
	Update(org *Org) error
	GetByName(name string) (*Org, error)
	GetById(id uuid.UUID) (*Org, error)
}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) *orgRepo {
	return &orgRepo{
		Db: db,
	}
}

func (r *orgRepo) Add(org *Org) error {

	d := r.Db.GetGormDb().Create(org)
	return d.Error
}

func (r *orgRepo) Update(org *Org) error {
	d := r.Db.GetGormDb().Where(&Org{Name: org.Name}).Updates(org)
	return d.Error
}

func (r *orgRepo) GetByName(name string) (*Org, error) {
	org := &Org{}
	d := r.Db.GetGormDb().Where(&Org{Name: name}).First(&org)
	return org, d.Error
}

func (r *orgRepo) GetById(id uuid.UUID) (*Org, error) {
	org := &Org{}
	d := r.Db.GetGormDb().Where(&Org{OrgId: id}).First(&org)
	return org, d.Error
}
