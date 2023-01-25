package db

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Org struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"uniqueIndex"`
	Deactivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Network struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"uniqueIndex:network_name_org_idx"`
	OrgID       uuid.UUID `gorm:"uniqueIndex:network_name_org_idx;type:uuid"`
	Org         *Org
	Deactivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Site struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"uniqueIndex:site_name_network_idx"`
	NetworkID   uuid.UUID `gorm:"uniqueIndex:site_name_network_idx;type:uuid"`
	Network     *Network
	Deactivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
