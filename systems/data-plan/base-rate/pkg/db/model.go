package db

import (
	"time"

	"gorm.io/gorm"
)

type Rate struct {
	gorm.Model
	Id           uint      `gorm:"primaryKey,where:deleted_at is null;size:23"`
	Country      string    `gorm:"type :string"`
	Network      string    `gorm:"type:string"`
	Vpmn         string    `gorm:"type:string"`
	Imsi         string    `gorm:"type:string"`
	Sms_mo       string    `gorm:"type:string"`
	Sms_mt       string    `gorm:"type:string"`
	Data         string    `gorm:"type:string"`
	X2g          string    `gorm:"type:string"`
	X3g          string    `gorm:"type:string"`
	X5g          string    `gorm:"type:string"`
	Lte          string    `gorm:"type:string"`
	Lte_m        string    `gorm:"type:string"`
	Apn          string    `gorm:"type:string"`
	Created_at   time.Time `gorm:"default:now();autoCreateTime:mili;not null"`
	Deleted_at   time.Time `gorm:"type:TIMESTAMP WITH TIME ZONE ;default:null"`
	Updated_at   time.Time `gorm:"default:now();autoUpdateTime:milli;not null"`
	Effective_at time.Time `gorm:";type:TIMESTAMP WITH TIME ZONE ;not null"`
	End_at       time.Time `gorm:"type:TIMESTAMP WITH TIME ZONE ;default:null"`
	Sim_type     string    `gorm:"type:string"`
}
