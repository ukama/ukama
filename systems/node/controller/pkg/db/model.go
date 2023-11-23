/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import "gorm.io/gorm"

type NodeLog struct {
	gorm.Model
	NodeId          string      `gorm:"type:string;uniqueIndex:idx_node_id_case_insensitive,where:deleted_at is null;size:23;not null"`
}