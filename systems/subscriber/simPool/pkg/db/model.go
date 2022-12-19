package db

import (
	"gorm.io/gorm"
)

type SimPool struct {
	gorm.Model
	network_id   string
	org_id       uint64
	iccid        string
	msisdn       string
	is_allocated bool
	sim_type     string
}
