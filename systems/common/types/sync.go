package types

import (
	"database/sql/driver"
	"strconv"
)

type SyncStatus uint8

const (
	SyncStatusUnknown SyncStatus = iota
	SyncStatusPending
	SyncStatusCompleted
	SyncStatusFailed
)

func (s *SyncStatus) Scan(value interface{}) error {
	*s = SyncStatus(uint8(value.(int64)))

	return nil
}

func (s SyncStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s SyncStatus) String() string {
	t := map[SyncStatus]string{0: "unknown", 1: "pending", 2: "completed", 3: "failed"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseStatus(value string) SyncStatus {
	i, err := strconv.Atoi(value)
	if err == nil {
		return SyncStatus(i)
	}

	t := map[string]SyncStatus{"unknown": 0, "pending": 1, "completed": 2, "failed": 3}

	v, ok := t[value]
	if !ok {
		return SyncStatus(0)
	}

	return SyncStatus(v)
}
