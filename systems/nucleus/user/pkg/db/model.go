package db

import (
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"not null;default:'unknown'"`
	Email       string    `gorm:"not null;unique"`
	Phone       string    `gorm:"unique"`
	Deactivated bool
	AuthId      uuid.UUID `gorm:"uniqueIndex:authid_unique,where:deleted_at is null;not null;type:uuid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
