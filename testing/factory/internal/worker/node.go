package worker

import (
	"time"

	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/testing/factory/internal"
)

func NewHNode() internal.Node {
	node := internal.Node{
		NodeID:       ukama.NewVirtualHomeNodeId(),
		Type:         ukama.NODE_ID_TYPE_HOMENODE,
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
		Type:         ukama.NODE_ID_TYPE_TOWERNODE,
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
		Type:         ukama.NODE_ID_TYPE_AMPLIFIERNODE,
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
