package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type SimPoolRepo interface {
	GetStats(SimType string) ([]SimPool, error)
	Add(simPools []SimPool) error
	Delete(Id uint64) error
}

type simPoolRepo struct {
	Db sql.Db
}

func NeSimPoolRepo(db sql.Db) *simPoolRepo {
	return &simPoolRepo{
		Db: db,
	}
}

func (s *simPoolRepo) GetStats(SimType string) ([]SimPool, error) {
	var simPool []SimPool
	result := s.Db.GetGormDb().Where(&SimPool{Sim_type: SimType}).Find(&simPool)

	if result.Error != nil {
		return nil, result.Error
	}

	return simPool, nil
}

func (s *simPoolRepo) Add(simPools []SimPool) error {
	e := s.Db.GetGormDb().Create(&simPools)
	if e != nil {
		return e.Error
	}

	return nil
}

func (s *simPoolRepo) Delete(Id uint64) error {
	result := s.Db.GetGormDb().Where("id = ?", Id).Delete(&SimPool{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
