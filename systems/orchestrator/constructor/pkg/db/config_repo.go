package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type ConfigsRepo interface {
	Create(c *Configs) error
	Get(name string) (*Configs, error)
	Delete(name string) error
	Update(n *Configs) error
	GetHistory(name string) (*[]Configs, error)
}

type configsRepo struct {
	Db sql.Db
}

func NewConfigsRepo(db sql.Db) *configsRepo {
	return &configsRepo{
		Db: db,
	}
}

func (d *configsRepo) Create(c *Configs) error {
	r := d.Db.GetGormDb().Create(c)
	return r.Error
}

func (d *configsRepo) Update(n *Configs) error {

	err := d.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		var c Configs
		r := tx.Where("name = ?", c.Name).First(&c)
		if r.Error != nil {
			return fmt.Errorf("failed to update %s config. Error while getting: %v", c.Name, r.Error)
		}

		r = tx.Delete(&c)
		if r.Error != nil {
			return fmt.Errorf("failed to update %s config. Error while deleting: %v", c.Name, r.Error)
		}

		r = tx.Create(n)
		if r.Error != nil {
			return fmt.Errorf("failed to  %s config. Error while creating: %v", c.Name, r.Error)
		}

		return nil
	})

	return err
}

func (d *configsRepo) Get(name string) (*Configs, error) {
	var c Configs
	result := d.Db.GetGormDb().Where("name = ?", name).First(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return &c, nil
}

func (d *configsRepo) Delete(name string) error {
	result := d.Db.GetGormDb().Where("name = ?", name).Delete(&Configs{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	}

	return fmt.Errorf("config %s missing", name)
}

func (d *configsRepo) GetHistory(name string) (*[]Configs, error) {
	var c []Configs
	result := d.Db.GetGormDb().Unscoped().Where("name = ?", name).Find(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return &c, nil
}
