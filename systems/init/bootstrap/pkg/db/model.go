/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Node struct {
	gorm.Model
	Id uuid.UUID  `gorm:"type:uuid;unique"`
	NodeId string `gorm:"type:string;uniqueIndex:node_id_idx_case_insensetive,expression:lower(node_id);size:23;not null"`
	MeshPodName string `gorm:"type:string;size:255;not null"`
}

