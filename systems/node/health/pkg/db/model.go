/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"encoding/json"
	"time"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
)

type Node struct {
	NodeID         string         `gorm:"column:nodeId;primaryKey;not null" json:"nodeId"`
	NodeType       ukama.NodeType `gorm:"column:nodeType;not null" json:"nodeType"`
	FirstSeenAt    time.Time      `gorm:"column:firstSeenAt;not null" json:"firstSeenAt"`
	LastSeenAt     time.Time      `gorm:"column:lastSeenAt;not null" json:"lastSeenAt"`
	LastReportedAt *time.Time     `gorm:"column:lastReportedAt" json:"lastReportedAt,omitempty"`
}

func (Node) TableName() string {
	return "health_nodes"
}

type HealthReport struct {
	ID            uuid.UUID       `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	NodeID        string          `gorm:"column:nodeId;not null;index" json:"nodeId"`
	NodeType      ukama.NodeType  `gorm:"column:nodeType;not null;index" json:"nodeType"`
	SchemaVersion string          `gorm:"column:schemaVersion;not null" json:"schemaVersion"`
	ReportedAt    time.Time       `gorm:"column:reportedAt;not null;index" json:"reportedAt"`
	ReceivedAt    time.Time       `gorm:"column:receivedAt;not null" json:"receivedAt"`
	ParseStatus   ukama.AppStatus `gorm:"column:parseStatus;not null" json:"parseStatus"`
	ParseError    string          `gorm:"column:parseError;not null;default:''" json:"parseError"`
	Payload       json.RawMessage `gorm:"column:payload;type:jsonb;not null" json:"payload"`
}

type NodeLatestHealth struct {
	NodeID        string          `gorm:"column:nodeId;primaryKey;not null" json:"nodeId"`
	NodeType      ukama.NodeType  `gorm:"column:nodeType;not null" json:"nodeType"`
	ReportID      uuid.UUID       `gorm:"column:reportId;type:uuid;not null" json:"reportId"`
	SchemaVersion string          `gorm:"column:schemaVersion;not null" json:"schemaVersion"`
	ReportedAt    time.Time       `gorm:"column:reportedAt;not null" json:"reportedAt"`
	ReceivedAt    time.Time       `gorm:"column:receivedAt;not null" json:"receivedAt"`
	ParseStatus   ukama.AppStatus `gorm:"column:parseStatus;not null" json:"parseStatus"`
	ParseError    string          `gorm:"column:parseError;not null;default:''" json:"parseError"`
	Payload       json.RawMessage `gorm:"column:payload;type:jsonb;not null" json:"payload"`
	UpdatedAt     time.Time       `gorm:"column:updatedAt;not null" json:"updatedAt"`
}
