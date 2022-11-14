package db

import (
	"gorm.io/gorm"
)

type Network struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex:network_name_org_idx"`
	OrgID       uint   `gorm:"uniqueIndex:network_name_org_idx"`
	Org         *Org
	Deactivated bool
}

type Org struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex"`
	Deactivated bool
}

type Site struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex"`
	Deactivated bool
}
