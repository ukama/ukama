package db

import (
	"fmt"

	"github.com/ukama/ukamaX/cloud/registry/pkg/validation"

	"github.com/ukama/ukamaX/common/sql"
)

type OrgRepo interface {
	Add(org *Org) error
	Get(id int) (*Org, error)
	GetByName(name string) (*Org, error)
}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) OrgRepo {
	return &orgRepo{
		Db: db,
	}
}

func (r *orgRepo) Add(org *Org) error {
	if !validation.IsValidDnsLabelName(org.Name) {
		return fmt.Errorf("invalid name must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	d := r.Db.GetGormDb().Create(org)
	return d.Error
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
