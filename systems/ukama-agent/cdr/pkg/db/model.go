/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
 
package db

import (
	"gorm.io/gorm"
)

type CDR struct {
	gorm.Model
	Session       uint64
	NodeId        string `gorm:"Index:cdr_node_idx;not null"`
	Imsi          string `gorm:"Index:cdr_imsi_idx;not null"`
	Policy        string `gorm:"Index:cdr_policy_idx;not null"`
	ApnName       string
	Ip            string
	StartTime     uint64
	EndTime       uint64
	LastUpdatedAt uint64
	TxBytes       uint64
	RxBytes       uint64
	TotalBytes    uint64
}

type Usage struct {
	gorm.Model
	Imsi             string `gorm:"Index:usage_imsi_idx,unique,not null;"`
	Historical       uint64 /* all data used till last session */
	Usage            uint64 /* usage till now (last session + current session */
	LastSessionUsage uint64 /* usage till last session for current package*/
	LastSessionId    uint64 /* usage till last session for current package*/
	LastNodeId       string
	LastCDRUpdatedAt uint64 /* timestamp for last CDR LasteUpdatedAt */
	Policy           string
}
