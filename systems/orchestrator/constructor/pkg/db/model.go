package db

import (
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Orgs struct {
	gorm.Model
	OrgID       uuid.UUID
	Name        string
	Deployments []Deployments
}

type Deployments struct {
	gorm.Model
	Name   string `gorm:"index:deployment_name_idx,not null"`
	Env    string
	Status uint8
	Values []string
	OrgID  uint
	Org    Orgs
	Config Config
}

type Config struct {
	gorm.Model
	Name          string
	Source        string `gorm:"type:string;index:source_name_idx,not null"`
	SourceVersion string
	Values        []string
}
