/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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
