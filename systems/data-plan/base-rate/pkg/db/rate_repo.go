package db

import (
	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/common/sql"
)

type BaseRateRepo interface {
	GetBaseRate(uuid uuid.UUID) (*BaseRate, error)
	GetBaseRates(country, network, effectiveAt string, simType SimType) ([]BaseRate, error)
	UploadBaseRates(rateList []BaseRate) error
}

type baseRateRepo struct {
	Db sql.Db
}

func NewBaseRateRepo(db sql.Db) *baseRateRepo {
	return &baseRateRepo{
		Db: db,
	}
}

func (u *baseRateRepo) GetBaseRate(uuid uuid.UUID) (*BaseRate, error) {
	rate := &BaseRate{}
	result := u.Db.GetGormDb().First(rate, "Uuid=?", uuid)
	if result.Error != nil {
		return nil, result.Error
	}
	return rate, nil
}

func (b *baseRateRepo) GetBaseRates(country, network, effectiveAt string, simType SimType) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Where(&BaseRate{Country: country, Network: network, SimType: simType, EffectiveAt: effectiveAt}).Find(&rates)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) UploadBaseRates(rateList []BaseRate) error {
	e := b.Db.GetGormDb().Create(&rateList)
	if e != nil {
		return e.Error
	}

	return nil
}
