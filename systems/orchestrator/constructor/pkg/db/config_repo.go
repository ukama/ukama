package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type ConfigRepo interface {
	Create(c *Config) error
	Get(name string) (*Config, error)
	GetAll() (*[]Config, error)
	Delete(name string) error
	AddOrg(name string, org Org) error
	DeleteOrg(name string, org Org) error
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
	err := d.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		r := tx.Where("name = ?", c.Name).First(&c)
		if r.Error != nil {
			return fmt.Errorf("failed to update %s config. Error while getting: %v", c.Name, r.Error)
		}

		t := tx.Where("name= ?", c.Name).Delete(&Config{})
		if t.RowsAffected > 0 {
			log.Debugf("Marking old state with delete_at for %s", c.Name)
		}

		result := tx.Create(c)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

/* Added on receiving events */
func (d *configsRepo) AddOrg(name string, org Org) error {

	err := d.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		var c Config
		r := tx.Where("name = ?", c.Name).First(&c)
		if r.Error != nil {
			return fmt.Errorf("failed to update %s config. Error while getting config for %s: %v", name, c.Name, r.Error)
		}

		err := tx.Model(&Config{Name: c.Name}).Omit("Orgs.*").Association("Orgs").Append(&org)
		if err != nil {
			return fmt.Errorf("failed to add %s org to %s config. Error while getting: %v", org.OrgId, c.Name, err)
		}

		return nil
	})

	return err
}

/* Deleted on receiving events */
func (d *configsRepo) DeleteOrg(name string, org Org) error {

	err := d.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		var c Config
		r := tx.Where("name = ?", c.Name).First(&c)
		if r.Error != nil {
			return fmt.Errorf("failed to update %s config. Error while getting config for %s: %v", name, c.Name, r.Error)
		}

		err := tx.Model(&Config{Name: c.Name}).Omit("Orgs.*").Association("Orgs").Delete(&org)
		if err != nil {
			return fmt.Errorf("failed to delete %s org from  %s config. Error while getting: %v", org.OrgId, c.Name, err)
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

func (d *configsRepo) GetAll() (*[]Config, error) {
	var c []Config
	result := d.Db.GetGormDb().Find(&c)
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
