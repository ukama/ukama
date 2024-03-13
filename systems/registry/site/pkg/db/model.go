package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
)

//Site model
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
    InstallDate    time.Time       
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      gorm.DeletedAt `gorm:"index"`
}
