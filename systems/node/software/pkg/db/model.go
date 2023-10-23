package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type Software struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string
	Tag         string
	Space 	 string
	Status      Status `gorm:"type:smallint" default:"0"`
	ReleaseDate time.Time
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null"`
	DeletedAt   *time.Time `gorm:"index"`
}

type Status uint8

const (
	Stable Status = 0
	Beta   Status = 1
	Alpha  Status = 2
)

func (e *Status) Scan(value interface{}) error {
	*e = Status(uint8(value.(int64)))

	return nil
}

func (e Status) Value() (uint8, error) {
	return uint8(e), nil
}
