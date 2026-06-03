/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

// Read-only mirror; source of truth: collector/pkg/db/model.go
// (contract: systems/analytics/docs/schema.md). The collector service is the
// only writer and the only service running AutoMigrate against the shared
// "analytics" database. These struct definitions must stay byte-identical
// to collector's.

import (
	"time"

	"github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/datatypes"
)

/* Foundation */

type EventLog struct {
	Id         uint64         `gorm:"primaryKey;autoIncrement"`
	RoutingKey string         `gorm:"index"`
	MsgId      string         `gorm:"uniqueIndex"`
	Payload    datatypes.JSON `gorm:"type:jsonb"`
	OccurredAt time.Time      `gorm:"index"`
	CreatedAt  time.Time
}

func (EventLog) TableName() string {
	return "analytics_event_logs"
}

/* Snapshots */

type SiteSnapshot struct {
	SiteId    uuid.UUID `gorm:"primaryKey;type:uuid"`
	NetworkId uuid.UUID `gorm:"index;type:uuid"`
	Name      string
	Status    string
	Latitude  float64
	Longitude float64
	NodeCount uint32
	UpdatedAt time.Time
}

func (SiteSnapshot) TableName() string {
	return "analytics_site_snapshots"
}

type SimSnapshot struct {
	SimId       string `gorm:"primaryKey"`
	Iccid       string `gorm:"index"`
	Status      string
	CustomerId  uuid.UUID `gorm:"index;type:uuid"`
	BatchId     string    `gorm:"index"`
	AllocatedAt *time.Time
	UpdatedAt   time.Time
}

func (SimSnapshot) TableName() string {
	return "analytics_sim_snapshots"
}

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

func (PackageSnapshot) TableName() string {
	return "analytics_package_snapshots"
}

type InventorySnapshot struct {
	ComponentId string `gorm:"primaryKey"`
	Type        string
	State       string
	NodeId      string `gorm:"index"`
	UpdatedAt   time.Time
}

func (InventorySnapshot) TableName() string {
	return "analytics_inventory_snapshots"
}

type BillingSnapshot struct {
	Id                  uint32 `gorm:"primaryKey"`
	Balance             float64
	PaymentMethodStatus string
	LastInvoiceAt       *time.Time
	UpdatedAt           time.Time
}

func (BillingSnapshot) TableName() string {
	return "analytics_billing_snapshots"
}

/* Facts */

type PaymentEvent struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	ExternalId string    `gorm:"uniqueIndex"`
	CustomerId uuid.UUID `gorm:"index;type:uuid"`
	PackageId  uuid.UUID `gorm:"index;type:uuid"`
	SiteId     uuid.UUID `gorm:"index;type:uuid"`
	NetworkId  uuid.UUID `gorm:"index;type:uuid"`
	Amount     float64
	Currency   string
	Status     string
	PaidAt     time.Time `gorm:"index"`
	CreatedAt  time.Time
}

func (PaymentEvent) TableName() string {
	return "analytics_payment_events"
}

/* Rollups */

type BusinessSalesRollupDaily struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement"`
	Day           time.Time `gorm:"index;uniqueIndex:idx_sales_day_network_site"`
	NetworkId     uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_sales_day_network_site"`
	SiteId        uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_sales_day_network_site"`
	Revenue       float64
	Purchases     uint32
	PaidCustomers uint32
	DataSoldMb    float64
}

func (BusinessSalesRollupDaily) TableName() string {
	return "analytics_business_sales_rollup_daily"
}

type BusinessPackageRollupDaily struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	Day        time.Time `gorm:"index;uniqueIndex:idx_package_day_package"`
	PackageId  uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_package_day_package"`
	SoldCount  uint32
	Revenue    float64
	DataUsedMb float64
}

func (BusinessPackageRollupDaily) TableName() string {
	return "analytics_business_package_rollup_daily"
}

type BusinessSiteRollupDaily struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement"`
	Day           time.Time `gorm:"index;uniqueIndex:idx_site_day_site"`
	SiteId        uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_site_day_site"`
	Revenue       float64
	Customers     uint32
	DataUsedMb    float64
	UptimePercent float64
}

func (BusinessSiteRollupDaily) TableName() string {
	return "analytics_business_site_rollup_daily"
}

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

type SiteHealthRollupHourly struct {
	Id                uint64    `gorm:"primaryKey;autoIncrement"`
	Hour              time.Time `gorm:"index;uniqueIndex:idx_site_health_hour_site"`
	SiteId            uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_site_health_hour_site"`
	UptimePercent     float64
	BackhaulLatencyMs float64
	BatteryPercent    float64
	Status            string
}

func (SiteHealthRollupHourly) TableName() string {
	return "analytics_site_health_rollup_hourly"
}
