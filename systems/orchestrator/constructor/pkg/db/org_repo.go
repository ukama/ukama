package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm/clause"
)

type OrgsRepo interface {
	Add(org *Orgs) error
	Get(id uuid.UUID) (*Orgs, error)
}

type orgsRepo struct {
	Db sql.Db
}

func NewOrgsRepo(db sql.Db) *orgsRepo {
	return &orgsRepo{
		Db: db,
	}
}

func (r *orgsRepo) Add(org *Orgs) error {

	d := r.Db.GetGormDb().Create(org)
	return d.Error
}

func (r *orgsRepo) Get(id uuid.UUID) (*Orgs, error) {
	org := &Orgs{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Where(&Orgs{OrgId: id}).First(org)
	return org, d.Error
}

func (r *orgsRepo) Delete(id uuid.UUID) error {
	org := &Orgs{
		OrgId: id,
	}
	d := r.Db.GetGormDb().Preload(clause.Associations, "Deployments").Delete(org)
	return d.Error
}
