package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm/clause"
)

type OrgRepo interface {
	Add(org *Org) error
	Get(id uuid.UUID) (*Org, error)
}

type orgsRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) *orgsRepo {
	return &orgsRepo{
		Db: db,
	}
}

func (r *orgsRepo) Add(org *Org) error {

	d := r.Db.GetGormDb().Create(org)
	return d.Error
}

func (r *orgsRepo) Get(id uuid.UUID) (*Org, error) {
	org := &Org{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Where(&Org{OrgId: id}).First(org)
	return org, d.Error
}

func (r *orgsRepo) Delete(id uuid.UUID) error {
	org := &Org{
		OrgId: id,
	}
	d := r.Db.GetGormDb().Preload(clause.Associations, "Deployments").Delete(org)
	return d.Error
}
