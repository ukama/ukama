package db

import (
	uuid "github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	PackageID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name         string
	SimType     string
	OrgID       uuid.UUID `gorm:"not null;type:uuid;index"`
	Active       bool
	Duration     uint
	SmsVolume   uint
	DataVolume  uint
	VoiceVolume uint
	OrgRatesID uint
}
