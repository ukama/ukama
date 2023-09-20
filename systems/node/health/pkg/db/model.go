package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type Health struct {
	Id        uuid.UUID  `gorm:"primaryKey;type:uuid"`
	NodeId    uuid.UUID  `gorm:"type:uuid"`
	Name      string     `gorm:"not null"`
	Version   string     `gorm:"not null"`
	Status    Status     `gorm:"not null"`
	Timestamp string     `gorm:"not null"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
}

type Status uint8

const (
	Unknown Status = 0
	Running Status = 1
	Failure Status = 2
)

func (e *Status) Scan(value interface{}) error {
	*e = Status(uint8(value.(int64)))

	return nil
}

func (e Status) Value() (uint8, error) {
	return uint8(e), nil
}
