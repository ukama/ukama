package db

import (
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/common/sql"
)

// declare interface so that we can mock it

type BaseRateRepo interface {
	GetBaseRate(Id uint64) (*Rate, error)
	GetBaseRates(country, network,effectiveAt, simType string) ([]Rate, error)
	UploadBaseRates(query string) error
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

func (b *baseRateRepo) GetBaseRates(country, network,effectiveAt, simType string) ([]Rate, error) {
	var rates []Rate
	date, error := time.Parse("2006-01-02", effectiveAt)
	if error != nil {
		fmt.Println(error)
	}
	result:= b.Db.GetGormDb().Where(&Rate{Country: country, Network: network,Sim_type:simType,Effective_at:date}).Find(&rates)

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
