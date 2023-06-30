package db

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeploymentRepo interface {
	Add(d *Deployment) error
	GetByName(name string) (*Deployment, error)
	GetAll() ([]Deployment, error)
	GetById(id uuid.UUID) (*Deployment, error)
	Delete(id uuid.UUID) error
	GetCount(name string) (int64, error)
	GetHistoryByName(name string) (*Deployment, error)
	UpdateStatus(name string, status DeploymentStatus) error
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
func NewDeploymentRepo(db sql.Db) *deploymentRepo {
	return &deploymentRepo{
		Db: db,
	}
}

func (r *deploymentRepo) Add(d *Deployment) error {

	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		t := tx.Where("name= ?", d.Name).Delete(&Deployment{})
		if t.RowsAffected > 0 {
			log.Debugf("Marking old state with delete_at for %s", d.Name)
		}

		res := tx.Model(&Deployment{}).Omit("Orgs").Omit("Configs").Create(d)
		if res.Error != nil {
			log.Errorf("failed to create new deployment %+v", d)
			return res.Error
		}

		err := tx.Model(&Deployment{Code: d.Code}).Omit("Orgs.*").Association("Orgs").Append(d.Org)
		if err != nil {
			log.Errorf("failed to update org in deployment %+v", d)
			return err
		}

		err = tx.Model(&Deployment{Code: d.Code}).Omit("Configs.*").Association("Configs").Append(d.Config)
		if err != nil {
			log.Errorf("failed to update org in deployment %+v", d)
			return err
		}

		return nil
	})
	//res := r.Db.GetGormDb().Model(&Deployment{}).Omit("Orgs.*").Omit("Configs.*").Association("Orgs", "Configs").Create(d)
	return err
}

// func (r *deploymentRepo) AddDeployment(d *Deployment, dep *Deployment) error {

// 	err := r.Db.GetGormDb().Association("Deployment").Append(dep)
// 	return err
// }

func (r *deploymentRepo) GetByName(name string) (*Deployment, error) {
	sys := &Deployment{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Where(&Deployment{Name: name}).Last(sys)
	return sys, d.Error
}

func (r *deploymentRepo) GetById(id uuid.UUID) (*Deployment, error) {
	sys := &Deployment{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Where(&Deployment{Code: id}).Last(sys)
	return sys, d.Error
}

func (r *deploymentRepo) GetAll() ([]Deployment, error) {
	sys := []Deployment{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Find(&sys)
	return sys, d.Error
}

func (r *deploymentRepo) Delete(id uuid.UUID) error {
	d := r.Db.GetGormDb().Delete(&Deployment{Code: id})
	return d.Error

}

func (r *deploymentRepo) GetCount(name string) (int64, error) {
	var count int64
	d := r.Db.GetGormDb().Find(&Deployment{}).Count(&count)
	if d.Error != nil {
		return 0, d.Error
	}
	return count, nil
	//c := r.Db.GetGormDb().Association("Deployment").Count()
	//return c
}

func (r *deploymentRepo) GetHistoryByName(name string) (*Deployment, error) {
	sys := &Deployment{}
	d := r.Db.GetGormDb().Unscoped().Preload(clause.Associations).Where(&Deployment{Name: name}).Last(sys)
	return sys, d.Error
}

func (r *deploymentRepo) UpdateStatus(name string, status DeploymentStatus) error {
	d := r.Db.GetGormDb().Where(&Deployment{Name: name}).Updates(&Deployment{Status: status})
	return d.Error
}
