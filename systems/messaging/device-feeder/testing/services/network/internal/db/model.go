package db

import (
	"fmt"

	"gorm.io/gorm"
)

type VNodeStatus string

const (
	VNodePreCheck VNodeStatus = "PreCheck"
	VNodeOn       VNodeStatus = "PowerOn"
	VNodeOff      VNodeStatus = "PowerOff"
)

type VNode struct {
	gorm.Model
	NodeID string `gorm:"unique;type:string;size:23;expression:lower(node_id);size:32;not null" json:"nodeID"`
	Status string `gorm:"size:32;not null" json:"status"`
}

func (s VNodeStatus) String() string {
	return string(s)
}

func VNodeState(s string) (*VNodeStatus, error) {
	vNodeStatus := map[VNodeStatus]struct{}{
		VNodePreCheck: {},
		VNodeOn:       {},
		VNodeOff:      {},
	}

	status := VNodeStatus(s)

	_, ok := vNodeStatus[status]
	if !ok {
		return nil, fmt.Errorf("%s is invalid virtual node status", s)
	}

	return &status, nil

}
