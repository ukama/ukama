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

type NodeConnectivity uint8

const (
	// TODO: (Issue: #983) Need to add a sentinel value. And update the logic in registry/node/api-gateway where we defined hardcode value.
	// TODO: The following value should be mapped to either unknown string or undefined string, not both.
	NodeConnectivityUndefined NodeConnectivity = iota
	NodeConnectivityOnline
	NodeConnectivityOffline
)

func (s *NodeConnectivity) Scan(value interface{}) error {
	*s = NodeConnectivity(uint8(value.(int64)))

	return nil
}

func (s NodeConnectivity) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s NodeConnectivity) String() string {
	t := map[NodeConnectivity]string{0: "unknown", 1: "online", 2: "offline"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseNodeConnectivity(value string) NodeConnectivity {
	i, err := strconv.Atoi(value)
	if err == nil {
		return NodeConnectivity(i)
	}

	t := map[string]NodeConnectivity{"unknown": 0, "online": 1, "offline": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return NodeConnectivity(0)
	}

	return NodeConnectivity(v)
}
