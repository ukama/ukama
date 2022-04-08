package db

import (
	"github.com/jackc/pgtype"
	"github.com/ukama/openIoR/services/common/ukama"
	"gorm.io/gorm"
)

/* Node Information */
type Node struct {
	gorm.Model
	NodeID     ukama.NodeID `gorm:"type:string;primaryKey;size:23;expression:lower(node_id);size:32;not null" json:"id" `
	Type       string
	PartNumber string
	Skew       string
	Mac        string
	SwVersion  string
	PSwVersion string
	AssmDate   string
	OemName    string
	Modules    []Module
}

/* Node Production test result */
type NodeStatus struct {
	gorm.Model
	NodeID         ukama.NodeID `gorm:"type:string;primaryKey;forigenKey;size:23;expression:lower(node_id);size:32;not null" json:"id" `
	ProdTestStatus string
	ProdReport     pgtype.Bytea /* Report for production test */
	Status         string       /* Assembled, UnderTesting, Tested, Stored, Allocated, Shipped */
}

/* Module Information */
type Module struct {
	gorm.Model
	ModuleId   ukama.NodeID `gorm:"type:string;primaryKey;size:23;expression:lower(node_id);size:32;not null" json:"id" `
	Type       string
	PartNumber string
	HwVersion  string
	Mac        string
	SwVersion  string
	PSwVersion string
	MfgDate    string
	MfgName    string
}

/* Module Production result and data */
type ModuleData struct {
	gorm.Model
	ModuleId           ukama.NodeID `gorm:"type:string;primaryKey;forigenKey;size:23;expression:lower(node_id);size:32;not null" json:"id" `
	ProdTestStatus     string
	ProdReport         pgtype.Bytea /* Report for production test */
	BootstrapCerts     pgtype.Bytea
	UserCalibrartion   pgtype.Bytea
	FactoryCalibration pgtype.Bytea
	UserConfig         pgtype.Bytea
	FactoryConfig      pgtype.Bytea
	InventoryData      pgtype.Bytea
}
