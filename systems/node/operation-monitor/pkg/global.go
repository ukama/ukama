/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import "time"

const (
	SystemName  = "node"
	ServiceName = "operationmonitor"
)

const (
	DefaultDeadlineTTL = 4 * time.Minute
	SweeperInterval    = 30 * time.Second
)

// Action → completion rule fallback when caller doesn't supply one.
// A reboot cycles the node substate on→reboot→on without changing the main
// state (Configured stays Configured), so completion is the node reporting
// substate=on again — state=Operational only ever fires after a "ready"
// event and never on reboot, which left locks held until the deadline.
// TODO: move to config or per-action proto when we add more actions.
var DefaultCompletionRule = map[string]string{
	"SendNodeCommand":      "substate=on",
	"RestartNode":          "substate=on",
	"ToggleRadio":          "substate=on",
	"ToggleInternetSwitch": "substate=on",
	"ToggleService":        "substate=on",
	"UpdateSoftware":       "substate=on",
}

var (
	IsDebugMode bool   = false
	InstanceId  string = ServiceName + "-debug"
)
