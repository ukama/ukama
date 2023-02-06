package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	Uuid         uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	Name         string
	Sim_type     string
	Org_id       uuid.UUID `gorm:"not null;type:uuid"`
	Active       bool
	Duration     uint
	Sms_volume   uint
	Data_volume  uint
	Voice_volume uint
	Org_rates_id uint
}
