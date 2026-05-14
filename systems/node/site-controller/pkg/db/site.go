package db

import (
	"errors"

	"gorm.io/gorm"
)

// ensureSite ensures a registry row exists for siteID before inserting FK children.
func ensureSite(tx *gorm.DB, siteID string) error {
	if siteID == "" {
		return errors.New("site_id is required")
	}
	s := Site{SiteID: siteID}
	return tx.Where(&Site{SiteID: siteID}).FirstOrCreate(&s).Error
}
