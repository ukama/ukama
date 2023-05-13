package db

import (
	"database/sql/driver"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Org struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"uniqueIndex"`
	Owner       uuid.UUID `gorm:"type:uuid"`
	Certificate string
	Members     []*User `gorm:"many2many:org_users"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Deactivated bool
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type User struct {
	Id          uint      `gorm:"primaryKey"`
	Uuid        uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	Deactivated bool
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type OrgUser struct {
	OrgId       uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserId      uint      `gorm:"primaryKey"`
	Uuid        uuid.UUID `gorm:"not null;type:uuid"`
	Deactivated bool
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Role        RoleType       `gorm:"type:uint;not null"`
}

type RoleType uint8

const (
	Undefined RoleType = 0
	Admin     RoleType = 1
	Member    RoleType = 2
	Vendor    RoleType = 3
)

func (e *RoleType) Scan(value interface{}) error {
	*e = RoleType(uint8(value.(int64)))

	return nil
}

func (e RoleType) Value() (driver.Value, error) {
	return uint8(e), nil
}
