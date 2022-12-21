package db

import (
	"gorm.io/gorm"
)

type SimPool struct {
	gorm.Model
	Iccid          string
	Msisdn         string
	Is_allocated   bool
	Sim_type       string
	SmDpAddress    string
	ActivationCode string
	QrCode         string
}
