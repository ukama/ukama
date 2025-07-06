/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/gorm"
)

type Sim struct {
	gorm.Model
	Iccid          string `gorm:"index:idx_iccid,unique"`
	Msisdn         string
	IsAllocated    bool
	IsFailed       bool
	SimType        ukama.SimType
	SmDpAddress    string
	ActivationCode string
	QrCode         string
	IsPhysical     bool
}
