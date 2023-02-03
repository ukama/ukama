package db

import (
	"gorm.io/gorm"
)

type Sim struct {
	gorm.Model
	Iccid          string `gorm:"index:idx_iccid,unique"`
	Msisdn         string
<<<<<<< HEAD
	IsAllocated    bool
	IsFailed       bool
	SimType        string
	SmDpAddress    string
	ActivationCode string
	QrCode         string
	IsPhysical     bool
=======
	Is_allocated   bool
	Sim_type       string
	SmDpAddress    string
	ActivationCode string
	QrCode         string
	Is_physical    bool
>>>>>>> subscriber-sys_sim-manager
}
