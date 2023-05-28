package db

import (
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Deployments struct {
	gorm.Model
	Name    string    `gorm:"index:deployment_name_idx,not null"`
	OrgID   uuid.UUID `gorm:"index:org_system_idx,not null;index:org_idx"`
	SysName string    `gorm:"uniqueIndex:system_name_idx,not null"` /* org specifc name for a system. has to be unique throughout the orgs */
	Env     string
	Status  uint8
	Values  []string
	Details string
}

/* System configurations */
type Systems struct {
	gorm.Model
	Name        string `gorm:"index:system_name_idx,not null"`
	SystemId    uint8  `gorm:"index:system_idx,not null` // registry, billing, subscriber etc
	Chart       string `gorm:"type:string;not null"`
	Version     string `gorm:"type:string;not null"`
	Values      []string
	Deployments []Deployments
}
