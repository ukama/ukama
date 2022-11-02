package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pb"
)

// declare interface so that we can mock it
type BaseRateRepo interface {
	GetBaseRate(Id int64) (*Rate, error)
	GetBaseRates(country ,network,simType string)([]*pb.Rate, error)
	UploadBaseRates(fileUrl, EffectiveAt,simType string)([]*pb.Rate, error)
}

type baseRateRepo struct {
	Db sql.Db
}

func NewBaseRateRepo(db sql.Db) *baseRateRepo {
	return &baseRateRepo{
		Db: db,
	}
}


func (u *baseRateRepo) GetBaseRate(rateId int64) (*Rate, error) {
	var rate Rate
	result := u.Db.GetGormDb().First(&rate, "Id=?",rateId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rate, nil
}

func (b *baseRateRepo) GetBaseRates(network , country, simType string) ([]*pb.Rate, error) {
	var rates []*pb.Rate
		result := b.Db.GetGormDb().Where("country = ? AND network = ?", country,network).Find(&rates)
	if result.Error != nil {
		return nil, result.Error
	}
	return rates, nil
}
func (b *baseRateRepo) UploadBaseRates(fileUrl , EffectiveAt, simType string) ([]*pb.Rate, error) {
	var rates []*pb.Rate
	
	return rates, nil
}