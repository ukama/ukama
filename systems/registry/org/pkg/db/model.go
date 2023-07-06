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

type OrgUser struct {
	OrgId       uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserId      uint      `gorm:"primaryKey"`
	Uuid        uuid.UUID `gorm:"not null;type:uuid"`
	Deactivated bool
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Role        RoleType       `gorm:"type:uint;not null;default:3"` // Set the default value to Member
}

type Invitation struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Org      string
	Link      string
	Email     string
	ExpiresAt time.Time
	Role 	RoleType `gorm:"type:uint;not null;default:3"` // Set the default value to Member
	Status    InvitationStatus `gorm:"type:uint;not null;default:0"` // Set the default value to Pending
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type RoleType uint8

const (
	Owner  RoleType = 0
	Admin  RoleType = 1
	Vendor RoleType = 2
	Member RoleType = 3
)

func (e *RoleType) Scan(value interface{}) error {
	*e = RoleType(uint8(value.(int64)))

	return nil
}

func (e RoleType) Value() (uint8, error) {
	return uint8(e), nil
}

type InvitationStatus uint8

const (
	Pending  InvitationStatus = 0
	Accepted InvitationStatus = 1
	Expired  InvitationStatus = 2
	Rejected InvitationStatus = 3
)

func (e *InvitationStatus) Scan(value interface{}) error {
	*e = InvitationStatus(uint8(value.(int64)))

	return nil
}

func (e InvitationStatus) Value() (uint8, error) {
	return uint8(e), nil
}
