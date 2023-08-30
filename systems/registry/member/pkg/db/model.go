package db

import (
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	UserId      uuid.UUID `gorm:"uniqueIndex:user_id_idx,where:deleted_at is null;not null;type:uuid"`
	Deactivated bool      `gorm:"default:false"`
	Role        RoleType  `gorm:"type:uint;not null;default:4"` // Set the default value to Member
}

type RoleType uint8

const (
	Owner    RoleType = 0
	Admin    RoleType = 1
	Employee RoleType = 2
	Vendor   RoleType = 3
	Users    RoleType = 4
)

func (e *RoleType) Scan(value interface{}) error {
	*e = RoleType(uint8(value.(int64)))

	return nil
}

func (e RoleType) Value() (uint8, error) {
	return uint8(e), nil
}
