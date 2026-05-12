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
	NodeID         string         `gorm:"column:node_id;primaryKey;not null" json:"nodeId"`
	NodeType       ukama.NodeType `gorm:"column:node_type;not null" json:"nodeType"`
	FirstSeenAt    time.Time      `gorm:"column:first_seen_at;not null" json:"firstSeenAt"`
	LastSeenAt     time.Time      `gorm:"column:last_seen_at;not null" json:"lastSeenAt"`
	LastReportedAt *time.Time     `gorm:"column:last_reported_at" json:"lastReportedAt,omitempty"`
}

func (Node) TableName() string {
	return "health_nodes"
}

type HealthReport struct {
	ID            uuid.UUID       `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	NodeID        string          `gorm:"column:node_id;not null;index" json:"nodeId"`
	NodeType      ukama.NodeType  `gorm:"column:node_type;not null;index" json:"nodeType"`
	SchemaVersion string          `gorm:"column:schema_version;not null" json:"schemaVersion"`
	ReportedAt    time.Time       `gorm:"column:reported_at;not null;index" json:"reportedAt"`
	ReceivedAt    time.Time       `gorm:"column:received_at;not null" json:"receivedAt"`
	ParseStatus   ukama.AppStatus `gorm:"column:parse_status;not null" json:"parseStatus"`
	ParseError    string          `gorm:"column:parse_error;not null;default:''" json:"parseError"`
	Payload       json.RawMessage `gorm:"column:payload;type:jsonb;not null" json:"payload"`
}

type NodeLatestHealth struct {
	NodeID        string          `gorm:"column:node_id;primaryKey;not null" json:"nodeId"`
	NodeType      ukama.NodeType  `gorm:"column:node_type;not null" json:"nodeType"`
	ReportID      uuid.UUID       `gorm:"column:report_id;type:uuid;not null" json:"reportId"`
	SchemaVersion string          `gorm:"column:schema_version;not null" json:"schemaVersion"`
	ReportedAt    time.Time       `gorm:"column:reported_at;not null" json:"reportedAt"`
	ReceivedAt    time.Time       `gorm:"column:received_at;not null" json:"receivedAt"`
	ParseStatus   ukama.AppStatus `gorm:"column:parse_status;not null" json:"parseStatus"`
	ParseError    string          `gorm:"column:parse_error;not null;default:''" json:"parseError"`
	Payload       json.RawMessage `gorm:"column:payload;type:jsonb;not null" json:"payload"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;not null" json:"updatedAt"`
}