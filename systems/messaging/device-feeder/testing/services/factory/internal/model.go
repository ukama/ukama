package internal

import (
	"time"

	"github.com/ukama/ukama/systems/common/ukama"
)

type Module struct {
	ModuleID   ukama.ModuleID `json:"moduleID"`
	Type       string         `form:"type" json:"type"`
	PartNumber string         `form:"partNumber" json:"partNumber"`
	HwVersion  string         `form:"hwVersion" json:"hwVersion"`
	Mac        string         `form:"mac" json:"mac"`
	SwVersion  string         `form:"swVersion" json:"swVersion"`
	PSwVersion string         `form:"pSwVersion" json:"mfgSwVersion"`
	MfgDate    time.Time      `form:"mfgDate" json:"mfgDate"`
	MfgName    string         `form:"mfgName" json:"mfgName"`
	Status     string         `form:"status" json:"status"`
}

type Node struct {
	NodeID        ukama.NodeID `json:"-"`
	Type          string       `json:"type"`
	PartNumber    string       `json:"partNumber"`
	Skew          string       `json:"skew"`
	Mac           string       `json:"mac"`
	SwVersion     string       `json:"swVersion"`
	PSwVersion    string       `json:"mfgSwVersion"`
	AssemblyDate  time.Time    `json:"assemblyDate"`
	OemName       string       `json:"oem"`
	Modules       []Module     `json:"-"`
	MfgTestStatus string       `json:"mfgTestStatus"`
	Status        string       `json:"status"`
}

type ReqNodeStausUpdate struct {
	Id     ukama.NodeID `gorm:"type:string;primaryKey;size:23" json:"id"`
	Name   string       `gorm:"size:255;" json:"name"`
	Type   string       `gorm:"size:255;not null" json:"type"`
	Status string       `gorm:"size:255;not null" json:"status"`
}
