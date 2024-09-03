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
 type Connectivity uint8

const (
	Unknown Connectivity = iota
	Offline Connectivity = 1 /* Not connected */
	Online  Connectivity = 2 /* Connected */
)
func (c *Connectivity) Scan(value interface{}) error {
	*c = Connectivity(uint8(value.(int64)))

	return nil
}

func (c Connectivity) Value() (driver.Value, error) {
	return int64(c), nil
}
func (c Connectivity) String() string {
	cs := map[Connectivity]string{
		Unknown: "unkown",
		Offline: "offline",
		Online:  "online",
	}

	return cs[c]
}

func ParseConnectivityState(s string) Connectivity {
	switch strings.ToLower(s) {
	case "offline":
		return Offline
	case "online":
		return Online
	default:
		return Unknown
	}
}
