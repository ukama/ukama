package db

import (
	"database/sql/driver"
	"strconv"
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Sim struct {
	Id                 uuid.UUID `gorm:"primaryKey;type:uuid"`
	SubscriberId       uuid.UUID `gorm:"not null;type:uuid"`
	NetworkId          uuid.UUID `gorm:"not null;type:uuid"`
	OrgId              uuid.UUID `gorm:"not null;type:uuid"`
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
	Id        uuid.UUID `gorm:"primaryKey;type:uuid"`
	SimId     uuid.UUID `gorm:"uniqueIndex:unique_sim_package_is_active,where:is_active is true;not null;type:uuid"`
	StartDate time.Time
	EndDate   time.Time
	PackageId uuid.UUID `gorm:"not null;type:uuid"`
	IsActive  bool      `gorm:"uniqueIndex:unique_sim_package_is_active,where:is_active is true;default:false"`
}

func (p Package) IsExpired() bool {
	return p.EndDate.Before(time.Now())
}

type SimType uint8

const (
	SimTypeUnknown SimType = iota
	SimTypeTest
	SimTypeOperatorData
	SimTypeUkamaData
)

func (s *SimType) Scan(value interface{}) error {
	*s = SimType(uint8(value.(int64)))
	return nil
}

func (s SimType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s SimType) String() string {
	t := map[SimType]string{0: "unknown", 1: "test", 2: "operator_data", 3: "ukama_data"}

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

	t := map[string]SimType{"unknown": 0, "test": 1, "operator_data": 2, "ukama_data": 3}

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
