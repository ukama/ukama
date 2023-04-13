package db

import (
	"database/sql/driver"
	"strconv"
)

type CallUnitType uint8

const (
	CallUnitTypeUnknown CallUnitType = iota
	CallUnitTypeSec                  = 1
	CallUnitTypeMin                  = 2
	CallUnitTypeHours                = 3
)

func (s *CallUnitType) Scan(value interface{}) error {
	*s = CallUnitType(uint8(value.(int64)))
	return nil
}

func (s CallUnitType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s CallUnitType) String() string {
	t := map[CallUnitType]string{0: "unknown", 1: "seconds", 2: "minutes", 3: "hours"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseCallUnitType(value string) CallUnitType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return CallUnitType(i)
	}

	t := map[string]CallUnitType{"unknown": 0, "seconds": 1, "minutes": 2, "hours": 3}

	v, ok := t[value]
	if !ok {
		return CallUnitType(0)
	}

	return CallUnitType(v)
}
