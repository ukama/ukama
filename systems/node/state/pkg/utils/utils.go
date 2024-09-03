/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"github.com/ukama/ukama/systems/common/ukama"
)

type NotificationType string
type SeverityType string

const (
	Event NotificationType = "event"
	Alert NotificationType = "alert"
)

const (
	Fatal    SeverityType = "fatal"
	Critical SeverityType = "critical"
	High     SeverityType = "high"
	Medium   SeverityType = "medium"
	Low      SeverityType = "low"
	Clean    SeverityType = "clean"
	Log      SeverityType = "log"
	Warning  SeverityType = "warning"
	Debug    SeverityType = "debug"
	Trace    SeverityType = "trace"
)

func GetNodeStateBySeverity(severity SeverityType) ukama.NodeStateEnum {
	switch severity {
	case Fatal, Critical, High:
		return ukama.StateFaulty
	case Medium:
		return ukama.StateOperational
	default:
		return ukama.StateOperational
	}
}

func ToSeverityType(severity string) SeverityType {
	return SeverityType(severity)
}
