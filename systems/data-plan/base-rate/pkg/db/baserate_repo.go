package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
)

type BaseRateRepo interface {
	GetBaseRateById(uuid uuid.UUID) (*BaseRate, error)
	GetBaseRatesHistoryByNetwork(country, network string) ([]BaseRate, error)
	GetBaseRatesByNetwork(country, network string, effectiveAt time.Time, simType SimType) ([]BaseRate, error)
	GetBaseRatesForPeriod(country, network string, from, to time.Time, simType SimType) ([]BaseRate, error)
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

func (u *baseRateRepo) GetBaseRateById(uuid uuid.UUID) (*BaseRate, error) {
	rate := &BaseRate{}
	result := u.Db.GetGormDb().First(rate, "Uuid=?", uuid)
	if result.Error != nil {
		return nil, result.Error
	}
	return rate, nil
}

func (b *baseRateRepo) GetBaseRatesHistoryByNetwork(country, network string) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Unscoped().Where(&BaseRate{Country: country, Network: network}).Find(&rates)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) GetBaseRatesByNetwork(country, network string, effectiveAt time.Time, simType SimType) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Model(BaseRate{}).Where("country = ?", country).Where("network = ?", network).
		Where("simType = ?", simType).Where("effective_at <= ?", time.Now()).Order("effective_at desc").Find(&rates).Limit(1)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) GetBaseRatesForPeriod(country, network string, from, to time.Time, simType SimType) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Model(BaseRate{}).Unscoped().Where("country = ?", country).Where("network = ?", network).
		Where("simType = ?", simType).Where("effective_at >= ?", from).Where("effective_at <= ?", to).Find(&rates)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) UploadBaseRates(rateList []BaseRate) error {
	err := b.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		for _, r := range rateList {
			o := &BaseRate{
				Country:     r.Country,
				Network:     r.Network,
				EffectiveAt: r.EffectiveAt,
				SimType:     r.SimType,
			}
			result := tx.Model(BaseRate{}).Delete(o)
			if result.Error != nil {
				if !sql.IsNotFoundError(result.Error) {
					return result.Error
				}
			}

			result = tx.Model(BaseRate{}).Create(r)
			if result.Error != nil {
				return result.Error
			}
		}

		return nil
	})

	return err
}
