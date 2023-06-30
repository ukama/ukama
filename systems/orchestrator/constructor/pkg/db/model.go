package db

import (
	"database/sql/driver"
	"strings"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Deployment struct {
	gorm.Model
	Code     uuid.UUID `gorm:"type=uuid;uniqueIndex:deployment_id,not null"`
	Name     string    `gorm:"index:deployment_name_idx,not null"` /* combination od org_name+system-name */
	OrgID    uuid.UUID `gorm:"index:org_idx,not null;"`
	Org      Org       `gorm:"references:OrgId"`
	Env      string
	Status   DeploymentStatus `gorm:"type=uint"`
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
	Source  string
	Chart   string `gorm:"type:string;not null"`
	Version string `gorm:"type:string;not null"`
	Values  []string
	Orgs    []Org
}

type DeploymentStatus uint8

const (
	Unkwown    DeploymentStatus = iota
	Waiting    DeploymentStatus = 1
	Working    DeploymentStatus = 2
	Failed     DeploymentStatus = 3
	Completetd DeploymentStatus = 4
	Scheduled  DeploymentStatus = 5
)

func (s *DeploymentStatus) Scan(value interface{}) error {
	*s = DeploymentStatus(uint8(value.(int64)))

	return nil
}

func (s DeploymentStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s DeploymentStatus) String() string {
	ns := map[DeploymentStatus]string{
		Unkwown:    "unkown",
		Waiting:    "waiting",
		Working:    "working-on-it",
		Failed:     "failed",
		Completetd: "completeted",
		Scheduled:  "scheduled",
	}

	return ns[s]
}

func ParseDeploymentStatus(s string) DeploymentStatus {
	switch strings.ToLower(s) {
	case "waiting":
		return Waiting
	case "working-on-it":
		return Working
	case "failed":
		return Failed
	case "completeted":
		return Completetd
	case "scheduled":
		return Scheduled
	default:
		return Unkwown
	}
}
