package db

import (
	"github.com/jackc/pgtype"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Orgs struct {
	gorm.Model
	OrgID uuid.UUID
	Deployments []Deployments
}

type Systems struct {
	gorm.Model
	OrgID uuid.UUID
	Deployments []Deployments
}

type Deployments struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex"`
	Source 		string
	Env         string 
	Nodes       []Node
	Systems     []System
}

type Org struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex"`
	Certificate string
	Ip          pgtype.Inet `gorm:"type:inet"`
	Nodes       []Node
	Systems     []System
}

type System struct {
	gorm.Model
	Name        string `gorm:"unique;type:string;uniqueIndex:name_idx_case_insensetive,expression:lower(name);not null"`
	Uuid        string `gorm:"type:uuid;unique"`
	Certificate string
	Ip          pgtype.Inet `gorm:"type:inet"`
	Port        int32
	OrgID       uint
	Org         Org
	Health      uint32 `gorm:"default:100"`
}
