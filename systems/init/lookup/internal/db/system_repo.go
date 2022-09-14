package db

import (
	"strings"

	"github.com/ukama/ukama/services/common/sql"
	"gorm.io/gorm/clause"
)

type SystemRepo interface {
	AddOrUpdate(sys *System) error
	Get(sys string) (*System, error)
	Delete(sys string, orgId uint) error
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
		DoUpdates: clause.AssignmentColumns([]string{"id", "certificate", "ip", "port", "org_id"}),
	}).Create(sys)
	return d.Error
}

func (s *systemRepo) Get(sys string) (*System, error) {
	var system System
	result := s.Db.GetGormDb().Preload(clause.Associations).First(&system, "name = ?", strings.ToLower(sys))
	if result.Error != nil {
		return nil, result.Error
	}
	return &system, nil
}

func (s *systemRepo) Delete(sys string, orgId uint) error {
	var system System
	result := s.Db.GetGormDb().Preload(clause.Associations).First(&system, "name = ?", strings.ToLower(sys))
	if result.Error != nil {
		return nil, result.Error
	}
	return &system, nil
}
