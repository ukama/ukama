package db

import (
	"database/sql/driver"
	"strconv"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type BaseRate struct {
	gorm.Model
	Uuid        uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	Country     string
	Network     string
	Vpmn        string
	Imsi        int64
	SmsMo       float64 `gorm:"type:float"`
	SmsMt       float64 `gorm:"type:float"`
	Data        float64 `gorm:"type:float"`
	X2g         bool    `gorm:"type:bool; default:false"`
	X3g         bool    `gorm:"type:bool; default:false"`
	X5g         bool    `gorm:"type:bool; default:false"`
	Lte         bool    `gorm:"type:bool; default:false"`
	LteM        bool    `gorm:"type:bool; default:false"`
	Apn         string
	EffectiveAt string
	EndAt       string
	SimType     SimType
	Currency    string `gorm:"not null; default:Dollar"`
}

type SimType uint8

const (
	SimTypeUnknown      SimType = iota
	SimTypeTest                 = 1
	SimTypeOperatorData         = 2
	SimTypeUkamaData            = 3
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
