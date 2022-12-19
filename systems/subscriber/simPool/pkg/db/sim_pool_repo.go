package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

type SimPoolRepo interface {
	GetStats(Id uint64, SimType string) ([]SimPool, error)
	Add(networkId string, orgId uint64, iccid string, msisdn string, isAllocated bool, simType string) ([]SimPool, error)
	Delete(Id uint64) error
	Upload(fileUrl string, simType string, orgId uint64) ([]SimPool, error)
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

func (b *simPoolRepo) Add(networkId string, orgId uint64, iccid string, msisdn string, isAllocated bool, simType string) ([]SimPool, error) {
	return nil, nil
}

func (b *simPoolRepo) Delete(Id uint64) error {
	return nil
}

func (b *simPoolRepo) Upload(fileUrl string, simType string, orgId uint64) ([]SimPool, error) {
	return nil, nil
}
