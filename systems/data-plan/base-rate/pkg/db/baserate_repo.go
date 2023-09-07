package db

import (
	"time"

	log "github.com/sirupsen/logrus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
)

type BaseRateRepo interface {
	GetBaseRateById(uuid uuid.UUID) (*BaseRate, error)
	GetBaseRatesHistoryByCountry(country, provider string, sType SimType) ([]BaseRate, error)
	GetBaseRatesByCountry(country, provider string, simType SimType) ([]BaseRate, error)
	GetBaseRatesForPeriod(country, provider string, from, to time.Time, simType SimType) ([]BaseRate, error)
	GetBaseRatesForPackage(country, provider string, from, to time.Time, simType SimType) ([]BaseRate, error)
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
	result := u.Db.GetGormDb().First(rate, "uuid=?", uuid)
	if result.Error != nil {
		return nil, result.Error
	}
	return rate, nil
}

func (b *baseRateRepo) GetBaseRatesHistoryByCountry(country, provider string, sType SimType) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Unscoped().Where(&BaseRate{Country: country, Provider: provider, SimType: sType}).Find(&rates)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) GetBaseRatesByCountry(country, provider string, simType SimType) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Model(BaseRate{}).Where("country = ?", country).Where("provider = ?", provider).
		Where("sim_type = ?", simType).Where("effective_at <= ?", time.Now()).Order("effective_at desc").Limit(1).Find(&rates)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) GetBaseRatesForPeriod(country, provider string, from, to time.Time, simType SimType) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Model(BaseRate{}).Unscoped().Where("country = ?", country).Where("provider = ?", provider).
		Where("sim_type = ?", simType).Where("effective_at >= ?", from).Where("effective_at <= ?", to).Find(&rates)

	if result.Error != nil {
		return nil, result.Error
	}

	return rates, nil
}

func (b *baseRateRepo) GetBaseRatesForPackage(country, provider string, from, to time.Time, simType SimType) ([]BaseRate, error) {
	var rates []BaseRate
	result := b.Db.GetGormDb().Model(BaseRate{}).Unscoped().Where("country = ?", country).Where("provider = ?", provider).
		Where("sim_type = ?", simType).Where("effective_at <= ?", from).Where("end_at >= ?", to).Order("created_at desc").Limit(1).Find(&rates)

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
				Provider:    r.Provider,
				EffectiveAt: r.EffectiveAt,
				SimType:     r.SimType,
			}
			result := tx.Model(BaseRate{}).Where("country = ?", r.Country).Where("provider = ?", r.Provider).
				Where("sim_type = ?", r.SimType).Where("effective_at = ?", r.EffectiveAt).Delete(o)
			if result.Error != nil {
				if !sql.IsNotFoundError(result.Error) {
					log.Errorf("Error deleting rate %+v . Error %s", o, result.Error.Error())
					return result.Error
				}
			}

			result = tx.Model(BaseRate{}).Create(&r)
			if result.Error != nil {
				log.Errorf("Error creating rate %+v . Error %s", r, result.Error.Error())
				return result.Error
			}
		}

		return nil
	})

	return err
}
