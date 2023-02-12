package db

import (
	"database/sql/driver"
	"strconv"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	PackageID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name         string
	SimType     SimType
	OrgID       uuid.UUID `gorm:"not null;type:uuid;index"`
	Active       bool
	Duration     uint
	SmsVolume   uint
	DataVolume  uint
	VoiceVolume uint
	OrgRatesID uint
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