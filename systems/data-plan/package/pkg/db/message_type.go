package db

import (
	"database/sql/driver"
	"strconv"
)

type MessageUnitType uint8

const (
	MessageUnitTypeUnknown MessageUnitType = iota
	MeesageUnitTypeInt                     = 1
)

func (s *MessageUnitType) Scan(value interface{}) error {
	*s = MessageUnitType(uint8(value.(int64)))
	return nil
}

func (s MessageUnitType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s MessageUnitType) String() string {
	t := map[MessageUnitType]string{0: "unknown", 1: "int"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseMessageType(value string) MessageUnitType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return MessageUnitType(i)
	}

	t := map[string]MessageUnitType{"unknown": 0, "int": 1}

	v, ok := t[value]
	if !ok {
		return MessageUnitType(0)
	}

	return MessageUnitType(v)
}
