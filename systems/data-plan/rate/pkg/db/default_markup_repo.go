package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

type DefaultMarkupRepo interface {
	GetDefaultMarkupRate() (*DefaultMarkup, error)
	CreateDefaultMarkupRate(markup float64) error
	DeleteDefaultMarkupRate() error
	UpdateDefaultMarkupRate(markup float64) error
	GetDefaultMarkupRateHistory() ([]DefaultMarkup, error)
}

type defaultMarkupRepo struct {
	Db sql.Db
}

func NewDefaultMarkupRepo(db sql.Db) *defaultMarkupRepo {
	return &defaultMarkupRepo{
		Db: db,
	}
}

func (m *defaultMarkupRepo) CreateDefaultMarkupRate(markup float64) error {
	rate := DefaultMarkup{
		Markup: markup,
	}
	result := m.Db.GetGormDb().Model(&DefaultMarkup{}).Create(&rate)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *defaultMarkupRepo) GetDefaultMarkupRate() (*DefaultMarkup, error) {
	rate := &DefaultMarkup{}
	result := m.Db.GetGormDb().Model(&DefaultMarkup{}).First(&rate)
	if result.Error != nil {
		return nil, result.Error
	}
	return rate, nil
}

func (m *defaultMarkupRepo) DeleteDefaultMarkupRate() error {
	rate := &DefaultMarkup{}
	result := m.Db.GetGormDb().Model(&DefaultMarkup{}).Where("deleted_at=?", nil).Delete(rate)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *defaultMarkupRepo) UpdateDefaultMarkupRate(markup float64) error {

	err := m.DeleteDefaultMarkupRate()
	if err != nil {
		return err
	}

	err = m.CreateDefaultMarkupRate(markup)
	if err != nil {
		return err
	}

	return nil
}

func (m *defaultMarkupRepo) GetDefaultMarkupRateHistory() ([]DefaultMarkup, error) {
	rate := []DefaultMarkup{}
	result := m.Db.GetGormDb().Model(&DefaultMarkup{}).Unscoped().Find(&rate)
	if result.Error != nil {
		return nil, result.Error
	}
	return rate, nil
}
