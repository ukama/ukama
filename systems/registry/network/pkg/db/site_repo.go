package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/validation"
	"gorm.io/gorm"
)

type SiteRepo interface {
	Add(site *Site) error
	Get(id uint) (*Site, error)
	GetByName(netID uint, siteName string) (*Site, error)
	GetByNetwork(netID uint) ([]Site, error)
	// Update(site *Site) error
	Delete(id uint) error

	// AttachNodes
	// DetachNodes
}

type siteRepo struct {
	Db sql.Db
}

func NewSiteRepo(db sql.Db) SiteRepo {
	return &siteRepo{
		Db: db,
	}
}

func (s siteRepo) Get(id uint) (*Site, error) {
	var site Site

	result := s.Db.GetGormDb().First(&site, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &site, nil
}

func (s siteRepo) GetByName(netID uint, siteName string) (*Site, error) {
	var site Site

	result := s.Db.GetGormDb().Where("sites.network_id = ? and sites.name = ?", netID, siteName).First(&site)
	if result.Error != nil {
		return nil, result.Error
	}

	return &site, nil
}

func (s siteRepo) GetByNetwork(netID uint) ([]Site, error) {
	var sites []Site
	db := s.Db.GetGormDb()

	result := db.Where(&Site{NetworkID: netID}).Find(&sites)
	if result.Error != nil {
		return nil, result.Error
	}

	return sites, nil
}

func (s siteRepo) Add(site *Site) error {
	if !validation.IsValidDnsLabelName(site.Name) {
		return fmt.Errorf("invalid name. must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	result := s.Db.GetGormDb().Create(site)

	return result.Error
}

func (s siteRepo) Delete(siteID uint) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("sites.id = ?", siteID).Delete(&Site{})

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}
