package server

import (
	"encoding/json"
	"time"

	"github.com/ukama/ukama/services/factory/nmr/internal/db"
)

type ReqAddOrUpdateNode struct {
	NodeID        string    `query:"node" validate:"required"`
	LookingTo     string    `query:"looking_to" validate:"required"`
	Type          string    `form:"type" json:"type"`
	PartNumber    string    `form:"part_number" json:"partNumber"`
	Skew          string    `form:"skew" json:"skew"`
	Mac           string    `form:"mac" json:"mac"`
	SwVersion     string    `form:"sw_version" json:"swVersion"`
	PSwVersion    string    `form:"p_sw_version" json:"mfgSwVersion"`
	AssemblyDate  time.Time `form:"assembly_date" json:"assemblyDate"`
	OemName       string    `form:"oem_name" json:"oemName"`
	Modules       []string  `form:"modules" json:"modules"`
	MfgTestStatus string    `form:"mfg_test_status" json:"mfgTestStatus"`
	Status        string    `form:"status" json:"status"`
}

type ReqUpdateNodeStatus struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
	Status    string `query:"status" validate:"required"`
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

type ReqGetNodeMfgStatus struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetNodeMfgStatus struct {
	Status     string           `json:"node_status"`
	ProdReport *json.RawMessage `form:"mfgReport" json:"mfgReport"`
}

type ReqUpdateNodeMfgStatus struct {
	NodeID        string           `query:"node" validate:"required"`
	LookingTo     string           `query:"looking_to" validate:"required"`
	MfgTestStatus string           `form:"mfgTestStatus" json:"mfgTestStatus"`
	MfgReport     *json.RawMessage `form:"mfgReport" json:"mfgReport"` /* Report for mfg function test */
	Status        string           `form:"status" json:"status"`
}

/* Modules */
type ReqAddOrUpdateModule struct {
	ModuleID   string    `query:"module" validate:"required"`
	LookingTo  string    `query:"looking_to" validate:"required"`
	Type       string    `form:"type" json:"type"`
	PartNumber string    `form:"partNumber" json:"partNumber"`
	HwVersion  string    `form:"hwVersion" json:"hwVersion"`
	Mac        string    `form:"mac" json:"mac"`
	SwVersion  string    `form:"swVersion" json:"swVersion"`
	PSwVersion string    `form:"pSwVersion" json:"mfgSwVersion"`
	MfgDate    time.Time `form:"mfgDate" json:"mfgDate"`
	MfgName    string    `form:"mfgName" json:"mfgName"`
	Status     string    `form:"status" json:"status"`
	UnitID     string    `form:"nodeID" json:"nodeID,omitempty"`
}

type ReqGetModule struct {
	ModuleID   string `query:"module" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetModule struct {
	Module *db.Module `json:"module"`
}

type ReqDeleteModule struct {
	ModuleID  string `query:"module" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
}

type ReqGetModuleList struct {
	ModuleID string `query:"module" validate:"required"`
}

type RespGetModuleList struct {
	Modules []db.Module `json:"modules"`
}

type ReqGetModuleMfgStatusData struct {
	ModuleID   string `query:"module" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetModuleMfgStatusData struct {
	Status string `form:"status" json:"status"`
}

type ReqUpdateModuleMfgStatusData struct {
	ModuleID  string `query:"module" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
	Status    string `query:"status" json:"status"`
}

type ReqGetModuleMfgData struct {
	ModuleID   string `query:"module" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetModuleMfgData struct {
	Status             string          `form:"status" json:"status"`
	MfgReport          json.RawMessage `form:"mfgReport" json:"mfgReport,omitempty"` /* Report for mfguction test */
	BootstrapCerts     json.RawMessage `form:"bootstrapCerts" json:"bootstrapCerts,omitempty"`
	UserCalibration    json.RawMessage `form:"userCalibration" json:"userCalibration,omitempty"`
	FactoryCalibration json.RawMessage `form:"factoryCalibration" json:"factoryCalibration,omitempty"`
	UserConfig         json.RawMessage `form:"userConfig" json:"userConfig,omitempty"`
	FactoryConfig      json.RawMessage `form:"factoryConfig" json:"factoryConfig,omitempty"`
	InventoryData      json.RawMessage `form:"inventoryData" json:"inventoryData,omitempty"`
}

// type ReqUpdateModuleStatusData struct {
// 	ModuleID  string `query:"module" validate:"required"`
// 	LookingTo string `query:"looking_for" validate:"required"`
// 	Status    string `query:"status" validate:"required"`
// }

type ReqGetModuleMfgField struct {
	ModuleID   string `query:"module" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
	Field      string `query:"field" validate:"required"`
}

type RespGetModuleMfgField struct {
	Field string `json:"field"`
	Data  []byte `json:"data"`
}

type ReqUpdateModuleMfgField struct {
	ModuleID  string `form:"module" query:"module" json:"module" validate:"required"`
	LookingTo string `form:"looking_to" query:"looking_to" json:"looking_to" validate:"required"`
	Field     string `form:"field" query:"field" json:"field" validate:"required"`
}

type ReqAssignModuleToNode struct {
	NodeID    string `query:"node" validate:"required"`
	ModuleID  string `query:"module" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
}

type ReqUpdateModuleBootStrapCerts struct {
	ModuleID  string `query:"module" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
}
