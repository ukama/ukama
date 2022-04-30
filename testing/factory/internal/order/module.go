package order

import (
	"time"

	"github.com/ukama/ukama/servcies/common/ukama"
)

type Module struct {
	MouduleID  ukama.ModuleID
	Type       string    `form:"type" json:"type"`
	PartNumber string    `form:"partNumber" json:"partNumber"`
	HwVersion  string    `form:"hwVersion" json:"hwVersion"`
	Mac        string    `form:"mac" json:"mac"`
	SwVersion  string    `form:"swVersion" json:"swVersion"`
	PSwVersion string    `form:"pSwVersion" json:"mfgSwVersion"`
	MfgDate    time.Time `form:"mfgDate" json:"mfgDate"`
	MfgName    string    `form:"mfgName" json:"mfgName"`
	Status     string    `form:"status" json:"status"`
}

func NewComModule() Module {
	module := Module{
		MouduleID:  ukama.NewVirtualComId(),
		Type:       ukama.MODULE_ID_TYPE_COM,
		PartNumber: "",
		HwVersion:  "",
		Mac:        "",
		SwVersion:  "",
		PSwVersion: "",
		MfgDate:    time.Now(),
		MfgName:    "",
		Status:     "MfgTestStatusPending",
	}
	return module
}

func NewTRXModule() Module {
	module := Module{
		MouduleID:  ukama.NewVirtualTRXId(),
		Type:       ukama.MODULE_ID_TYPE_TRX,
		PartNumber: "",
		HwVersion:  "",
		Mac:        "",
		SwVersion:  "",
		PSwVersion: "",
		MfgDate:    time.Now(),
		MfgName:    "",
		Status:     "MfgTestStatusPending",
	}
	return module
}

func NewRFModule() Module {
	module := Module{
		MouduleID:  ukama.NewVirtualRFId(),
		Type:       ukama.MODULE_ID_TYPE_RF,
		PartNumber: "",
		HwVersion:  "",
		Mac:        "",
		SwVersion:  "",
		PSwVersion: "",
		MfgDate:    time.Now(),
		MfgName:    "",
		Status:     "MfgTestStatusPending",
	}
	return module
}

func NewCtrlModule() Module {
	module := Module{
		MouduleID:  ukama.NewVirtualCtrlId(),
		Type:       ukama.MODULE_ID_TYPE_CTRL,
		PartNumber: "",
		HwVersion:  "",
		Mac:        "",
		SwVersion:  "",
		PSwVersion: "",
		MfgDate:    time.Now(),
		MfgName:    "",
		Status:     "MfgTestStatusPending",
	}
	return module
}
