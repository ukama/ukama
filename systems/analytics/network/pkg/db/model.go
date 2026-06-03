/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

// read-only mirror; source of truth: collector/pkg/db/model.go
package db

import (
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/datatypes"
)

/* Snapshots (current state; upserted by collector) */

type NetworkSnapshot struct {
	NetworkId uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name      string
	Status    string
	UpdatedAt time.Time
}

func (NetworkSnapshot) TableName() string { return "analytics_network_snapshots" }

type SiteSnapshot struct {
	SiteId    uuid.UUID `gorm:"primaryKey;type:uuid"`
	NetworkId uuid.UUID `gorm:"index;type:uuid"`
	Name      string
	Status    string /* online/degraded/offline */
	Latitude  float64
	Longitude float64
	NodeCount uint32
	UpdatedAt time.Time
}

func (SiteSnapshot) TableName() string { return "analytics_site_snapshots" }

type NodeSnapshot struct {
	NodeId          string    `gorm:"primaryKey"`
	SiteId          uuid.UUID `gorm:"index;type:uuid"`
	NetworkId       uuid.UUID `gorm:"index;type:uuid"`
	Name            string
	Type            string
	Status          string /* online/offline/configuring/needs_attention */
	Connectivity    string
	LastTelemetryAt *time.Time
	UpdatedAt       time.Time
}

func (NodeSnapshot) TableName() string { return "analytics_node_snapshots" }

type CustomerSnapshot struct {
	CustomerId      uuid.UUID `gorm:"primaryKey;type:uuid"`
	NetworkId       uuid.UUID `gorm:"index;type:uuid"`
	Name            string
	Email           string `gorm:"index"`
	Status          string /* active/inactive/expired */
	PackageId       uuid.UUID `gorm:"type:uuid"`
	PackageName     string
	PackageStatus   string
	SimIccid        string
	SimStatus       string
	SiteId          uuid.UUID `gorm:"index;type:uuid"`
	LastSeenAt      *time.Time
	SourceCreatedAt *time.Time
	UpdatedAt       time.Time
}

func (CustomerSnapshot) TableName() string { return "analytics_customer_snapshots" }

type InventorySnapshot struct {
	ComponentId string `gorm:"primaryKey"`
	Type        string
	State       string /* available/deployed/rma */
	NodeId      string `gorm:"index"`
	UpdatedAt   time.Time
}

func (InventorySnapshot) TableName() string { return "analytics_inventory_snapshots" }

/* Facts (append-only; written by collector event handlers) */

type MetricSample struct {
	Id           uint64 `gorm:"primaryKey;autoIncrement"`
	Metric       string `gorm:"index"`
	ResourceType string
	ResourceId   string `gorm:"index"`
	Value        float64
	Unit         string
	SampledAt    time.Time `gorm:"index"`
}

func (MetricSample) TableName() string { return "analytics_metric_samples" }

type AlarmEvent struct {
	Id                uint64 `gorm:"primaryKey;autoIncrement"`
	AlarmId           string `gorm:"uniqueIndex"`
	Severity          string /* critical/warning */
	State             string /* open/closed */
	ResourceType      string /* site/node/backhaul/power/radio */
	ResourceId        string `gorm:"index"`
	Description       string
	CustomersAffected uint32
	RevenueAtRisk     float64
	RecommendedAction string
	OpenedAt          time.Time `gorm:"index"`
	ClosedAt          *time.Time
}

func (AlarmEvent) TableName() string { return "analytics_alarm_events" }

type EventLog struct {
	Id         uint64 `gorm:"primaryKey;autoIncrement"`
	RoutingKey string `gorm:"index"`
	MsgId      string `gorm:"uniqueIndex"`
	Payload    datatypes.JSON
	OccurredAt time.Time `gorm:"index"`
	CreatedAt  time.Time
}

func (EventLog) TableName() string { return "analytics_event_logs" }

/* Intervals (derived by collector from state events) */

type NodeStateInterval struct {
	Id              uint64 `gorm:"primaryKey;autoIncrement"`
	NodeId          string `gorm:"index"`
	State           string
	StartAt         time.Time `gorm:"index"`
	EndAt           *time.Time
	DurationSeconds float64
}

func (NodeStateInterval) TableName() string { return "analytics_node_state_intervals" }

