package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type ConfigRepo interface {
	Create(c *Config) error
	Get(name string) (*Config, error)
	Delete(name string) error
	Update(n *Config) error
	GetHistory(name string) (*[]Config, error)
}

type configsRepo struct {
	Db sql.Db
}

func NewConfigRepo(db sql.Db) *configsRepo {
	return &configsRepo{
		Db: db,
	}
}

func (d *configsRepo) Create(c *Config) error {
	r := d.Db.GetGormDb().Create(c)
	return r.Error
}

func (d *configsRepo) Update(n *Config) error {

	err := d.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		var c Config
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

func (d *configsRepo) Get(name string) (*Config, error) {
	var c Config
	result := d.Db.GetGormDb().Where("name = ?", name).First(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return &c, nil
}

func (d *configsRepo) Delete(name string) error {
	result := d.Db.GetGormDb().Where("name = ?", name).Delete(&Config{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	}

	return fmt.Errorf("config %s missing", name)
}

func (d *configsRepo) GetHistory(name string) (*[]Config, error) {
	var c []Config
	result := d.Db.GetGormDb().Unscoped().Where("name = ?", name).Find(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return &c, nil
}
