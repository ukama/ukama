/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
)

type State struct {
	Id              uuid.UUID     `gorm:"primaryKey;type:string;uniqueIndex:idx_nodestate_id_case_insensitive,expression:lower(id),where:deleted_at is null;size:23;not null"`
	NodeId          string        `gorm:"type:string;uniqueIndex:idx_nodestate_node_id_case_insensitive,expression:lower(node_id),where:deleted_at is null;size:23;not null"`
	State           ukama.NodeStateEnum `gorm:"type:uint;not null"`
	LastHeartbeat   time.Time
	LastStateChange time.Time
	Connectivity   Connectivity `gorm:"type:uint;not null"`
	Type            string `gorm:"type:string;not null"`
	Version         string `gorm:"type:string"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}


type Connectivity uint8

const (
	Unknown Connectivity = iota
	Offline Connectivity = 1 /* Not connected */
	Online  Connectivity = 2 /* Connected */
)
func (c *Connectivity) Scan(value interface{}) error {
	*c = Connectivity(uint8(value.(int64)))

	return nil
}

func (c Connectivity) Value() (driver.Value, error) {
	return int64(c), nil
}
func (c Connectivity) String() string {
	cs := map[Connectivity]string{
		Unknown: "unkown",
		Offline: "offline",
		Online:  "online",
	}

	return cs[c]
}

func ParseConnectivityState(s string) Connectivity {
	switch strings.ToLower(s) {
	case "offline":
		return Offline
	case "online":
		return Online
	default:
		return Unknown
	}
}
