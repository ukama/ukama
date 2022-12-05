package db

import (
	"gorm.io/gorm"
)

type Org struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex"`
	Deactivated bool
}

type Network struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex:network_name_org_idx"`
	OrgID       uint   `gorm:"uniqueIndex:network_name_org_idx"`
	Org         *Org
	Deactivated bool
}

type Site struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex:site_name_network_idx"`
	NetworkID   uint   `gorm:"uniqueIndex:site_name_network_idx"`
	Network     *Network
	Deactivated bool
}
