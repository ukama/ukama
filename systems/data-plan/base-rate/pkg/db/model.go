package db

import (
	"time"

	"gorm.io/gorm"
)

type Rate struct {
	gorm.Model
	Id           uint   `gorm:"primaryKey,where:deleted_at is null;size:23"`
	Country      string `gorm:"type :string"`
	Network      string `gorm:"type:string"`
	Vpmn         string `gorm:"type:string"`
	Imsi         string `gorm:"type:string"`
	Sms_mo       string `gorm:"type:string"`
	Sms_mt       string `gorm:"type:string"`
	Data         string `gorm:"type:string"`
	X2g          string `gorm:"type:string"`
	X3g          string `gorm:"type:string"`
	X5g          string `gorm:"type:string"`
	Lte          string `gorm:"type:string"`
	Lte_m        string `gorm:"type:string"`
	Apn          string `gorm:"type:string"`
	// Created_at   time.Time  `gorm:"autoCreateTime:true"`
	Effective_at time.Time `gorm:"autoCreateTime:false"`
	End_at       time.Time `gorm:"autoCreateTime:false"`
	// Deleted_at time.Time `gorm:"autoCreateTime:false;null"`
	Sim_type     string `gorm:"type:string"`
}
