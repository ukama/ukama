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
	Undefined NodeState = iota
	Onboarded
	Configured
	Active
	Maintenance
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
	t := map[NodeState]string{0: "undefined", 1: "onboarded", 2: "configured", 3: "active", 4: "maintenance", 5: "faulty"}

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

	t := map[string]CdrType{"undefined": 0, "onboarded": 1, "configured": 2, "active": 3, "maintenance": 4, "faulty": 5}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return NodeState(0)
	}

	return NodeState(v)
}
