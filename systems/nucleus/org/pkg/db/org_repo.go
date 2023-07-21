package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
)

type OrgRepo interface {
	/* Orgs */
	Add(org *Org, nestedFunc func(*Org, *gorm.DB) error) error
	Get(id uuid.UUID) (*Org, error)
	GetByName(name string) (*Org, error)
	GetByOwner(uuid uuid.UUID) ([]Org, error)
	GetByMember(uuid uuid.UUID) ([]Org, error)
	GetAll() ([]Org, error)
	GetOrgCount() (int64, int64, error)
	// Update(id uint) error
	// Deactivate(id uint) error
	// Delete(id uint) error

}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) OrgRepo {
	return &orgRepo{
		Db: db,
	}
}

func (r *orgRepo) Add(org *Org, nestedFunc func(*Org, *gorm.DB) error) (err error) {
	if !validation.IsValidDnsLabelName(org.Name) {
		return fmt.Errorf("invalid name must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	err = r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(org, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		d := tx.Create(org)
		if d.Error != nil {
			return d.Error
		}

		return nil
	})

	return err
}

func (r *orgRepo) Get(id uuid.UUID) (*Org, error) {
	var org Org

	result := r.Db.GetGormDb().First(&org, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

func (r *orgRepo) GetByName(name string) (*Org, error) {
	var org Org

	result := r.Db.GetGormDb().Where(&Org{Name: name}).First(&org)
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

func (r *orgRepo) GetByOwner(uuid uuid.UUID) ([]Org, error) {
	var orgs []Org

	result := r.Db.GetGormDb().Where(&Org{Owner: uuid}).Find(&orgs)
	if result.Error != nil {
		return nil, result.Error
	}

	return orgs, nil
}

func (r *orgRepo) GetByMember(uuid uuid.UUID) ([]Org, error) {
	var membOrgs []Org
	//TODO: fix this
	// result := r.Db.GetGormDb().Where(&Org{Uuid: uuid}).Find(&membOrgs)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	return membOrgs, nil
}

func (r *orgRepo) GetAll() ([]Org, error) {
	var orgs []Org

	result := r.Db.GetGormDb().Where(&Org{}).Find(&orgs)
	if result.Error != nil {
		return nil, result.Error
	}

	return orgs, nil
}

func (r *orgRepo) GetOrgCount() (int64, int64, error) {
	var activeOrgCount int64
	var deactiveOrgCount int64

	result := r.Db.GetGormDb().Model(&Org{}).
		Where("deactivated = ?", false).Count(&activeOrgCount)

	if result.Error != nil {
		return 0, 0, result.Error
	}

	result = r.Db.GetGormDb().Model(&Org{}).
		Where("deactivated = ?", true).Count(&deactiveOrgCount)

	if result.Error != nil {
		return 0, 0, result.Error
	}

	return activeOrgCount, deactiveOrgCount, nil
}
