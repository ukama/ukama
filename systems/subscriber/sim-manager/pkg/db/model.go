package db

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sim struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid"`
	SubscriberID uuid.UUID `gorm:"not null;type:uuid"`
	NetworkID    uuid.UUID `gorm:"not null;type:uuid"`
	Iccid        string    `gorm:"index:idx_iccid,unique"`
	Msisdn       string
	Type         SimType
	Status       SimStatus
	IsPhysical   bool
	AllocatedAt  int64 `gorm:"autoCreateTime"`
	UpdatedAt    time.Time
	TerminatedAt gorm.DeletedAt `gorm:"index"`
}

type SimType uint8

const (
	SimTypeUnknown SimType = iota
	SimTypeUkama
	SimTypeTelna
	SimTypeGigs
)

func (s SimType) String() string {
	return []string{"unknown", "ukama", "telna", "gigs"}[s]
}

func (s *SimType) Scan(value interface{}) error {
	*s = SimType(uint8(value.(int64)))
	return nil
}

func (s SimType) Value() (driver.Value, error) {
	return int64(s), nil
}

type SimStatus uint8

const (
	SimStatusUnknown SimStatus = iota
	SimStatusActive
	SimStatusInactive
	SimStatusTerminated
)

func (s SimStatus) String() string {
	return []string{"unknown", "active", "inactive", "terminated"}[s]
}

func (s *SimStatus) Scan(value interface{}) error {
	*s = SimStatus(uint8(value.(int64)))
	return nil
}

func (s SimStatus) Value() (driver.Value, error) {
	return int64(s), nil
}
