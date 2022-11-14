package db

import (
	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	Name         string
	Sim_type     string
	Org_id       uint
	Active       bool
	Duration     uint
	Sms_volume   uint
	Data_volume  uint
	Voice_volume uint
	Org_rates_id uint
}
