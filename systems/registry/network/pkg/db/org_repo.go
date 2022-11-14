package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/errors"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/validation"

	"github.com/ukama/ukama/systems/common/sql"
)

type OrgRepo interface {
	Add(org *Org, nestedFunc ...func() error) error
	Get(id uint) (*Org, error)
	GetByName(name string) (*Org, error)
	MakeUserOrgExist(orgName string) (*Org, error)
}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) OrgRepo {
	return &orgRepo{
		Db: db,
	}
}

// Add adds an organisation. Add nestedFunc to execute an action inside transations
// if one of the nestedFunc returns error then Add action is rolled back
func (r *orgRepo) Add(org *Org, nestedFunc ...func() error) (err error) {
	if !validation.IsValidDnsLabelName(org.Name) {
		return fmt.Errorf("invalid name must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	err = r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(org)
	}, nestedFunc...)

	return err
}

func (r *orgRepo) Get(id uint) (*Org, error) {
	var org Org

	result := r.Db.GetGormDb().First(&org, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

func (r *orgRepo) GetByName(name string) (*Org, error) {
	var org Org

	result := r.Db.GetGormDb().First(&org, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

func (r *orgRepo) MakeUserOrgExist(orgName string) (*Org, error) {
	org := Org{
		Name: orgName,
	}

	d := r.Db.GetGormDb().First(&org, "name = ?", orgName)
	if d.Error != nil {
		if sql.IsNotFoundError(d.Error) {
			d2 := r.Db.GetGormDb().Create(&org)
			if d2.Error != nil {
				return nil, errors.Wrap(d2.Error, "error adding the org")
			}

		} else {
			return nil, errors.Wrap(d.Error, "error finding the org")
		}
	}

	return &org, nil
}
