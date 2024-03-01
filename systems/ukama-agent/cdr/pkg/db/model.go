package db

import (
	"gorm.io/gorm"
)

type CDR struct {
	gorm.Model
	Session       uint64
	NodeId        string `gorm:"Index:cdr_node_idx;not null"`
	Imsi          string `gorm:"Index:cdr_imsi_idx;not null"`
	Policy        string `gorm:"Index:cdr_policy_idx;not null"`
	ApnName       string
	Ip            string
	StartTime     uint64
	EndTime       uint64
	LastUpdatedAt uint64
	TxBytes       uint64
	RxBytes       uint64
	TotalBytes    uint64
}

type Usage struct {
	gorm.Model
	Imsi  string `gorm:"Index:cdr_imsi_idx;not null"`
	Usage uint64
}
