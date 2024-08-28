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
	"strings"
)
type NodeStateEnum uint8

const (
	StateUnknown NodeStateEnum = iota
	StateConfigure
	StateOperational
	StateFaulty
)

func (e *NodeStateEnum) Scan(value interface{}) error {
	*e = NodeStateEnum(uint8(value.(int64)))
	return nil
}

func (e NodeStateEnum) Value() (driver.Value, error) {
	return int64(e), nil
}

func (e NodeStateEnum) String() string {
	ns := map[NodeStateEnum]string{
		StateUnknown:     "unknown",
		StateConfigure:   "configure",
		StateOperational: "operational",
		StateFaulty:      "faulty",
	}
	return ns[e]
}

func ParseNodeStateEnum(s string) NodeStateEnum {
	switch strings.ToLower(s) {
	case "configure":
		return StateConfigure
	case "operational":
		return StateOperational
	case "faulty":
		return StateFaulty
	default:
		return StateUnknown
	}
}