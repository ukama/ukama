/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import "time"

const (
	SystemName  = "operation"
	ServiceName = "manager"
)

const (
	DefaultLeaseTTL = 5 * time.Minute
	SweeperInterval = 30 * time.Second
)

var (
	IsDebugMode bool   = false
	InstanceId  string = ServiceName + "-debug"
)
