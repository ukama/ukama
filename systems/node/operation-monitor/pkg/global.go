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
	ServiceName = "operation-monitor"
)

const (
	DefaultDeadlineTTL = 10 * time.Minute
	SweeperInterval    = 30 * time.Second
)

// Action → completion rule fallback when caller doesn't supply one.
// TODO: move to config or per-action proto when we add more actions.
var DefaultCompletionRule = map[string]string{
	"RestartNode":          "state=Operational",
	"RestartSite":          "state=Operational",
	"RestartNodes":         "state=Operational",
	"ToggleRfSwitch":       "state=Operational",
	"ToggleInternetSwitch": "state=Operational",
	"ToggleNodeService":    "state=Operational",
	"UpdateSoftware":       "state=Operational",
}

var (
	IsDebugMode bool   = false
	InstanceId  string = ServiceName + "-debug"
)
