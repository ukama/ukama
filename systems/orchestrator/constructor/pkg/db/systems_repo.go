package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm/clause"
)

type SystemsRepo interface {
	Add(sys *Systems) error
	AddDeployment(sys *Systems, dep *Deployments) error
	GetByName(name string) (*Systems, error)
	GetAll() ([]Systems, error)
	DeleteSystem(name string) error
	GetCount(name string) int64
}

type systemsRepo struct {
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
func NewSystemsRepo(db sql.Db) *systemsRepo {
	return &systemsRepo{
		Db: db,
	}
}

func (r *systemsRepo) Add(sys *Systems) error {

	d := r.Db.GetGormDb().Omit("Deployments").Create(sys)
	return d.Error
}

func (r *systemsRepo) AddDeployment(sys *Systems, dep *Deployments) error {

	err := r.Db.GetGormDb().Association("Deployments").Append(dep)
	return err
}

func (r *systemsRepo) GetByName(name string) (*Systems, error) {
	sys := &Systems{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Where(&Systems{Name: name}).Last(sys)
	return sys, d.Error
}

func (r *systemsRepo) GetById(id uint8) (*Systems, error) {
	sys := &Systems{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Where(&Systems{SystemId: id}).Last(sys)
	return sys, d.Error
}

func (r *systemsRepo) GetAll() ([]Systems, error) {
	sys := []Systems{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Find(&sys)
	return sys, d.Error
}

func (r *systemsRepo) DeleteSystem(name string) error {
	c := r.GetCount(name)
	if c != 0 {
		return fmt.Errorf("system is still have %d deployments", c)
	}

	d := r.Db.GetGormDb().Delete(&Systems{Name: name})
	return d.Error

}

func (r *systemsRepo) GetCount(name string) int64 {
	c := r.Db.GetGormDb().Association("Deployments").Count()
	return c
}
