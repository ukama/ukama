package db

import (
	"time"

	"github.com/jackc/pgtype"
	"github.com/ukama/openIoR/services/common/ukama"
	"gorm.io/gorm"
)

/* Node Information */
type Node struct {
	gorm.Model
	NodeID         ukama.NodeID   `gorm:"type:string;primaryKey;size:23;expression:lower(node_id);size:32;not null" json:"nodeId" `
	Type           string         `gorm:"size:32;not null" json:"type"`
	PartNumber     string         `gorm:"size:32;not null" json:"partNumber"`
	Skew           string         `gorm:"size:32;not null" json:"skew"`
	Mac            string         `gorm:"size:32;not null" json:"mac"`
	SwVersion      string         `gorm:"size:32;not null" json:"swVersion"`
	PSwVersion     string         `gorm:"size:32;not null" json:"pSwVersion"`
	AssmDate       time.Time      `gorm:"type:pgtype.Time;size:32;not null" json:"assmDate"`
	OemName        string         `gorm:"size:32;not null" json:"oemName"`
	Modules        []ukama.NodeID `gorm:"type:string;foreignkey;expression:lower(node_id);;not null" json:"modules" `
	ProdTestStatus string         `gorm:"size:32;not null" json:"prodTestStatus"`
	ProdReport     []byte         `gorm:"type:pgtype.Bytea;" json:"ProdReport"` /* Report for production test */
	Status         string         `gorm:"size:32;not null" json:"status"`
}

/* Module Information */
type Module struct {
	gorm.Model
	ModuleId           ukama.NodeID `gorm:"type:string;primaryKey;size:23;expression:lower(node_id);size:32;not null" json:"moduleId" `
	Type               string       `gorm:"size:32;not null" json:"type"`
	PartNumber         string       `gorm:"size:32;not null" json:"partNumber"`
	HwVersion          string       `gorm:"size:32;not null" json:"hwVersion"`
	Mac                string       `gorm:"size:32;not null" json:"mac"`
	SwVersion          string       `gorm:"size:32;not null" json:"swVersion"`
	PSwVersion         string       `gorm:"size:32;not null" json:"pSwVersion"`
	MfgDate            time.Time    `gorm:"type:pgtype.Time size:32;not null" json:"assmDate"`
	MfgName            string       `gorm:"size:32;not null" json:"mfgName"`
	ProdTestStatus     string       `gorm:"size:32;not null" json:"prodTestStatus"`
	ProdReport         pgtype.Bytea `gorm:"type:pgtype.Bytea;" json:"ProdReport"` /* Report for production test */
	BootstrapCerts     pgtype.Bytea `gorm:"type:pgtype.Bytea;" json:"BootstrapCerts"`
	UserCalibrartion   pgtype.Bytea `gorm:"type:pgtype.Bytea;" json:"UserCalibrartion"`
	FactoryCalibration pgtype.Bytea `gorm:"type:pgtype.Bytea;" json:"FactoryCalibration"`
	UserConfig         pgtype.Bytea `gorm:"type:pgtype.Bytea;" json:"UserConfig"`
	FactoryConfig      pgtype.Bytea `gorm:"type:pgtype.Bytea;" json:"FactoryConfig"`
	InventoryData      pgtype.Bytea `gorm:"type:pgtype.Bytea;" json:"InventoryData"`
}
