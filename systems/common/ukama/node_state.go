/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukama

import (
	"database/sql/driver"
	"strconv"
	"strings"
)

type NodeState uint8

const (
	NodeStateUnknown NodeState = iota
	NodeStateConfiguring
	NodeStateOperational
	NodeStateFaulty
	NodeStateInitializing
	NodeStateReady
	NodeStateUpdating
	NodeStateOffboarded
)

var nodeStateNames = map[NodeState]string{
	NodeStateUnknown:      "unknown",
	NodeStateConfiguring:  "configuring",
	NodeStateOperational:  "operational",
	NodeStateFaulty:       "faulty",
	NodeStateInitializing: "initializing",
	NodeStateReady:        "ready",
	NodeStateUpdating:     "updating",
	NodeStateOffboarded:   "offboarded",
}

var nodeStateValues = map[string]NodeState{
	"unknown":      NodeStateUnknown,
	"configuring":  NodeStateConfiguring,
	"operational":  NodeStateOperational,
	"faulty":       NodeStateFaulty,
	"initializing": NodeStateInitializing,
	"ready":        NodeStateReady,
	"updating":     NodeStateUpdating,
	"offboarded":   NodeStateOffboarded,
}

func (s *NodeState) Scan(value interface{}) error {
	*s = NodeState(uint8(value.(int64)))

	return nil
}

func (s NodeState) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s NodeState) String() string {
	v, ok := nodeStateNames[s]
	if !ok {
		return nodeStateNames[NodeStateUnknown]
	}

	return v
}

func ParseNodeState(value string) NodeState {
	i, err := strconv.Atoi(value)
	if err == nil {
		return NodeState(i)
	}

	v, ok := nodeStateValues[strings.ToLower(value)]
	if !ok {
		return NodeStateUnknown
	}

	return v
}

// Is unknown considered as valid node state?
func IsValidNodeState(value string) bool {
	_, ok := nodeStateValues[strings.ToLower(value)]

	return ok
}
