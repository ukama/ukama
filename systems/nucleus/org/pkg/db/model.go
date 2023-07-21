package db

import (
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
	Deactivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type User struct {
	Id          uint      `gorm:"primaryKey"`
	Uuid        uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	Deactivated bool
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// type OrgUser struct {
// 	OrgId       uuid.UUID `gorm:"primaryKey;type:uuid"`
// 	UserId      uint      `gorm:"primaryKey"`
// 	Uuid        uuid.UUID `gorm:"not null;type:uuid"`
// 	Deactivated bool
// 	CreatedAt   time.Time
// 	DeletedAt   gorm.DeletedAt `gorm:"index"`
// 	Role        RoleType       `gorm:"type:uint;not null;default:3"` // Set the default value to Member
// }
