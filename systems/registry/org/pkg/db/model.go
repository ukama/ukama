package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Org struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"uniqueIndex"`
	Owner       uuid.UUID `gorm:"type:uuid"`
	Certificate string
	Members     []*User `gorm:"many2many:org_members;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type User struct {
	ID   uint      `gorm:"primaryKey"`
	Uuid uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
}

type OrgUser struct {
	OrgID     uint `gorm:"primaryKey"`
	UserID    uint `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
	// Role (owner, admin, vendor)
}
