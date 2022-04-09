package db

import (
	"time"

	"github.com/ukama/openIoR/services/common/ukama"
	"gorm.io/gorm"
)

/* Node Information */
type Node struct {
	gorm.Model
	NodeID       ukama.NodeID `gorm:"type:string;primaryKey;unique;size:23;expression:lower(node_id);size:32;not null" json:"node_id" `
	Type         string       `gorm:"size:32;not null" json:"type"`
	PartNumber   string       `gorm:"size:32;not null" json:"part_number"`
	Skew         string       `gorm:"size:32;not null" json:"skew"`
	Mac          string       `gorm:"size:32;not null" json:"mac"`
	SwVersion    string       `gorm:"size:32;not null" json:"sw_version"`
	PSwVersion   string       `gorm:"size:32;not null" json:"p_sw_version"`
	AssemblyDate time.Time    `gorm:"type:time; size:32;not null" json:"assembly_date"`
	OemName      string       `gorm:"size:32;not null" json:"oem_name"`
	//Modules        []Module     `gorm:"ForeignKey:node;AssociationForeignKey:NodeID" json:"modules"`
	Modules        []Module `gorm:"ForeignKey:UnitID;references:NodeID" json:"modules"`
	ProdTestStatus string   `gorm:"size:32;not null" json:"prod_test_status"`
	ProdReport     []byte   `gorm:"type:bytes;" json:"prod_report"` /* Report for production test */
	Status         string   `gorm:"size:32;not null" json:"status"`
}

/* Module Information */
type Module struct {
	gorm.Model
	ModuleID           ukama.NodeID `gorm:"type:string;primaryKey;size:23;expression:lower(module_id);size:32;not null" json:"module_id" `
	Type               string       `gorm:"size:32;not null" json:"type"`
	PartNumber         string       `gorm:"size:32;not null" json:"partNumber"`
	HwVersion          string       `gorm:"size:32;not null" json:"hwVersion"`
	Mac                string       `gorm:"size:32;not null" json:"mac"`
	SwVersion          string       `gorm:"size:32;not null" json:"swVersion"`
	PSwVersion         string       `gorm:"size:32;not null" json:"pSwVersion"`
	MfgDate            time.Time    `gorm:"type:time; size:32;not null" json:"assmDate"`
	MfgName            string       `gorm:"size:32;not null" json:"mfgName"`
	ProdTestStatus     string       `gorm:"size:32;not null" json:"prodTestStatus"`
	ProdReport         []byte       `gorm:"type:bytes;" json:"ProdReport"` /* Report for production test */
	BootstrapCerts     []byte       `gorm:"type:bytes;" json:"BootstrapCerts"`
	UserCalibrartion   []byte       `gorm:"type:bytes;" json:"UserCalibrartion"`
	FactoryCalibration []byte       `gorm:"type:bytes;" json:"FactoryCalibration"`
	UserConfig         []byte       `gorm:"type:bytes;" json:"UserConfig"`
	FactoryConfig      []byte       `gorm:"type:bytes;" json:"FactoryConfig"`
	InventoryData      []byte       `gorm:"type:bytes;" json:"InventoryData"`
	UnitID             ukama.NodeID `gorm:"type:string;" json:"NodeId"`
}
