package worker

import (
	"time"

	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/testing/factory/internal"
)

func NewComModule() internal.Module {
	module := internal.Module{
		ModuleID:   ukama.NewVirtualComId(),
		Type:       ukama.MODULE_ID_TYPE_COMP,
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
		ModuleID:   ukama.NewVirtualTRXId(),
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
		ModuleID:   ukama.NewVirtualFEId(),
		Type:       ukama.MODULE_ID_TYPE_FE,
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
		ModuleID:   ukama.NewVirtualCtrlId(),
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
