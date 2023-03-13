package db

import (
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uuid        uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	Name        string    `gorm:"not null;default:'unknown'"`
	Email       string    `gorm:"not null;unique"`
	Phone       string    `gorm:"not null;unique"`
	Deactivated bool
}
