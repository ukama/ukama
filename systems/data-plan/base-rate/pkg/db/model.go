package db

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Rate struct {
	gorm.Model
	Uuid         uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	Country      string
	Network      string
	Vpmn         string
	Imsi         string
	Sms_mo       string
	Sms_mt       string
	Data         string
	X2g          string
	X3g          string
	X5g          string
	Lte          string
	Lte_m        string
	Apn          string
	Effective_at string
	End_at       string
	Sim_type     string
}
