/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

// Source of truth for the shared "analytics" database schema (see
// systems/analytics/docs/schema.md). The collector is the only writer and the
// only service that runs AutoMigrate. Business/customer/network keep
// read-only mirrors of the structs they consume.
package schema

import (
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/datatypes"
)

/* Foundation (collector-only) */

type EventLog struct {
	Id         uint64 `gorm:"primaryKey;autoIncrement"`
	RoutingKey string `gorm:"index"`
	MsgId      string `gorm:"uniqueIndex"` /* idempotency key */
	Payload    datatypes.JSON
	OccurredAt time.Time `gorm:"index"`
	CreatedAt  time.Time
}

func (EventLog) TableName() string { return "analytics_event_logs" }

type EventError struct {
	Id         uint64 `gorm:"primaryKey;autoIncrement"`
	RoutingKey string
	Reason     string
	Payload    datatypes.JSON
	CreatedAt  time.Time
}

func (EventError) TableName() string { return "analytics_event_errors" }

type RefreshState struct {
	Source        string `gorm:"primaryKey"` /* registry/subscriber/dataplan/metrics/node/inventory/billing */
	Status        string /* ok/running/failed/stale */
	Detail        string
	LastRunAt     time.Time
	LastSuccessAt time.Time
}

func (RefreshState) TableName() string { return "analytics_refresh_states" }

type RollupState struct {
	Rollup    string `gorm:"primaryKey"` /* e.g. business_sales_daily */
	Watermark time.Time
	Dirty     bool
}

func (RollupState) TableName() string { return "analytics_rollup_states" }

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

type SimSnapshot struct {
	SimId       string `gorm:"primaryKey"`
	Iccid       string `gorm:"index"`
	Status      string    /* available/assigned/active/suspended/faulty */
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

type InventorySnapshot struct {
	ComponentId string `gorm:"primaryKey"`
	Type        string
	State       string /* available/deployed/rma */
	NodeId      string `gorm:"index"`
	UpdatedAt   time.Time
}

func (InventorySnapshot) TableName() string { return "analytics_inventory_snapshots" }

type BillingSnapshot struct {
	Id                  uint32 `gorm:"primaryKey"` /* always 1, org-level */
	Balance             float64
	PaymentMethodStatus string
	LastInvoiceAt       *time.Time
	UpdatedAt           time.Time
}

func (BillingSnapshot) TableName() string { return "analytics_billing_snapshots" }

type HealthReportSnapshot struct {
	NodeId     string `gorm:"primaryKey"`
	ReportedAt time.Time
	Payload    datatypes.JSON
	UpdatedAt  time.Time
}

func (HealthReportSnapshot) TableName() string { return "analytics_health_report_snapshots" }

/* Facts (append-only; written by collector event handlers) */

type PaymentEvent struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	ExternalId string    `gorm:"uniqueIndex"`
	CustomerId uuid.UUID `gorm:"index;type:uuid"`
	PackageId  uuid.UUID `gorm:"index;type:uuid"`
	SiteId     uuid.UUID `gorm:"index;type:uuid"`
	NetworkId  uuid.UUID `gorm:"index;type:uuid"`
	Amount     float64
	Currency   string
	Status     string    /* success/failed */
	PaidAt     time.Time `gorm:"index"`
	CreatedAt  time.Time
}

func (PaymentEvent) TableName() string { return "analytics_payment_events" }

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

type NodeStateEvent struct {
	Id         uint64 `gorm:"primaryKey;autoIncrement"`
	NodeId     string `gorm:"index"`
	State      string
	OccurredAt time.Time `gorm:"index"`
}

func (NodeStateEvent) TableName() string { return "analytics_node_state_events" }

type SiteStateEvent struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	SiteId     uuid.UUID `gorm:"index;type:uuid"`
	State      string
	OccurredAt time.Time `gorm:"index"`
}

func (SiteStateEvent) TableName() string { return "analytics_site_state_events" }

type CustomerEvent struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	CustomerId uuid.UUID `gorm:"index;type:uuid"`
	Kind       string    /* create/update/delete/activation_failed */
	OccurredAt time.Time `gorm:"index"`
}

func (CustomerEvent) TableName() string { return "analytics_customer_events" }

type SimEvent struct {
	Id         uint64 `gorm:"primaryKey;autoIncrement"`
	SimId      string `gorm:"index"`
	Kind       string /* allocate/activate/add_package/active_package/remove_package/delete/upload */
	OccurredAt time.Time `gorm:"index"`
}

