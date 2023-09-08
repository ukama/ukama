package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm/clause"
)

type ResourceRepo interface {
	Add(r *Resource) error
	Get(id uuid.UUID) (*Resource, error)
	Update(r *Resource) error
	Delete(id uuid.UUID) error
}

type resourceRepo struct {
	Db sql.Db
}

func NewResourceRepo(db sql.Db) *resourceRepo {
	return &resourceRepo{
		Db: db,
	}
}

func (r *resourceRepo) Add(n *Resource) error {
	d := r.Db.GetGormDb().Create(n)

	return d.Error
}

func (r *resourceRepo) Get(id uuid.UUID) (*Resource, error) {
	Resource := Resource{}

	result := r.Db.GetGormDb().Preload(clause.Associations).
		First(&Resource, "id = ?", id.String())

	if result.Error != nil {
		return nil, result.Error
	}

	return &Resource, nil
}

func (r *resourceRepo) Update(n *Resource) error {
	d := r.Db.GetGormDb().Updates(n)

	return d.Error
}

func (r *resourceRepo) Delete(id uuid.UUID) error {
	result := r.Db.GetGormDb().Where("id = ?", id.String()).Delete(&Resource{})

	return result.Error
}
