/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

// TODO: This service should have be called intent-monitor
// or just intent, not operation-monitor

const (
	SystemName  = "node"
	ServiceName = "operationmonitor"
)

var (
	IsDebugMode bool   = false
	InstanceId  string = ServiceName + "-debug"
)

const (
	RuleKeySubstate = "substate"
	SubstateOn      = "on"
	RuleSubstateOn  = RuleKeySubstate + "=" + SubstateOn
)
