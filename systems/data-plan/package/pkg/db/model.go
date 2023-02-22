package db

import (
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	Uuid        uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	Name        string
	SimType     string
	OrgId       uuid.UUID `gorm:"not null;type:uuid"`
	Active      bool
	Duration    uint
	SmsVolume   uint
	DataVolume  uint
	VoiceVolume uint
	OrgRatesId  uint
}
