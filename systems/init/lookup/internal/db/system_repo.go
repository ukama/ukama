package db

import (
	"strings"

	"github.com/ukama/ukama/services/common/sql"
	"gorm.io/gorm/clause"
)

type SystemRepo interface {
	AddOrUpdate(sys *System) error
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

func (s *systemRepo) AddOrUpdate(sys *System) error {
	d := s.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "lower(name::text)", Raw: true}},
		DoUpdates: clause.AssignmentColumns([]string{"uuid", "certificate", "ip", "port", "org_id", "deleted_at"}),
	}).Create(sys)
	return d.Error
}

func (s *systemRepo) Delete(sys string) error {
	var system System
	result := s.Db.GetGormDb().Preload(clause.Associations).Delete(&system, "name = ?", strings.ToLower(sys))
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
