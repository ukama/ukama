package db

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/validation"

	"github.com/ukama/ukama/systems/common/sql"
)

type OrgRepo interface {
	Add(org *Org, nestedFunc ...func() error) error
	Get(id int) (*Org, error)
	GetByName(name string) (*Org, error)
	Delete(name string) error
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

func (r *orgRepo) Get(id int) (*Org, error) {
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

func (r *orgRepo) Delete(name string) error {
	return r.Db.GetGormDb().Delete(&Org{}, "name = ?", name).Error
}
