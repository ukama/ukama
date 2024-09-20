/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import "github.com/ukama/ukama/systems/common/uuid"


type NodeState struct {
	Id            uuid.UUID `gorm:"primaryKey;type:uuid"`
    NodeId          string         `gorm:"not null"`
    PreviousStateId *uuid.UUID          `gorm:"default:null"`
    CurrentState    string         `gorm:"not null"`
    SubState        string         `gorm:"not null"`
    Events          []string       `gorm:"type:json"` 
    Severity        string         `gorm:"not null"`
}
