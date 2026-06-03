/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

// Package db exposes repository code over the shared analytics schema.
//
// The actual GORM model structs live in systems/analytics/schema so the
// collector and read API use one database contract. Keep repository code
// service-local, but do not duplicate table structs here.
package db

import schema "github.com/ukama/ukama/systems/analytics/schema"

type EventLog = schema.EventLog
type EventError = schema.EventError
type RefreshState = schema.RefreshState
type RollupState = schema.RollupState
type NetworkSnapshot = schema.NetworkSnapshot
type SiteSnapshot = schema.SiteSnapshot
type NodeSnapshot = schema.NodeSnapshot
type CustomerSnapshot = schema.CustomerSnapshot
type SimSnapshot = schema.SimSnapshot
type SimBatchSnapshot = schema.SimBatchSnapshot
type PackageSnapshot = schema.PackageSnapshot
type InventorySnapshot = schema.InventorySnapshot
type BillingSnapshot = schema.BillingSnapshot
type HealthReportSnapshot = schema.HealthReportSnapshot
type PaymentEvent = schema.PaymentEvent
type UsageEvent = schema.UsageEvent
type MetricSample = schema.MetricSample
type AlarmEvent = schema.AlarmEvent
type NodeStateEvent = schema.NodeStateEvent
type SiteStateEvent = schema.SiteStateEvent
type CustomerEvent = schema.CustomerEvent
type SimEvent = schema.SimEvent
type PackageEvent = schema.PackageEvent
type InventoryEvent = schema.InventoryEvent
type NodeStateInterval = schema.NodeStateInterval
type SiteStateInterval = schema.SiteStateInterval
type CustomerPackageInterval = schema.CustomerPackageInterval
type SimStateInterval = schema.SimStateInterval
type MaintenanceWindow = schema.MaintenanceWindow
type BusinessSalesRollupDaily = schema.BusinessSalesRollupDaily
type BusinessPackageRollupDaily = schema.BusinessPackageRollupDaily
type BusinessSiteRollupDaily = schema.BusinessSiteRollupDaily
type BusinessInventoryRollupDaily = schema.BusinessInventoryRollupDaily
type BusinessBillingRollupDaily = schema.BusinessBillingRollupDaily
type CustomerUsageRollupDaily = schema.CustomerUsageRollupDaily
type CustomerStateRollupDaily = schema.CustomerStateRollupDaily
type NetworkHealthRollupHourly = schema.NetworkHealthRollupHourly
type SiteHealthRollupHourly = schema.SiteHealthRollupHourly
type NodeHealthRollupHourly = schema.NodeHealthRollupHourly
type MetricRollupHourly = schema.MetricRollupHourly
type AlarmRollupDaily = schema.AlarmRollupDaily
type RadioRollupHourly = schema.RadioRollupHourly
type BackhaulRollupHourly = schema.BackhaulRollupHourly
type PowerRollupHourly = schema.PowerRollupHourly
