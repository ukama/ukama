package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
)

// Site model
type Site struct {
    Id            uuid.UUID       `gorm:"primaryKey;type:uuid"`
    Name          string          
    NetworkId     uuid.UUID       `gorm:"uniqueIndex;type:uuid;index"`
    BackhaulId    uuid.UUID       `gorm:"type:uuid"`
    PowerId       uuid.UUID       `gorm:"type:uuid"`
    AccessId      uuid.UUID       `gorm:"type:uuid"`
    SwitchId      uuid.UUID       `gorm:"type:uuid"`
    IsDeactivated bool
    Latitude      float64
    Longitude     float64
    InstallDate   string
    CreatedAt     time.Time
    UpdatedAt     time.Time
    DeletedAt     gorm.DeletedAt `gorm:"index"`
}
