package db

import (
	"errors"
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
	Role        string         `gorm:"type:enum('owner', 'admin', 'member','vendor')"`
}

func (u *OrgUser) Validate(db *gorm.DB) {
	if !containsRole(u.Role) {
		db.AddError(errors.New("invalid role"))
	}
}

func containsRole(role string) bool {
	roles := []string{"owner", "admin", "member", "vendor"}
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
