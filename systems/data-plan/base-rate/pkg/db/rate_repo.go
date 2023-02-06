package db

import (
	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/sql"
)

type BaseRateRepo interface {
	GetBaseRate(uuid uuid.UUID) (*Rate, error)
	GetBaseRates(country, network, effectiveAt, simType string) ([]Rate, error)
	UploadBaseRates(rateList []Rate) error
}

type baseRateRepo struct {
	Db sql.Db
}

func NewBaseRateRepo(db sql.Db) *baseRateRepo {
	return &baseRateRepo{
		Db: db,
	}
}

func (u *baseRateRepo) GetBaseRate(uuid uuid.UUID) (*Rate, error) {
	rate := &Rate{}
	result := u.Db.GetGormDb().First(rate, "Uuid=?", uuid)
	if result.Error != nil {
		return nil, result.Error
	}
	return rate, nil
}

func (b *baseRateRepo) GetBaseRates(country, network, effectiveAt, simType string) ([]Rate, error) {
	var rates []Rate
	result := b.Db.GetGormDb().Where(&Rate{Country: country, Network: network, Sim_type: simType, Effective_at: effectiveAt}).Find(&rates)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) UploadBaseRates(rateList []Rate) error {
	e := b.Db.GetGormDb().Create(&rateList)
	if e != nil {
		return e.Error
	}

	return nil
}
