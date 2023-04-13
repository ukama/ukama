package db

import (
	"database/sql/driver"
	"strconv"
)

type PackageType uint8

const (
	PackageTypeUnknown  PackageType = iota
	PackageTypePrepaid              = 1
	PackageTypePostpaid             = 2
)

func (s *PackageType) Scan(value interface{}) error {
	*s = PackageType(uint8(value.(int64)))
	return nil
}

func (s PackageType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s PackageType) String() string {
	t := map[PackageType]string{0: "unknown", 1: "prepaid", 2: "postpaid"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParsePackageType(value string) PackageType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return PackageType(i)
	}

	t := map[string]PackageType{"unknown": 0, "prepaid": 1, "postpaid": 2}

	v, ok := t[value]
	if !ok {
		return PackageType(0)
	}

	return PackageType(v)
}
