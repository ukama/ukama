package db

import (
	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/common/sql"
)

type MarkupsRepo interface {
	GetMarkupRate(uuid uuid.UUID) (*Markups, error)
	DeleteMarkupRate(uuid uuid.UUID) error
	UpdatedMarkupRate(uuid uuid.UUID) error
}

type markupsRepo struct {
	Db sql.Db
}

func NewMarkupsRepo(db sql.Db) *markupsRepo {
	return &markupsRepo{
		Db: db,
	}
}

func (u *markupsRepo) GetMarkupRate(uuid uuid.UUID) (*Markups, error) {
	rate := &Markups{}
	result := u.Db.GetGormDb().Model(&Markups{}).First("owner_id=?", uuid)
	if result.Error != nil {
		return nil, result.Error
	}
	return rate, nil
}

func (u *markupsRepo) DeleteMarkupRate(uuid uuid.UUID) error {
	rate := &Markups{}
	result := u.Db.GetGormDb().Model(&Markups{}).Where("owner_id=?", uuid).Delete(rate)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *markupsRepo) UpdateMarkupRate(uuid uuid.UUID, markup int64) error {
	rate := Markups{
		OwnerId: uuid,
		Markup:  markup,
	}
	result := u.Db.GetGormDb().Model(&Markups{}).Updates(rate)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
