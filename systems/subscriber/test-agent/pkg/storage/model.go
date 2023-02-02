package storage

import (
	"errors"
	"strconv"
)

var ErrNotFound = errors.New("sim not found")

type Storage interface {
	Get(key string) (*SimInfo, error)
	Put(string, *SimInfo) error
	Delete(key string) error
}

type SimInfo struct {
	Iccid  string `json:"iccid"`
	Imsi   string `json:"imsi"`
	Status string `json:"status"`
}

type SimStatus uint8

const (
	SimStatusUnknown SimStatus = iota
	SimStatusActive
	SimStatusInactive
)

func (s SimStatus) String() string {
	t := map[SimStatus]string{0: "unknown", 1: "active", 2: "inactive"}

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

	t := map[string]SimStatus{"unknown": 0, "active": 1, "inactive": 2}

	v, ok := t[value]
	if !ok {
		return SimStatus(0)
	}

	return SimStatus(v)
}
