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

type SimSnapshot struct {
	SimId       string `gorm:"primaryKey"`
	Iccid       string `gorm:"index"`
	Status      string /* available/assigned/active/suspended/faulty */
	CustomerId  uuid.UUID `gorm:"index;type:uuid"`
	BatchId     string    `gorm:"index"`
	AllocatedAt *time.Time
	UpdatedAt   time.Time
}

func (SimSnapshot) TableName() string { return "analytics_sim_snapshots" }

type SimBatchSnapshot struct {
	BatchId    string `gorm:"primaryKey"`
	Quantity   uint32
	Assigned   uint32
	UploadedAt *time.Time
	UpdatedAt  time.Time
}

func (SimBatchSnapshot) TableName() string { return "analytics_sim_batch_snapshots" }

type PackageSnapshot struct {
	PackageId         uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name              string
	Price             float64
	Currency          string
	DurationDays      uint32
	DataQuotaMb       float64
	Status            string
	ActiveSubscribers uint32
	UpdatedAt         time.Time
}

func (PackageSnapshot) TableName() string { return "analytics_package_snapshots" }

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

/* Facts (append-only; written by collector event handlers) */

type UsageEvent struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	CustomerId uuid.UUID `gorm:"index;type:uuid"`
	SimId      string    `gorm:"index"`
	BytesUsed  uint64
	StartAt    time.Time `gorm:"index"`
	EndAt      time.Time
	CreatedAt  time.Time
}

func (UsageEvent) TableName() string { return "analytics_usage_events" }

type SimEvent struct {
	Id         uint64 `gorm:"primaryKey;autoIncrement"`
	SimId      string `gorm:"index"`
	Kind       string /* allocate/activate/add_package/active_package/remove_package/delete/upload */
	OccurredAt time.Time `gorm:"index"`
}

func (SimEvent) TableName() string { return "analytics_sim_events" }

type CustomerEvent struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	CustomerId uuid.UUID `gorm:"index;type:uuid"`
	Kind       string    /* create/update/delete/activation_failed */
	OccurredAt time.Time `gorm:"index"`
}

func (CustomerEvent) TableName() string { return "analytics_customer_events" }

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

type CustomerPackageInterval struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	CustomerId uuid.UUID `gorm:"index;type:uuid"`
	PackageId  uuid.UUID `gorm:"type:uuid"`
	State      string    /* active/expired */
	StartAt    time.Time
	EndAt      *time.Time
}

func (CustomerPackageInterval) TableName() string { return "analytics_customer_package_intervals" }

/* Rollups (rebuilt by collector; read by owning service) */

type CustomerUsageRollupDaily struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	Day        time.Time `gorm:"index;uniqueIndex:idx_customer_usage_day_customer"`
	CustomerId uuid.UUID `gorm:"index;type:uuid;uniqueIndex:idx_customer_usage_day_customer"`
	DataUsedMb float64
}

func (CustomerUsageRollupDaily) TableName() string { return "analytics_customer_usage_rollup_daily" }

type CustomerStateRollupDaily struct {
	Id                uint64    `gorm:"primaryKey;autoIncrement"`
	Day               time.Time `gorm:"index;uniqueIndex:idx_customer_state_day_network"`
	NetworkId         uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_customer_state_day_network"`
	Total             uint32
	Active            uint32
	New               uint32
	Expired           uint32
	FailedActivations uint32
}

func (CustomerStateRollupDaily) TableName() string { return "analytics_customer_state_rollup_daily" }

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
