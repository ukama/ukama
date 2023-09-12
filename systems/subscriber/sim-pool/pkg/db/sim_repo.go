package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type SimRepo interface {
	Get(isPhysicalSim bool, simType SimType) (*Sim, error)
	GetByIccid(iccid string) (*Sim, error)
	GetSimsByType(simType SimType) ([]Sim, error)
	Add(sims []Sim) error
	Delete(id []uint64) error
}

type simRepo struct {
	Db sql.Db
}

func NewSimRepo(db sql.Db) *simRepo {
	return &simRepo{
		Db: db,
	}
}

func (s *simRepo) Get(isPhysicalSim bool, simType SimType) (*Sim, error) {
	var sim Sim
	result := s.Db.GetGormDb().Where("is_allocated = ?", false).Where("is_physical = ?", isPhysicalSim).Where("sim_type = ?", simType).First(&sim)

	if result.Error != nil {
		return nil, result.Error
	}

	return &sim, nil
}

func (s *simRepo) GetByIccid(iccid string) (*Sim, error) {
	var sim Sim
	result := s.Db.GetGormDb().Where("is_allocated = ? AND iccid = ?", false, iccid).First(&sim)

	if result.Error != nil {
		return nil, result.Error
	}

	return &sim, nil
}

func (s *simRepo) GetSimsByType(simType SimType) ([]Sim, error) {
	var sim []Sim
	result := s.Db.GetGormDb().Where("sim_type = ?", simType).Find(&sim)

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
