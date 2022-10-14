package db

import (
	"github.com/ukama/ukama/systems/common/errors"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

func makeUserOrgExist(db *gorm.DB, orgName string) (*Org, error) {
	org := Org{
		Name: orgName,
	}

	d := db.First(&org, "name = ?", orgName)
	if d.Error != nil {
		if sql.IsNotFoundError(d.Error) {
			d2 := db.Create(&org)
			if d2.Error != nil {
				return nil, errors.Wrap(d2.Error, "error adding the org")
			}
		} else {
			return nil, errors.Wrap(d.Error, "error finding the org")
		}
	}

	return &org, nil
}