func (SimEvent) TableName() string { return "analytics_sim_events" }

type PackageEvent struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	PackageId  uuid.UUID `gorm:"index;type:uuid"`
	Kind       string    /* create/update/delete */
	OccurredAt time.Time `gorm:"index"`
}

func (PackageEvent) TableName() string { return "analytics_package_events" }

type InventoryEvent struct {
	Id          uint64 `gorm:"primaryKey;autoIncrement"`
	ComponentId string `gorm:"index"`
	Kind        string /* sync */
	OccurredAt  time.Time `gorm:"index"`
}

func (InventoryEvent) TableName() string { return "analytics_inventory_events" }

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

type CustomerPackageInterval struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	CustomerId uuid.UUID `gorm:"index;type:uuid"`
	PackageId  uuid.UUID `gorm:"type:uuid"`
	State      string    /* active/expired */
	StartAt    time.Time
	EndAt      *time.Time
}

func (CustomerPackageInterval) TableName() string { return "analytics_customer_package_intervals" }

type SimStateInterval struct {
	Id      uint64 `gorm:"primaryKey;autoIncrement"`
	SimId   string `gorm:"index"`
	State   string
	StartAt time.Time
	EndAt   *time.Time
}

func (SimStateInterval) TableName() string { return "analytics_sim_state_intervals" }

type MaintenanceWindow struct {
	Id           uint64 `gorm:"primaryKey;autoIncrement"`
	ResourceType string
	ResourceId   string `gorm:"index"`
	StartAt      time.Time
	EndAt        time.Time
	Reason       string
}

func (MaintenanceWindow) TableName() string { return "analytics_maintenance_windows" }

/* Rollups (rebuilt by collector; read by owning service) */

/* Business */

type BusinessSalesRollupDaily struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement"`
	Day           time.Time `gorm:"index;uniqueIndex:uniq_sales_day_net_site"`
	NetworkId     uuid.UUID `gorm:"type:uuid;uniqueIndex:uniq_sales_day_net_site"`
	SiteId        uuid.UUID `gorm:"type:uuid;uniqueIndex:uniq_sales_day_net_site"`
	Revenue       float64
	Purchases     uint32
	PaidCustomers uint32
	DataSoldMb    float64
}

func (BusinessSalesRollupDaily) TableName() string { return "analytics_business_sales_rollup_daily" }

type BusinessPackageRollupDaily struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	Day        time.Time `gorm:"index;uniqueIndex:uniq_pkg_day_pkg"`
	PackageId  uuid.UUID `gorm:"type:uuid;uniqueIndex:uniq_pkg_day_pkg"`
	SoldCount  uint32
	Revenue    float64
	DataUsedMb float64
}

func (BusinessPackageRollupDaily) TableName() string {
	return "analytics_business_package_rollup_daily"
}

type BusinessSiteRollupDaily struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement"`
	Day           time.Time `gorm:"index;uniqueIndex:uniq_bsite_day_site"`
	SiteId        uuid.UUID `gorm:"type:uuid;uniqueIndex:uniq_bsite_day_site"`
	Revenue       float64
	Customers     uint32
	DataUsedMb    float64
	UptimePercent float64
}

func (BusinessSiteRollupDaily) TableName() string { return "analytics_business_site_rollup_daily" }

type BusinessInventoryRollupDaily struct {
	Id             uint64    `gorm:"primaryKey;autoIncrement"`
	Day            time.Time `gorm:"uniqueIndex"`
	AvailableSims  uint32
	ActiveSims     uint32
	AvailableNodes uint32
	DeployedNodes  uint32
}

func (BusinessInventoryRollupDaily) TableName() string {
	return "analytics_business_inventory_rollup_daily"
}

type BusinessBillingRollupDaily struct {
	Id             uint64    `gorm:"primaryKey;autoIncrement"`
	Day            time.Time `gorm:"uniqueIndex"`
	InvoicedAmount float64
	InvoiceCount   uint32
}

func (BusinessBillingRollupDaily) TableName() string {
	return "analytics_business_billing_rollup_daily"
}

/* Customer */

type CustomerUsageRollupDaily struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	Day        time.Time `gorm:"index;uniqueIndex:uniq_cusage_day_cust"`
	CustomerId uuid.UUID `gorm:"index;type:uuid;uniqueIndex:uniq_cusage_day_cust"`
	DataUsedMb float64
}

