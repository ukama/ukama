package db

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

/* Node Information */
type Node struct {
	gorm.Model
	NodeID       string    `gorm:"unique;type:string;size:23;expression:lower(node_id);size:32;not null" json:"nodeID" `
	Type         string    `gorm:"size:32;not null" json:"type"`
	PartNumber   string    `gorm:"size:32;not null" json:"partNumber"`
	Skew         string    `gorm:"size:32;not null" json:"skew"`
	Mac          string    `gorm:"size:32;not null" json:"mac"`
	SwVersion    string    `gorm:"size:32;not null" json:"swVersion"`
	PSwVersion   string    `gorm:"size:32;not null" json:"pSwVersion"`
	AssemblyDate time.Time `gorm:"type:time; size:32;not null" json:"assemblyDate"`
	OemName      string    `gorm:"size:32;not null" json:"oemName"`
	//Modules        []Module     `gorm:"ForeignKey:node;AssociationForeignKey:NodeID" json:"modules"`
	Modules       []Module `gorm:"foreignKey:UnitID;references:NodeID;" json:"modules"`
	MfgTestStatus string   `gorm:"size:32;not null" json:"mfgTestStatus"`
	MfgReport     *[]byte  `gorm:"type:bytes;" json:"mfgReport"` /* Report for production test */
	Status        string   `gorm:"size:32;not null" json:"status"`
}

/* Module Information */
type Module struct {
	gorm.Model
	ModuleID           string         `gorm:"unique;type:string;size:23;expression:lower(module_id);size:32;not null" json:"moduleID" `
	Type               string         `gorm:"size:32;not null" json:"type"`
	PartNumber         string         `gorm:"size:32;not null" json:"partNumber"`
	HwVersion          string         `gorm:"size:32;not null" json:"hwVersion"`
	Mac                string         `gorm:"size:32;not null" json:"mac"`
	SwVersion          string         `gorm:"size:32;not null" json:"swVersion"`
	PSwVersion         string         `gorm:"size:32;not null" json:"mfgSwVersion"`
	MfgDate            time.Time      `gorm:"type:time; size:32;not null" json:"mfgDate"`
	MfgName            string         `gorm:"size:32;not null" json:"mfgName"`
	MfgTestStatus      string         `gorm:"size:32;not null" json:"mfgTestStatus"`
	MfgReport          []byte         `gorm:"type:bytes;" json:"mfgReport"` /* Report for production test */
	BootstrapCerts     []byte         `gorm:"type:bytes;" json:"bootstrapCerts"`
	UserCalibration    []byte         `gorm:"type:bytes;" json:"userCalibration"`
	FactoryCalibration []byte         `gorm:"type:bytes;" json:"factoryCalibration"`
	UserConfig         []byte         `gorm:"type:bytes;" json:"userConfig"`
	FactoryConfig      []byte         `gorm:"type:bytes;" json:"factoryConfig"`
	InventoryData      []byte         `gorm:"type:bytes;" json:"inventoryData"`
	UnitID             sql.NullString `gorm:"column:unit_id;type:string;default:null;" json:"nodeId"`
}
