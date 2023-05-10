package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SystemsRepo interface {
	Add(sys *Systems) error
	Get(name string) (*Systems, error)
	GetById(id uint) (*Systems, error)
}

type systemsRepo struct {
	Db sql.Db
}

func NewSystemsRepo(db sql.Db) *systemsRepo {
	return &systemsRepo{
		Db: db,
	}
}

func (r *systemsRepo) Add(sys *Systems) error {

	d := r.Db.GetGormDb().Create(sys)
	return d.Error
}

func (r *systemsRepo) Get(name string) (*Systems, error) {
	sys := &Systems{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Where(&Systems{Name: name}).Last(sys)
	return sys, d.Error
}

func (r *systemsRepo) GetById(id uint) (*Systems, error) {
	sys := &Systems{}
	d := r.Db.GetGormDb().Preload(clause.Associations).Where(&Systems{Model: gorm.Model{ID: id}}).Last(sys)
	return sys, d.Error
}
