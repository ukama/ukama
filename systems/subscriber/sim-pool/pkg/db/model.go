package db

import (
	"database/sql/driver"
	"strconv"

	"gorm.io/gorm"
)

type Sim struct {
	gorm.Model
	Iccid          string `gorm:"index:idx_iccid,unique"`
	Msisdn         string
	IsAllocated    bool
	IsFailed       bool
	SimType        SimType
	SmDpAddress    string
	ActivationCode string
	QrCode         string
	IsPhysical     bool
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
