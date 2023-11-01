package db

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/types"
	"github.com/ukama/ukama/systems/common/uuid"
)

type Org struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"uniqueIndex"`
	Deactivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Network struct {
	Id               uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name             string    `gorm:"uniqueIndex:network_name_org_idx"`
	OrgId            uuid.UUID `gorm:"uniqueIndex:network_name_org_idx;type:uuid"`
	Org              *Org
	Deactivated      bool
	AllowedCountries pq.StringArray `gorm:"type:varchar(64)[]" json:"allowed_countries"`
	AllowedNetworks  pq.StringArray `gorm:"type:varchar(64)[]" json:"allowed_networks"`
	Budget           float64
	Overdraft        float64
	TrafficPolicy    uint32
	PaymentLinks     bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	SyncStatus       types.SyncStatus
}

type Site struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"uniqueIndex:site_name_network_idx"`
	NetworkId   uuid.UUID `gorm:"uniqueIndex:site_name_network_idx;type:uuid"`
	Network     *Network
	Deactivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
