package db

import (
	"gorm.io/gorm"
)

type Sim struct {
	gorm.Model
	Iccid          string `gorm:"index:idx_iccid,unique"`
	Msisdn         string
	Is_allocated   bool
	Is_failed   bool
	Sim_type       string
	SmDpAddress    string
	ActivationCode string
	QrCode         string
	Is_physical    bool
}
