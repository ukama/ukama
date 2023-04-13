package db

import (
	"database/sql/driver"
	"strconv"
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	Uuid         uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name         string
	SimType      SimType
	OrgId        uuid.UUID       `gorm:"not null;type:uuid;index"`
	Active       bool            `gorm:"not null; default:false"`
	Duration     uint            `gorm:"not null; default:0"`
	SmsVolume    uint            `gorm:"not null; default:0"`
	DataVolume   uint            `gorm:"not null; default:0"`
	VoiceVolume  uint            `gorm:"not null; default:0"`
	Dlbr         uint64          `gorm:"not null; default:0"`
	Ulbr         uint64          `gorm:"not null; default:0"`
	Rate         *Rate           `gorm:"not null"`
	Type         PackageType     `gorm:"not null; type: uint; default:2"`
	DataUnits    DataUnitType    `gorm:"not null; type: uint; default:3"`
	CallUnits    CallUnitType    `gorm:"not null; type: uint; default:2"`
	MessageUnits MessageUnitType `gorm:"not null; type: uint; default:1"`
	Flatrate     bool            `gorm:"not null; default:false"`
	Currency     string          `gorm:"not null; default:Dollar"`
	EffectiveAt  time.Time       `gorm:"not null"`
	EndAt        time.Time
	Country      string `gorm:"not null;type:string"`
	Provider     string `gorm:"not null;type:string"`
	Apn          string
}

type Rate struct {
	gorm.Model
	BaseRateId uuid.UUID `gorm:"uniqueIndex:uuid_idx,where:deleted_at is null;not null;type:uuid"`
	Markup     float64   `gorm:"type:float"`
	Flatrate   float64   `gorm:"type:float"`
	Vpmn       string
	Imsi       int64
	SmsMo      float64 `gorm:"type:float"`
	SmsMt      float64 `gorm:"type:float"`
	Data       float64 `gorm:"type:float"`
	X2g        bool    `gorm:"type:bool; default:false"`
	X3g        bool    `gorm:"type:bool; default:false"`
	X5g        bool    `gorm:"type:bool; default:false"`
	Lte        bool    `gorm:"type:bool; default:false"`
	LteM       bool    `gorm:"type:bool; default:false"`
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