type SiteStateInterval struct {
	Id              uint64    `gorm:"primaryKey;autoIncrement"`
	SiteId          uuid.UUID `gorm:"index;type:uuid"`
	State           string
	StartAt         time.Time `gorm:"index"`
	EndAt           *time.Time
	DurationSeconds float64
}

func (SiteStateInterval) TableName() string { return "analytics_site_state_intervals" }

/* Rollups (rebuilt by collector; read by owning service) */

type NetworkHealthRollupHourly struct {
	Id             uint64    `gorm:"primaryKey;autoIncrement"`
	Hour           time.Time `gorm:"index;uniqueIndex:idx_network_health_hour_network"`
	NetworkId      uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_network_health_hour_network"`
	SitesOnline    uint32
	SitesTotal     uint32
	NodesOnline    uint32
	NodesTotal     uint32
	UptimePercent  float64
	OpenAlarms     uint32
	CriticalAlarms uint32
}

func (NetworkHealthRollupHourly) TableName() string { return "analytics_network_health_rollup_hourly" }

type SiteHealthRollupHourly struct {
	Id                uint64    `gorm:"primaryKey;autoIncrement"`
	Hour              time.Time `gorm:"index;uniqueIndex:idx_site_health_hour_site"`
	SiteId            uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_site_health_hour_site"`
	UptimePercent     float64
	BackhaulLatencyMs float64
	BatteryPercent    float64
	Status            string
}

func (SiteHealthRollupHourly) TableName() string { return "analytics_site_health_rollup_hourly" }

type NodeHealthRollupHourly struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement"`
	Hour          time.Time `gorm:"index;uniqueIndex:idx_node_health_hour_node"`
	NodeId        string    `gorm:"uniqueIndex:idx_node_health_hour_node"`
	UptimePercent float64
	Status        string
}

func (NodeHealthRollupHourly) TableName() string { return "analytics_node_health_rollup_hourly" }

type MetricRollupHourly struct {
	Id           uint64    `gorm:"primaryKey;autoIncrement"`
	Hour         time.Time `gorm:"index;uniqueIndex:idx_metric_rollup_hour_metric_resource"`
	Metric       string    `gorm:"uniqueIndex:idx_metric_rollup_hour_metric_resource"`
	ResourceType string
	ResourceId   string `gorm:"uniqueIndex:idx_metric_rollup_hour_metric_resource"`
	Avg          float64
	Min          float64
	Max          float64
	Count        uint32
}

func (MetricRollupHourly) TableName() string { return "analytics_metric_rollup_hourly" }

type AlarmRollupDaily struct {
	Id       uint64    `gorm:"primaryKey;autoIncrement"`
	Day      time.Time `gorm:"uniqueIndex"`
	Opened   uint32
	Closed   uint32
	Critical uint32
	Warning  uint32
}

func (AlarmRollupDaily) TableName() string { return "analytics_alarm_rollup_daily" }

type RadioRollupHourly struct {
	Id               uint64    `gorm:"primaryKey;autoIncrement"`
	Hour             time.Time `gorm:"index;uniqueIndex:idx_radio_rollup_hour_node"`
	NodeId           string    `gorm:"uniqueIndex:idx_radio_rollup_hour_node"`
	ActiveUes        uint32
	AttachFailures   uint32
	DlThroughputMbps float64
	UlThroughputMbps float64
	SignalDbm        float64
}

func (RadioRollupHourly) TableName() string { return "analytics_radio_rollup_hourly" }

type BackhaulRollupHourly struct {
	Id                uint64    `gorm:"primaryKey;autoIncrement"`
	Hour              time.Time `gorm:"index;uniqueIndex:idx_backhaul_rollup_hour_site"`
	SiteId            uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_backhaul_rollup_hour_site"`
	LatencyMs         float64
	DlMbps            float64
	UlMbps            float64
	PacketLossPercent float64
}

func (BackhaulRollupHourly) TableName() string { return "analytics_backhaul_rollup_hourly" }

type PowerRollupHourly struct {
	Id             uint64    `gorm:"primaryKey;autoIncrement"`
	Hour           time.Time `gorm:"index;uniqueIndex:idx_power_rollup_hour_site"`
	SiteId         uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_power_rollup_hour_site"`
	BatteryPercent float64
	BatteryVoltage float64
	LoadWatts      float64
	SolarWatts     float64
	TemperatureC   float64
}

func (PowerRollupHourly) TableName() string { return "analytics_power_rollup_hourly" }
