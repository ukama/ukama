package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

// declare interface so that we can mock it

type BaseRateRepo interface {
	GetBaseRate(Id uint64) (*Rate, error)
	GetBaseRates(country, network, simType string) (RateList, error)
	UploadBaseRates(query string) error
	GetAllBaseRates(effectiveAt string) (RateList, error)
}

type baseRateRepo struct {
	Db sql.Db
}


func NewBaseRateRepo(db sql.Db) *baseRateRepo {
	return &baseRateRepo{
		Db: db,
	}
}

func (u *baseRateRepo) GetBaseRate(rateId uint64) (*Rate, error) {
	var rate *Rate
	result := u.Db.GetGormDb().First(&rate, "Id=?", rateId)
	if result.Error != nil {
		return nil, result.Error
	}
	return rate, nil
}

func (b *baseRateRepo) GetBaseRates(country, network, simType string) (RateList, error) {
	var rates RateList
	result:= b.Db.GetGormDb().Where(&Rate{Country: country, Network: network,Sim_type:simType}).Find(&rates)

		if result.Error != nil {
			return nil, result.Error
		}

	return rates, nil
}

func (b *baseRateRepo) GetAllBaseRates(effectiveAt string) (RateList, error) {
	var rates RateList
	result := b.Db.GetGormDb().Where("effective_at", effectiveAt).Find(&rates)
	if result.Error != nil {
		return nil, result.Error
	}
	return rates, nil
}

func (b *baseRateRepo) UploadBaseRates(query string) error {
	e := b.Db.GetGormDb().Exec(query)
	if e != nil {
		return e.Error
	}

	return nil
}
