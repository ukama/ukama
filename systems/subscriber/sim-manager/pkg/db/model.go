package db

import (
	"database/sql/driver"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sim struct {
	ID                 uuid.UUID `gorm:"primaryKey;type:uuid"`
	SubscriberID       uuid.UUID `gorm:"not null;type:uuid"`
	Package            Package
	Iccid              string `gorm:"index:idx_iccid,unique"`
	Msisdn             string
	Imsi               string
	Type               SimType
	Status             SimStatus
	IsPhysical         bool
	ActivationsCount   uint64 `gorm:"default:0"`
	DeactivationsCount uint64 `gorm:"default:0"`
	FirstActivatedOn   time.Time
	LastActivatedOn    time.Time
	AllocatedAt        int64 `gorm:"autoCreateTime"`
	UpdatedAt          time.Time
	TerminatedAt       gorm.DeletedAt `gorm:"index"`
}

type Package struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	SimID     uuid.UUID `gorm:"not null;type:uuid"`
	StartDate time.Time
	EndDate   time.Time
	PlanID    uuid.UUID      `gorm:"not null;type:uuid"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SimType uint8

const (
	SimTypeUnknown SimType = iota
	SimTypeInterNone
	SimTypeInterMnoAll
	SimTypeInterMnoData
	SimTypeInterUkamaAll
)

func (s *SimType) Scan(value interface{}) error {
	*s = SimType(uint8(value.(int64)))
	return nil
}

func (s SimType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s SimType) String() string {
	t := map[SimType]string{0: "unknown", 1: "inter_none", 2: "inter_mno_all", 3: "inter_mno_data", 4: "inter_ukama_all"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseType(value string) SimType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return SimType(i)
	}

	t := map[string]SimType{"unknown": 0, "inter_none": 1, "inter_mno_all": 2, "inter_mno_data": 3, "inter_ukama_all": 4}

	v, ok := t[value]
	if !ok {
		return SimType(0)
	}

	return SimType(v)
}

type SimStatus uint8

const (
	SimStatusUnknown SimStatus = iota
	SimStatusActive
	SimStatusInactive
	SimStatusTerminated
)

func (s *SimStatus) Scan(value interface{}) error {
	*s = SimStatus(uint8(value.(int64)))
	return nil
}

func (s SimStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s SimStatus) String() string {
	t := map[SimStatus]string{0: "unknown", 1: "active", 2: "inactive", 3: "terminated"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseStatus(value string) SimStatus {
	i, err := strconv.Atoi(value)
	if err == nil {
		return SimStatus(i)
	}

	t := map[string]SimStatus{"unknown": 0, "active": 1, "inactive": 2, "terminated": 3}

	v, ok := t[value]
	if !ok {
		return SimStatus(0)
	}

	return SimStatus(v)
}
