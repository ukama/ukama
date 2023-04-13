package db

import (
	"database/sql/driver"
	"strconv"
)

type DataUnitType uint8

const (
	DataUnitTypeUnknown DataUnitType = iota
	DataUnitTypeB                    = 1
	DataUnitTypeKB                   = 2
	DataUnitTypeMB                   = 3
	DataUnitTypeGB                   = 4
)

func (s *DataUnitType) Scan(value interface{}) error {
	*s = DataUnitType(uint8(value.(int64)))
	return nil
}

func (s DataUnitType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s DataUnitType) String() string {
	t := map[DataUnitType]string{0: "unknown", 1: "Bytes", 2: "Kilobytes", 3: "Megabyts", 4: "Gigabytes"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseDataUnitType(value string) DataUnitType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return DataUnitType(i)
	}

	t := map[string]DataUnitType{"unknown": 0, "Bytes": 1, "Kilobytes": 2, "Megabyts": 3, "Gigabytes": 4}

	v, ok := t[value]
	if !ok {
		return DataUnitType(0)
	}

	return DataUnitType(v)
}
