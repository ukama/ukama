package db

import (
	"gorm.io/gorm"
)

type Sim struct {
	gorm.Model
	Iccid          string `gorm:"index:idx_iccid,unique"`
	Msisdn         string
	IsAllocated    bool
	IsFailed       bool
	SimType        string
	SmDpAddress    string
	ActivationCode string
	QrCode         string
	IsPhysical     bool
}
