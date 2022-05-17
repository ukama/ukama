package worker

import (
	"strings"
	"time"

	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/testing/services/factory/internal"
)

func NewComModule() internal.Module {
	module := internal.Module{
		ModuleID:   ukama.NewVirtualComId(),
		Type:       strings.ToLower(ukama.MODULE_ID_TYPE_COMP),
		PartNumber: "",
		HwVersion:  "",
		Mac:        "",
		SwVersion:  "",
		PSwVersion: "",
		MfgDate:    time.Now(),
		MfgName:    "",
		Status:     "StatusAssemblyCompleted",
	}
	return module
}

func NewTRXModule() internal.Module {
	module := internal.Module{
		ModuleID:   ukama.NewVirtualTRXId(),
		Type:       strings.ToLower(ukama.MODULE_ID_TYPE_TRX),
		PartNumber: "",
		HwVersion:  "",
		Mac:        "",
		SwVersion:  "",
		PSwVersion: "",
		MfgDate:    time.Now(),
		MfgName:    "",
		Status:     "StatusAssemblyCompleted",
	}
	return module
}

func NewRFModule() internal.Module {
	module := internal.Module{
		ModuleID:   ukama.NewVirtualFEId(),
		Type:       strings.ToLower(ukama.MODULE_ID_TYPE_FE),
		PartNumber: "",
		HwVersion:  "",
		Mac:        "",
		SwVersion:  "",
		PSwVersion: "",
		MfgDate:    time.Now(),
		MfgName:    "",
		Status:     "StatusAssemblyCompleted",
	}
	return module
}

func NewCtrlModule() internal.Module {
	module := internal.Module{
		ModuleID:   ukama.NewVirtualCtrlId(),
		Type:       strings.ToLower(ukama.MODULE_ID_TYPE_CTRL),
		PartNumber: "",
		HwVersion:  "",
		Mac:        "",
		SwVersion:  "",
		PSwVersion: "",
		MfgDate:    time.Now(),
		MfgName:    "",
		Status:     "StatusAssemblyCompleted",
	}
	return module
}
