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
	System Systems
}

type Systems struct {
	gorm.Model
	Name         string
	Chart        string `gorm:"type:string;index:chart_name_idx,not null"`
	Version      string
	Values       []string
	DeploymentID uint
	Deployment   []Deployments
}
