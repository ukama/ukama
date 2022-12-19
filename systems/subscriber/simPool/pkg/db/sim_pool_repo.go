package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

type SimPoolRepo interface {
	GetStats(Id uint64, SimType string) ([]SimPool, error)
	Add(simPools []SimPool) ([]SimPool, error)
	Delete(Id uint64) error
}

type simPoolRepo struct {
	Db sql.Db
}

func GetSimPoolRepo(db sql.Db) *simPoolRepo {
	return &simPoolRepo{
		Db: db,
	}
}

func (u *simPoolRepo) GetStats(Id uint64, SimType string) (*SimPool, error) {
	return nil, nil
}

func (b *simPoolRepo) Add(simPools []SimPool) ([]SimPool, error) {
	return nil, nil
}

func (b *simPoolRepo) Delete(Id uint64) error {
	return nil
}
