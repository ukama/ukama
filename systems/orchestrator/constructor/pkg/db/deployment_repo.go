package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm/clause"
)

type DeploymentsRepo interface {
	Create(org uuid.UUID, name string) error
	Get(org uuid.UUID, name string) (*Deployments, error)
	Delete(org uuid.UUID, name string) error
}

type deploymentsRepo struct {
	Db sql.Db
}

func NewDeploymentsRepo(db sql.Db) *deploymentsRepo {
	return &deploymentsRepo{
		Db: db,
	}
}

func (d *deploymentsRepo) Create(Deployment *Deployments) error {
	r := d.Db.GetGormDb().Create(Deployment)
	return r.Error
}

func (d *deploymentsRepo) Get(org uuid.UUID, name string) (*Deployments, error) {
	var Deployment Deployments
	result := d.Db.GetGormDb().Where("name = ?", name).Where("org = ?", org).First(&Deployment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Deployment, nil
}

func (d *deploymentsRepo) Delete(org uuid.UUID, name string) error {
	var Deployment Deployments
	result := d.Db.GetGormDb().Where("name = ?", name).Where("org = ?", org).Delete(&Deployment)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	}

	return fmt.Errorf("deployment %s missing for org %s", name, org.String())
}
