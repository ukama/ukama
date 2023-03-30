package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	"gorm.io/gorm"
)

type SiteRepo interface {
	Add(site *Site, nestedFunc func(*Site, *gorm.DB) error) error
	Get(id uuid.UUID) (*Site, error)
	GetByName(netID uuid.UUID, siteName string) (*Site, error)
	GetByNetwork(netID uuid.UUID) ([]Site, error)
	// Update(site *Site) error
	Delete(id uuid.UUID) error
	GetSiteCount(netID uuid.UUID) (int64, error)

	// AttachNodes
	// DetachNodes
}

type siteRepo struct {
	Db sql.Db
}
git 
	}
}

func (s siteRepo) Get(id uuid.UUID) (*Site, error) {
	var site Site

	result := s.Db.GetGormDb().First(&site, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &site, nil
}

func (s siteRepo) GetByName(netID uuid.UUID, siteName string) (*Site, error) {
	var site Site

	result := s.Db.GetGormDb().Where("sites.network_id = ? and sites.name = ?", netID, siteName).First(&site)
	if result.Error != nil {
		return nil, result.Error
	}

	return &site, nil
}

func (s siteRepo) GetByNetwork(netID uuid.UUID) ([]Site, error) {
	var sites []Site
	db := s.Db.GetGormDb()

	result := db.Where(&Site{NetworkId: netID}).Find(&sites)
	if result.Error != nil {
		return nil, result.Error
	}

	return sites, nil
}

func (s siteRepo) Add(site *Site, nestedFunc func(site *Site, tx *gorm.DB) error) error {
	if !validation.IsValidDnsLabelName(site.Name) {
		return fmt.Errorf("invalid name. must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(site, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(site)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (s siteRepo) Delete(siteID uuid.UUID) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("sites.id = ?", siteID).Delete(&Site{})

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (s siteRepo) GetSiteCount(netID uuid.UUID) (int64, error) {
	var count int64
	result := s.Db.GetGormDb().Model(&Site{}).Where("network_id = ?", netID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
