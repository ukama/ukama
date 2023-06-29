package db

import (
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Deployment struct {
	gorm.Model
	Code     uuid.UUID `gorm:"type=uuid;uniqueIndex:deployment_id,not null"`
	Name     string    `gorm:"index:deployment_name_idx,not null"` /* combination od org_name+system-name */
	OrgID    uuid.UUID `gorm:"index:org_idx,not null;"`
	Org      Org       `gorm:"references:OrgId"`
	SysName  string    `gorm:"uniqueIndex:system_name_idx,not null"` /* org specifc name for a system. has to be unique throughout the orgs */
	Env      string
	Status   uint8
	Values   []string
	Details  string
	ConfigID uint
	Config   Config
}

type Org struct {
	gorm.Model
	OrgId   uuid.UUID `gorm:"index:org_idx,not null;"`
	OrgName string    /* Get this from registry. No need to store here just to verify if the irg is valid */
	Values  []string
	Config  []Config
}

/* System configurations */
type Config struct {
	gorm.Model
	Name    string `gorm:"uniqueIndex:config_name_idx,not null"`
	Chart   string `gorm:"type:string;not null"`
	Version string `gorm:"type:string;not null"`
	Values  []string
	Orgs    []Org
}
