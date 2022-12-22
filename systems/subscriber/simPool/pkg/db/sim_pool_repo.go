package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type SimRepo interface {
	GetStats(SimType string) ([]Sim, error)
	Add(sims []Sim) error
	Delete(Id []uint64) error
}

type simRepo struct {
	Db sql.Db
}

func NewSimRepo(db sql.Db) *simRepo {
	return &simRepo{
		Db: db,
	}
}

func (s *simRepo) GetStats(SimType string) ([]Sim, error) {
	var sim []Sim
	result := s.Db.GetGormDb().Where(&Sim{Sim_type: SimType}).Find(&sim)

	if result.Error != nil {
		return nil, result.Error
	}

	return sim, nil
}

func (s *simRepo) Add(sims []Sim) error {
	e := s.Db.GetGormDb().Create(&sims)
	if e != nil {
		return e.Error
	}

	return nil
}

func (s *simRepo) Delete(Id []uint64) error {
	result := s.Db.GetGormDb().Where("id IN (?)", Id).Delete(&Sim{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
