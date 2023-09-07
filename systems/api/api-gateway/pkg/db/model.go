package db

import (
	"database/sql/driver"
	"strconv"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Resource struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Status    ResourceStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ResourceStatus uint8

const (
	ResourceStatusUnknown ResourceStatus = iota
	ResourceStatusPending
	ResourceStatusCompleted
	ResourceStatusFailed
)

func (s *ResourceStatus) Scan(value interface{}) error {
	*s = ResourceStatus(uint8(value.(int64)))
	return nil
}

func (s ResourceStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s ResourceStatus) String() string {
	t := map[ResourceStatus]string{0: "unknown", 1: "pending", 2: "completed", 3: "failed"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseStatus(value string) ResourceStatus {
	i, err := strconv.Atoi(value)
	if err == nil {
		return ResourceStatus(i)
	}

	t := map[string]ResourceStatus{"unknown": 0, "pending": 1, "completed": 2, "failed": 3}

	v, ok := t[value]
	if !ok {
		return ResourceStatus(0)
	}

	return ResourceStatus(v)
}
