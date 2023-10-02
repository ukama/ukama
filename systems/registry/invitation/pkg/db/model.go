package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Invitation struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Org       string
	Link      string
	Email     string
	Name      string
	ExpiresAt time.Time
	Role      RoleType         `gorm:"type:uint;not null;default:3"` // Set the default value to Member
	Status    InvitationStatus `gorm:"type:uint;not null;default:0"` // Set the default value to Pending
	UserId    string           `gorm:"type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type InvitationStatus uint8

const (
	Pending  InvitationStatus = 0
	Accepted InvitationStatus = 1
	Declined InvitationStatus = 2
)

func (e *InvitationStatus) Scan(value interface{}) error {
	*e = InvitationStatus(uint8(value.(int64)))

	return nil
}

func (e InvitationStatus) Value() (uint8, error) {
	return uint8(e), nil
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