func (CustomerUsageRollupDaily) TableName() string { return "analytics_customer_usage_rollup_daily" }

type CustomerStateRollupDaily struct {
	Id                uint64    `gorm:"primaryKey;autoIncrement"`
	Day               time.Time `gorm:"index;uniqueIndex:uniq_cstate_day_net"`
	NetworkId         uuid.UUID `gorm:"type:uuid;uniqueIndex:uniq_cstate_day_net"`
	Total             uint32
	Active            uint32
	New               uint32
	Expired           uint32
	FailedActivations uint32
}

func (CustomerStateRollupDaily) TableName() string { return "analytics_customer_state_rollup_daily" }

/* Network */

type NetworkHealthRollupHourly struct {
	Id             uint64    `gorm:"primaryKey;autoIncrement"`
	Hour           time.Time `gorm:"index;uniqueIndex:uniq_nethealth_hour_net"`
	NetworkId      uuid.UUID `gorm:"type:uuid;uniqueIndex:uniq_nethealth_hour_net"`
	SitesOnline    uint32
	SitesTotal     uint32
	NodesOnline    uint32
	NodesTotal     uint32
	UptimePercent  float64
	OpenAlarms     uint32
	CriticalAlarms uint32
}

func (NetworkHealthRollupHourly) TableName() string {
	return "analytics_network_health_rollup_hourly"
}

type SiteHealthRollupHourly struct {
	Id                uint64    `gorm:"primaryKey;autoIncrement"`
	Hour              time.Time `gorm:"index;uniqueIndex:uniq_sitehealth_hour_site"`
	SiteId            uuid.UUID `gorm:"type:uuid;uniqueIndex:uniq_sitehealth_hour_site"`
	UptimePercent     float64
	BackhaulLatencyMs float64
	BatteryPercent    float64
	Status            string
}

func (SiteHealthRollupHourly) TableName() string { return "analytics_site_health_rollup_hourly" }

type NodeHealthRollupHourly struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement"`
	Hour          time.Time `gorm:"index;uniqueIndex:uniq_nodehealth_hour_node"`
	NodeId        string    `gorm:"uniqueIndex:uniq_nodehealth_hour_node"`
	UptimePercent float64
	Status        string
}

func (NodeHealthRollupHourly) TableName() string { return "analytics_node_health_rollup_hourly" }

type MetricRollupHourly struct {
	Id           uint64    `gorm:"primaryKey;autoIncrement"`
	Hour         time.Time `gorm:"index;uniqueIndex:uniq_metric_hour_metric_res"`
	Metric       string    `gorm:"uniqueIndex:uniq_metric_hour_metric_res"`
	ResourceType string
	ResourceId   string `gorm:"uniqueIndex:uniq_metric_hour_metric_res"`
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
	Hour             time.Time `gorm:"index;uniqueIndex:uniq_radio_hour_node"`
	NodeId           string    `gorm:"uniqueIndex:uniq_radio_hour_node"`
	ActiveUes        uint32
	AttachFailures   uint32
	DlThroughputMbps float64
	UlThroughputMbps float64
	SignalDbm        float64
}

func (RadioRollupHourly) TableName() string { return "analytics_radio_rollup_hourly" }

type BackhaulRollupHourly struct {
	Id                uint64    `gorm:"primaryKey;autoIncrement"`
	Hour              time.Time `gorm:"index;uniqueIndex:uniq_backhaul_hour_site"`
	SiteId            uuid.UUID `gorm:"type:uuid;uniqueIndex:uniq_backhaul_hour_site"`
	LatencyMs         float64
	DlMbps            float64
	UlMbps            float64
	PacketLossPercent float64
}

func (BackhaulRollupHourly) TableName() string { return "analytics_backhaul_rollup_hourly" }

type PowerRollupHourly struct {
	Id             uint64    `gorm:"primaryKey;autoIncrement"`
	Hour           time.Time `gorm:"index;uniqueIndex:uniq_power_hour_site"`
	SiteId         uuid.UUID `gorm:"type:uuid;uniqueIndex:uniq_power_hour_site"`
	BatteryPercent float64
	BatteryVoltage float64
	LoadWatts      float64
	SolarWatts     float64
	TemperatureC   float64
}

func (PowerRollupHourly) TableName() string { return "analytics_power_rollup_hourly" }
