package db

import (
	"strings"

	"github.com/ukama/ukama/services/common/sql"
	"gorm.io/gorm/clause"
)

type SystemRepo interface {
	Add(sys *System) error
	Update(sys *System) error
	Delete(sys string) error
	GetByName(sys string) (*System, error)
}

type systemRepo struct {
	Db sql.Db
}

func NewSystemRepo(db sql.Db) *systemRepo {
	return &systemRepo{
		Db: db,
	}
}

func (s *systemRepo) Add(sys *System) error {
	d := s.Db.GetGormDb().Create(sys)
	return d.Error
}

func (s *systemRepo) Update(sys *System) error {
	d := s.Db.GetGormDb().Where(&System{Name: sys.Name}).Updates(sys)
	return d.Error
}

func (s *systemRepo) Delete(sys string) error {
	var system System
	result := s.Db.GetGormDb().Preload(clause.Associations).Unscoped().Delete(&system, "name = ?", strings.ToLower(sys))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *systemRepo) GetByName(sys string) (*System, error) {
	system := &System{}
	result := s.Db.GetGormDb().Preload(clause.Associations).First(system, "name = ?", strings.ToLower(sys))
	if result.Error != nil {
		return nil, result.Error
	}
	return system, nil
}
