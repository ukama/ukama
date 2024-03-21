package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/lib/pq"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
)

// Site model
type Site struct {
    ID            uuid.UUID       `gorm:"primaryKey;type:uuid"`
    Name           string          `gorm:"uniqueIndex:site_name_network_idx"`
    NetworkID      uuid.UUID       `gorm:"uniqueIndex:site_name_network_idx;type:uuid"`
    BackhaulID     uuid.UUID      `gorm:"type:uuid"` 
    PowerID       uuid.UUID        `gorm:"type:uuid"` 
    AccessID      uuid.UUID        `gorm:"type:uuid"` 
    SwitchID      uuid.UUID       `gorm:"type:uuid"` 
    IsDeactivated  bool            
    Latitude       float64         
    Longitude      float64         
    InstallDate    string       
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      gorm.DeletedAt `gorm:"index"`
	Network       *Network         
}

// Network model
type Network struct {
	ID               uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name             string    `gorm:"uniqueIndex:network_name_org_idx"`
	OrgID            uuid.UUID `gorm:"uniqueIndex:network_name_org_idx;type:uuid"`
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
	SyncStatus       ukama.StatusType
}
