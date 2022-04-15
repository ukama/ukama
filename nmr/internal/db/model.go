package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MfgStatus string

const (
	StatusUnkown                      MfgStatus = ""
	StatusLabelGenrated               MfgStatus = "StatusLabelGenrated"
	StatusModuleTest                  MfgStatus = "StatusModuleTest"
	StatusModuleTestPass              MfgStatus = "StatusModuleTestPass"
	StatusModuleTestFailure           MfgStatus = "StatusModuleTestFailure"
	StatusModuleNeedAssistance        MfgStatus = "StatusModuleNeedAssistance"
	StatusModuleWaitngForShipment     MfgStatus = "StatusModuleWaitngForShipment"
	StatusModuleInTransit             MfgStatus = "StatusModuleInTransit"
	StatusModuleReadyForAssembly      MfgStatus = "StatusModuleReadyForAssembly"
	StatusPendingAssembly             MfgStatus = "StatusPendingAssembly"
	StatusUnderAssembly               MfgStatus = "StatusUnderAssembly"
	StatusAssemblyCompleted           MfgStatus = "StatusAssemblyCompleted"
	StatusPendingProductionTest       MfgStatus = "StatusPendingProductionTest"
	StatusProductionTest              MfgStatus = "StatusProductionTest"
	StatusProductionTestPass          MfgStatus = "StatusProductionTestPass"
	StatusProductionTestFail          MfgStatus = "StatusProductionTestFail"
	StatusProductionTestNeedAssitance MfgStatus = "StatusProductionTestNeedAssitance"
	StatusProductionTestCompleted     MfgStatus = "StatusProductionTestCompleted"
	StatusIntransitToWareHouse        MfgStatus = "StatusIntransitToWareHouse"
	StatusNodeAllocated               MfgStatus = "StatusNodeAllocated"
	StatusNodeWaitngForShipment       MfgStatus = "StatusNodeWaitngForShipment"
	StatusNodeIntransit               MfgStatus = "StateNodeIntransit"
)

type MfgTestStatus string

const (
	MfgTestStatusPending   MfgTestStatus = "MfgTestStatusPending"
	MfgTestStatusUnderTest MfgTestStatus = "MfgTestStatusUnderTest"
	MfgTestStatusPass      MfgTestStatus = "MfgTestStatusPass"
	MfgTestStatusFail      MfgTestStatus = "MfgTestStatusFail"
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
	ModuleID           string    `gorm:"unique;type:string;size:23;expression:lower(module_id);size:32;not null" json:"moduleID" `
	Type               string    `gorm:"size:32;not null" json:"type"`
	PartNumber         string    `gorm:"size:32;not null" json:"partNumber"`
	HwVersion          string    `gorm:"size:32;not null" json:"hwVersion"`
	Mac                string    `gorm:"size:32;not null" json:"mac"`
	SwVersion          string    `gorm:"size:32;not null" json:"swVersion"`
	PSwVersion         string    `gorm:"size:32;not null" json:"mfgSwVersion"`
	MfgDate            time.Time `gorm:"type:time; size:32;not null" json:"mfgDate"`
	MfgName            string    `gorm:"size:32;not null" json:"mfgName"`
	Status             string    `gorm:"size:32;not null" json:"Status"`
	MfgReport          *[]byte   `gorm:"type:bytes;" json:"mfgReport,omitempty"` /* Report for production test */
	BootstrapCerts     *[]byte   `gorm:"type:bytes;" json:"bootstrapCerts,omitempty"`
	UserCalibration    *[]byte   `gorm:"type:bytes;" json:"userCalibration,omitempty"`
	FactoryCalibration *[]byte   `gorm:"type:bytes;" json:"factoryCalibration,omitempty"`
	UserConfig         *[]byte   `gorm:"type:bytes;" json:"userConfig,omitempty"`
	FactoryConfig      *[]byte   `gorm:"type:bytes;" json:"factoryConfig,omitempty"`
	InventoryData      *[]byte   `gorm:"type:bytes;" json:"inventoryData,omitempty"`
	//	CloudCerts         *[]byte        `gorm:"type:bytes;" json:"cloudCerts"`
	UnitID sql.NullString `gorm:"column:unit_id;type:string;default:null;" json:"nodeId"`
}

func (m MfgTestStatus) String() string {
	return string(m)
}

func (s MfgStatus) String() string {
	return string(s)
}

func MfgTestState(s string) (*MfgTestStatus, error) {
	mfgTestStatus := map[MfgTestStatus]struct{}{
		MfgTestStatusPending:   {},
		MfgTestStatusUnderTest: {},
		MfgTestStatusPass:      {},
		MfgTestStatusFail:      {},
	}

	status := MfgTestStatus(s)

	_, ok := mfgTestStatus[status]
	if !ok {
		return nil, fmt.Errorf("%s is invalid mfg test status", s)
	}

	return &status, nil

}

func MfgState(s string) (*MfgStatus, error) {

	mfgStatus := map[MfgStatus]struct{}{
		StatusLabelGenrated:               {},
		StatusModuleTest:                  {},
		StatusModuleTestPass:              {},
		StatusModuleTestFailure:           {},
		StatusModuleNeedAssistance:        {},
		StatusModuleWaitngForShipment:     {},
		StatusModuleInTransit:             {},
		StatusModuleReadyForAssembly:      {},
		StatusPendingAssembly:             {},
		StatusUnderAssembly:               {},
		StatusAssemblyCompleted:           {},
		StatusPendingProductionTest:       {},
		StatusProductionTest:              {},
		StatusProductionTestPass:          {},
		StatusProductionTestFail:          {},
		StatusProductionTestNeedAssitance: {},
		StatusProductionTestCompleted:     {},
		StatusIntransitToWareHouse:        {},
		StatusNodeAllocated:               {},
		StatusNodeWaitngForShipment:       {},
		StatusNodeIntransit:               {},
	}

	status := MfgStatus(s)

	_, ok := mfgStatus[status]
	if !ok {
		return nil, fmt.Errorf("%s is invalid mfg test status", s)
	}

	return &status, nil
}
