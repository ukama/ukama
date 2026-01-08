/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"github.com/jackc/pgtype"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Node struct {
	gorm.Model
	NodeID string `gorm:"type:string;uniqueIndex:node_id_idx_case_insensetive,expression:lower(node_id);size:23;not null"`
	OrgID  uint
	Org    Org
}

type Org struct {
	gorm.Model
	Name        string    `gorm:"uniqueIndex"`
	OrgId       uuid.UUID `gorm:"type:uuid;uniqueIndex:org_id_unique_index,where:deleted_at is null;not null;column_name:org_org_id;"`
	Certificate string
	Ip          pgtype.Inet `gorm:"type:inet"`
	Nodes       []Node
	Systems     []System
}

type System struct {
	gorm.Model
	Name        string `gorm:"type:string;index:sys_idx,unique,composite:sys_idx,expression:lower(name);not null"`
	Uuid        string `gorm:"type:uuid;unique"`
	Certificate string
	ApiGwIp     pgtype.Inet `gorm:"type:inet"`
	ApiGwUrl    string
	ApiGwPort   int32
	NodeGwIp    pgtype.Inet `gorm:"type:inet;default:'0.0.0.0'"`
	NodeGwPort  int32 `gorm:"default:8080"`
	NodeGwHealth uint32 `gorm:"default:100"`
	OrgID       uint `gorm:"type:string;index:sys_idx,unique,composite:sys_idx;not null"`
	Org         Org
	ApiGwHealth uint32 `gorm:"default:100"`
}
