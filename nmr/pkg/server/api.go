package server

import (
	"github.com/ukama/openIoR/services/factory/nmr/internal/db"
)

type ReqAddOrUpdateNode struct {
	NodeID         string      `query:"nodeid" validate:"required"`
	LookingTo      string      `query:"looking_to" validate:"required"`
	Type           string      `form:"type" json:"type"`
	PartNumber     string      `form:"part_number" json:"part_number"`
	Skew           string      `form:"skew" json:"skew"`
	Mac            string      `form:"mac" json:"mac"`
	SwVersion      string      `form:"sw_version" json:"sw_version"`
	PSwVersion     string      `form:"p_sw_version" json:"p_sw_version"`
	AssmDate       string      `form:"assm_date" json:"assm_date"`
	OemName        string      `form:"oem_name" json:"oem_name"`
	Modules        []db.Module `form:"module" json:"module"`
	ProdTestStatus string      `form:"prod_test_status" json:"prod_test_status"`
	ProdReport     []byte      `form:"prod_report" json:"prod_report"` /* Report for production test */
	Status         string      `form:"status" json:"status"`
}

type ReqUpdateNodeStatus struct {
	NodeID    string `query:"nodeid" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
	Status    string `form:"status"`
}

type ReqDeleteNode struct {
	NodeID    string `query:"nodeid" validate:"required"`
	LookingTo string `query:"looking_to" validate:"required"`
}

type ReqGetNode struct {
	NodeID     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetNode struct {
	Node *db.Node `json:"node"`
}

type ReqGetNodeList struct {
	NodeID     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetNodeList struct {
	NodeList []db.Node `json:"nodes"`
}

type ReqGetNodeStatus struct {
	NodeID     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetNodeStatus struct {
	Status string `json:"node_status"`
}

type ReqUpdateNodeProdStatus struct {
	NodeID         string `query:"nodeid" validate:"required"`
	LookingTo      string `query:"looking_to" validate:"required"`
	ProdTestStatus string `form:"prod_test_status" json:"prod_test_status"`
	ProdReport     []byte `form:"prod_report" json:"prod_report"` /* Report for production test */
	Status         string `form:"status" json:"status"`
}

/* Modules */
type ReqAddOrUpdateModule struct {
	NodeID             string `query:"nodeid" validate:"required"`
	LookingTo          string `query:"looking_to" validate:"required"`
	Type               string `form:"type" json:"type"`
	PartNumber         string `form:"partnumber" json:"partnumber"`
	HwVersion          string `form:"hw_version" json:"hw_version"`
	Mac                string `form:"mac" json:"mac"`
	SwVersion          string `form:"sw_version" json:"sw_version"`
	PSwVersion         string `form:"p_sw_version" json:"p_sw_version"`
	MfgDate            string `form:"mfg_date" json:"mfg_date"`
	MfgName            string `form:"mfg_name" json:"mfg_name"`
	ProdTestStatus     string `form:"prod_test_status" json:"prod_test_status"`
	ProdReport         []byte `form:"prod_report" json:"prod_report"` /* Report for production test */
	BootstrapCerts     []byte `form:"bootstrap_certs" json:"bootstrap_certs"`
	UserCalibrartion   []byte `form:"user_calibration" json:"user_calibration"`
	FactoryCalibration []byte `form:"factory_calibration" json:"factory_calibration"`
	UserConfig         []byte `form:"user_config" json:"user_config"`
	FactoryConfig      []byte `form:"factory_config" json:"factory_config"`
	InventoryData      []byte `form:"inventory_data" json:"inventory_data"`
}

type ReqGetModule struct {
	NodeID     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespGetModule struct {
	Module *db.Module `json: "module"`
}

type ReqDeleteModule struct {
	NodeID     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type ReqGetModuleList struct {
	NodeID string `query:"nodeid" validate:"required"`
}

type RespGetModuleList struct {
	Modules []db.Module `json:"modules"`
}

type ReqUpdateModuleProdStatusData struct {
	NodeID    string `query:"nodeid" validate:"required"`
	LookingTo string `query:"looking_for" validate:"required"`
}

type RespUpdateModuleProdStatusData struct {
	ProdTestStatus     string `form:"prod_test_status" json:"prod_test_status"`
	ProdReport         []byte `form:"prod_report" json:"prod_report"` /* Report for production test */
	BootstrapCerts     []byte `form:"bootstrap_certs" json:"bootstrap_certs"`
	UserCalibrartion   []byte `form:"user_calibration" json:"user_calibration"`
	FactoryCalibration []byte `form:"factory_calibration" json:"factory_calibration"`
	UserConfig         []byte `form:"user_config" json:"user_config"`
	FactoryConfig      []byte `form:"factory_config" json:"factory_config"`
	InventoryData      []byte `form:"inventory_data" json:"inventory_data"`
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
