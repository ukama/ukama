package db

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Subscriber struct {
	gorm.Model
	SubscriberID          uuid.UUID `gorm:"type:uuid"`
	FullName              string
	Email                 string
	PhoneNumber           string
	DOB                   *time.Time
	ProofOfIdentification string
	IdSerial              string
	Address               string
	Sims                  []*Sim `gorm:"one2many:attached_sims"`
}

type Sim struct {
	gorm.Model
	SimID                uuid.UUID `gorm:"type:uuid"`
	NetworkID            uuid.UUID `gorm:"type:uuid"`
	SubscriberID         uuid.UUID `gorm:"type:uuid"`
	OrgID                uuid.UUID `gorm:"type:uuid"`
	ActivePackageID      uuid.UUID `gorm:"type:uuid"`
	IMSI                 string
	SimManager           string
	Packages             []*Package `gorm:"many2many:attached_packages"`
	ActivationsCount     int64
	DeactivationsCount   int64
	LastActivationDate   *time.Time
	LastDeactivationDate *time.Time
	ICCID                string
	MSISDN               string
	State                SimState `gorm:"type:varchar(255)"`
	IsPrepaid            bool
	SimType              string
}

type Package struct {
	gorm.Model
	Status                     bool
	SimID                      uuid.UUID `gorm:"type:uuid"`
	PackageID                  uuid.UUID `gorm:"type:uuid"`
	PackageStartActivationDate *time.Time
	PackageEndActivationDate   *time.Time
}
type SimState string

const (
	SimStateReady     SimState = "Ready"
	SimStateNew       SimState = "New"
	SimStateActive    SimState = "Active"
	SimStateInactive  SimState = "Inactive"
	SimStateSuspended SimState = "Suspended"
	SimStateInvalid   SimState = "Invalid"
	SimStateExpired   SimState = "Expired"
	SimStateStolen    SimState = "Stolen"
)
