package order

import (
	"time"

	"github.com/ukama/ukama/servcies/common/ukama"
)

type Node struct {
	NodeID        ukama.NodeID `json:-`
	Type          string       `json:"type"`
	PartNumber    string       `json:"partNumber"`
	Skew          string       `json:"skew"`
	Mac           string       `json:"mac"`
	SwVersion     string       `json:"swVersion"`
	PSwVersion    string       `json:"mfgSwVersion"`
	AssemblyDate  time.Time    `json:"assemblyDate"`
	OemName       string       `json:"oemName"`
	Modules       []Module     `json:-`
	MfgTestStatus string       `json:"mfgTestStatus"`
	Status        string       `json:"status"`
}

func NewHNode() Node {
	node := Node{
		NodeID:       ukama.NewVirtualHomeNodeId(),
		Type:         ukama.NODE_ID_TYPE_HOMENODE,
		PartNumber:   "",
		Skew:         "",
		Mac:          "",
		SwVersion:    "",
		PSwVersion:   "",
		AssemblyDate: time.Now(),
		OemName:      "",
		Modules: []Module{
			NewTRXModule(),
		},
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}
	return node
}

func NewTNode() Node {
	node := Node{
		NodeID:       ukama.NewVirtualTowerNodeId(),
		Type:         ukama.NODE_ID_TYPE_TOWERNODE,
		PartNumber:   "",
		Skew:         "",
		Mac:          "",
		SwVersion:    "",
		PSwVersion:   "",
		AssemblyDate: time.Now(),
		OemName:      "",
		Modules: []Module{
			NewComModule(),
			NewTRXModule(),
		},
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}
	return node
}

func NewANode() Node {
	node := Node{
		NodeID:       ukama.NewVirtualAmplifierNodeId(),
		Type:         ukama.NODE_ID_TYPE_AMPLIFIERNODE,
		PartNumber:   "",
		Skew:         "",
		Mac:          "",
		SwVersion:    "",
		PSwVersion:   "",
		AssemblyDate: time.Now(),
		OemName:      "",
		Modules: []Module{
			NewCtrlModule(),
			NewRFModule(),
		},
		MfgTestStatus: "MfgTestStatusPending",
		Status:        "StatusLabelGenerated",
	}
	return node
}
