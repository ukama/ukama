package worker

import (
	"strings"
	"time"

	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/testing/services/factory/internal"
)

func NewHNode() internal.Node {
	node := internal.Node{
		NodeID:       ukama.NewVirtualHomeNodeId(),
		Type:         strings.ToLower(ukama.NODE_ID_TYPE_HOMENODE),
		PartNumber:   "",
		Skew:         "",
		Mac:          "",
		SwVersion:    "",
		PSwVersion:   "",
		AssemblyDate: time.Now(),
		OemName:      "",
		Modules: []internal.Module{
			NewTRXModule(),
		},
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}
	return node
}

func NewTNode() internal.Node {
	node := internal.Node{
		NodeID:       ukama.NewVirtualTowerNodeId(),
		Type:         strings.ToLower(ukama.NODE_ID_TYPE_TOWERNODE),
		PartNumber:   "",
		Skew:         "",
		Mac:          "",
		SwVersion:    "",
		PSwVersion:   "",
		AssemblyDate: time.Now(),
		OemName:      "",
		Modules: []internal.Module{
			NewComModule(),
			NewTRXModule(),
		},
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}
	return node
}

func NewANode() internal.Node {
	node := internal.Node{
		NodeID:       ukama.NewVirtualAmplifierNodeId(),
		Type:         strings.ToLower(ukama.NODE_ID_TYPE_AMPNODE),
		PartNumber:   "",
		Skew:         "",
		Mac:          "",
		SwVersion:    "",
		PSwVersion:   "",
		AssemblyDate: time.Now(),
		OemName:      "",
		Modules: []internal.Module{
			NewCtrlModule(),
			NewRFModule(),
		},
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}
	return node
}
