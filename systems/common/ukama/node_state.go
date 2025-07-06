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
	Unknown NodeState = iota
	Configured
	Operational
	Faulty
)

func (s *NodeState) Scan(value interface{}) error {
	*s = NodeState(uint8(value.(int64)))

	return nil
}

func (s NodeState) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s NodeState) String() string {
	t := map[NodeState]string{0: "unknown", 1: "configured", 2: "operational", 3: "faulty"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseNodeState(value string) NodeState {
	i, err := strconv.Atoi(value)
	if err == nil {
		return NodeState(i)
	}

	t := map[string]NodeState{"unknown": 0, "configured": 1, "operational": 2, "faulty": 3}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return NodeState(0)
	}

	return NodeState(v)
}

// Is unknown considered as valid node state?
func IsValidNodeState(value string) bool {
	t := map[string]NodeState{"unknown": 0, "configured": 1, "operational": 2, "faulty": 3}

	_, ok := t[strings.ToLower(value)]
	return ok
}
