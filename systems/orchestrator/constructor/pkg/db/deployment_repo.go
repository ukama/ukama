package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm/clause"
)

type DeploymentsRepo interface {
	Add(d *Deployments) error
	GetByName(name string) (*Deployments, error)
	GetAll() ([]Deployments, error)
	GetById(id uuid.UUID) (*Deployments, error)
	Delete(id uuid.UUID) error
	GetCount(name string) (int64, error)
	GetHistoryByName(name string) (*Deployments, error)
}

type deploymentRepo struct {
	Db sql.Db
}

/*
Add systems without any deployment on startup

	Add depoyments per org for a system on request
	For update to system :
	1) Create a new system
	2) Create a new deploymenet one by one for every org and simultaneously remove it from last system deployment association
	3) Once all deployments are cretaed and removed from the last system delete the last system

	For deleting system :
	1) Deleete all teh deployments , remove from association.
	2) delete system
*/
func NewDeploymentsRepo(db sql.Db) *deploymentRepo {
	return &deploymentRepo{
		Db: db,
	}
}

func (r *deploymentRepo) Add(d *Deployments) error {

	res := r.Db.GetGormDb().Omit("Orgs").Create(d)
	return res.Error
}

// func (r *deploymentRepo) AddDeployment(d *Deployments, dep *Deployments) error {

// 	err := r.Db.GetGormDb().Association("Deployments").Append(dep)
// 	return err
// }

func (r *deploymentRepo) GetByName(name string) (*Deployments, error) {
	sys := &Deployments{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Where(&Deployments{Name: name}).Last(sys)
	return sys, d.Error
}

func (r *deploymentRepo) GetById(id uuid.UUID) (*Deployments, error) {
	sys := &Deployments{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Where(&Deployments{Code: id}).Last(sys)
	return sys, d.Error
}

func (r *deploymentRepo) GetAll() ([]Deployments, error) {
	sys := []Deployments{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Find(&sys)
	return sys, d.Error
}

func (r *deploymentRepo) Delete(id uuid.UUID) error {
	d := r.Db.GetGormDb().Delete(&Deployments{Code: id})
	return d.Error

}

func (r *deploymentRepo) GetCount(name string) (int64, error) {
	var count int64
	d := r.Db.GetGormDb().Find(&Deployments{}).Count(&count)
	if d.Error != nil {
		return 0, d.Error
	}
	return count, nil
	//c := r.Db.GetGormDb().Association("Deployments").Count()
	//return c
}

func (r *deploymentRepo) GetHistoryByName(name string) (*Deployments, error) {
	sys := &Deployments{}
	d := r.Db.GetGormDb().Unscoped().Preload(clause.Associations).Where(&Deployments{Name: name}).Last(sys)
	return sys, d.Error
}
