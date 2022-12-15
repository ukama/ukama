package db

import (
	"time"

	"gorm.io/gorm"
)

type Subscriber struct {
	gorm.Model
	SubscriberId string `gorm:"type:string;uniqueIndex:subscriber_id, where:deleted_at is null;size:23"`
	FullName    string
	Email   string
	PhoneNumber   string
	DateOfBirth *time.Time
	PassportNumber string
	Address string
	Sims []*Sim `gorm:"many2many:attached_sims"`
}

type Sim struct {
	gorm.Model
	SimId string `gorm:"type:string;uniqueIndex:sim_id, where:deleted_at is null;size:23"`
	NetworkId string
	SubscriberId string
	OrgId string
	Imsi string
	SimProvider string
	Packages []*Package `gorm:"many2many:attached_packages"`
	ActivationsCount int64
	DeactivationsCount int64
	LastActivationDate *time.Time
	LastDeactivationDate *time.Time
	Iccid string
	Msisdn string
	Status bool
	IsPrepaid bool
	SimType string
}
type Package struct {
	gorm.Model
	status bool
	PackageId string `gorm:"type:string;uniqueIndex:package_id, where:deleted_at is null;size:23"`
	PackageStartActivationDate *time.Time
}
