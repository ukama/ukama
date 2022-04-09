package server

import (
	"time"

	"github.com/ukama/openIoR/services/factory/nmr/internal/db"
)

type ReqAddOrUpdateNode struct {
	NodeID         string    `query:"node" validate:"required"`
	LookingTo      string    `query:"looking_to" validate:"required"`
	Type           string    `form:"type" json:"type"`
	PartNumber     string    `form:"part_number" json:"partNumber"`
	Skew           string    `form:"skew" json:"skew"`
	Mac            string    `form:"mac" json:"mac"`
	SwVersion      string    `form:"sw_version" json:"swVersion"`
	PSwVersion     string    `form:"p_sw_version" json:"prodSwVersion"`
	AssemblyDate   time.Time `form:"assembly_date" json:"assemblyDate"`
	OemName        string    `form:"oem_name" json:"oemName"`
	Modules        []string  `form:"modules" json:"modules"`
	ProdTestStatus string    `form:"prod_test_status" json:"prodTestStatus"`
	ProdReport     []byte    `form:"prod_report" json:"prodReport"` /* Report for production test */
	Status         string    `form:"status" json:"status"`
}

type ReqUpdateNodeStatus struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
	Status    string `form:"status"`
}

type ReqDeleteNode struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
}

type ReqGetNode struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetNode struct {
	Node *db.Node `json:"node"`
}

type ReqGetNodeList struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetNodeList struct {
	NodeList []db.Node `json:"nodes"`
}

type ReqGetNodeStatus struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetNodeStatus struct {
	Status string `json:"node_status"`
}

type ReqUpdateNodeProdStatus struct {
	NodeID         string `query:"node" validate:"required"`
	LookingTo      string `query:"looking_to" validate:"required"`
	ProdTestStatus string `form:"prodTestStatus" json:"prodTestStatus"`
	ProdReport     []byte `form:"prodReport" json:"prodReport"` /* Report for production test */
	Status         string `form:"status" json:"status"`
}

/* Modules */
type ReqAddOrUpdateModule struct {
	ModuleID           string `query:"module" validate:"required"`
	LookingTo          string `query:"looking_to" validate:"required"`
	Type               string `form:"type" json:"type"`
	PartNumber         string `form:"partNumber" json:"partNumber"`
	HwVersion          string `form:"hwVersion" json:"hwVersion"`
	Mac                string `form:"mac" json:"mac"`
	SwVersion          string `form:"swVersion" json:"swVersion"`
	PSwVersion         string `form:"pSwVersion" json:"prodSwVersion"`
	MfgDate            string `form:"mfgDate" json:"mfgDate"`
	MfgName            string `form:"mfgName" json:"mfgName"`
	ProdTestStatus     string `form:"prodTestStatus" json:"prodTestStatus"`
	ProdReport         []byte `form:"prodTeport" json:"prodReport"` /* Report for production test */
	BootstrapCerts     []byte `form:"bootstrapCerts" json:"bootstrapCerts"`
	UserCalibrartion   []byte `form:"userCalibration" json:"userCalibration"`
	FactoryCalibration []byte `form:"factoryCalibration" json:"factoryCalibration"`
	UserConfig         []byte `form:"userConfig" json:"userConfig"`
	FactoryConfig      []byte `form:"factoryConfig" json:"factoryConfig"`
	InventoryData      []byte `form:"inventoryData" json:"inventoryData"`
}

type ReqGetModule struct {
	ModuleID   string `query:"module" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetModule struct {
	Module *db.Module `json:"module"`
}

type ReqDeleteModule struct {
	ModuleID   string `query:"module" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type ReqGetModuleList struct {
	ModuleID string `query:"module" validate:"required"`
}

type RespGetModuleList struct {
	Modules []db.Module `json:"modules"`
}

type ReqGetModuleStatusData struct {
	ModuleID  string `query:"module" validate:"required"`
	LookingTo string `query:"looking_for" validate:"required"`
}

type RespUpdateModuleStatusData struct {
	ProdTestStatus     string `form:"prodTestStatus" json:"prodTestStatus"`
	ProdReport         []byte `form:"prodReport" json:"prodReport"` /* Report for production test */
	BootstrapCerts     []byte `form:"bootstrapCerts" json:"bootstrapCerts"`
	UserCalibrartion   []byte `form:"userCalibration" json:"userCalibration"`
	FactoryCalibration []byte `form:"factoryCalibration" json:"factoryCalibration"`
	UserConfig         []byte `form:"userConfig" json:"userConfig"`
	FactoryConfig      []byte `form:"factoryConfig" json:"factoryConfig"`
	InventoryData      []byte `form:"inventoryData" json:"inventoryData"`
}

type ReqUpdateModuleStatusData struct {
	ModuleID  string `query:"module" validate:"required"`
	LookingTo string `query:"looking_for" validate:"required"`
}

type ReqUpdateModuleData struct {
	ModuleID  string `query:"module" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
	Field     string `query:"field" validate:"required"`
}

type RespUpdateModuleData struct {
	Field string `json:"field"`
	Data  []byte `json:"data"`
}

type ReqUpdateModuleBootStrapCerts struct {
	NodeID    string `query:"nodeid" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
}

type ReqUpdateModuleUserConfig struct {
	NodeID    string `query:"nodeid" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
}

type ReqUpdateModuleFactoryConfig struct {
	NodeID    string `query:"nodeid" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
}

type ReqUpdateModuleUserCalibration struct {
	NodeID    string `query:"nodeid" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
}

type ReqUpdateModuleFactoryCalibration struct {
	NodeID    string `query:"nodeid" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
}
