package db

import (
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
)

type MarkupsRepo interface {
	GetMarkupRate(uuid uuid.UUID) (*Markups, error)
	CreateMarkupRate(uuid uuid.UUID, markup float64) error
	DeleteMarkupRate(uuid uuid.UUID) error
	UpdateMarkupRate(uuid uuid.UUID, markup float64) error
	GetMarkupRateHistory(uuid uuid.UUID) ([]Markups, error)
}

type markupsRepo struct {
	Db sql.Db
}

func NewMarkupsRepo(db sql.Db) *markupsRepo {
	return &markupsRepo{
		Db: db,
	}
}

func (m *markupsRepo) CreateMarkupRate(uuid uuid.UUID, markup float64) error {
	rate := Markups{
		OwnerId: uuid,
		Markup:  markup,
	}
	result := m.Db.GetGormDb().Model(&Markups{}).Create(&rate)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *markupsRepo) GetMarkupRate(uuid uuid.UUID) (*Markups, error) {
	rate := &Markups{}
	result := m.Db.GetGormDb().Model(&Markups{}).Where("owner_id=?", uuid).First(rate)
	if result.Error != nil {
		return nil, result.Error
	}
	return rate, nil
}

func (m *markupsRepo) DeleteMarkupRate(uuid uuid.UUID) error {
	rate := &Markups{}
	result := m.Db.GetGormDb().Model(&Markups{}).Where("owner_id=?", uuid).Delete(rate)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *markupsRepo) UpdateMarkupRate(uuid uuid.UUID, mrate float64) error {
	err := m.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		markup := &Markups{
			OwnerId: uuid,
		}
		result := tx.Model(Markups{}).Where("owner_id = ?", uuid).Delete(markup)
		if result.Error != nil {
			if !sql.IsNotFoundError(result.Error) {
				return result.Error
			}
		}

		new := &Markups{
			OwnerId: uuid,
			Markup:  mrate,
		}
		result = tx.Model(Markups{}).Create(new)
		if result.Error != nil {
			return result.Error
		}

		return nil

	})

	return err
}

func (m *markupsRepo) GetMarkupRateHistory(uuid uuid.UUID) ([]Markups, error) {
	rates := []Markups{}
	result := m.Db.GetGormDb().Model(&Markups{}).Unscoped().Where("owner_id=?", uuid).Find(&rates)
	if result.Error != nil {
		return nil, result.Error
	}
	return rates, nil
}
