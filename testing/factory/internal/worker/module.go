package worker

import (
	"time"

	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/testing/factory/internal"
)

func NewComModule() internal.Module {
	module := internal.Module{
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

func NewTRXModule() internal.Module {
	module := internal.Module{
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

func NewRFModule() internal.Module {
	module := internal.Module{
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

func NewCtrlModule() internal.Module {
	module := internal.Module{
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
