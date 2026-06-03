# Analytics DB schema (contract)

Single database **`analytics`** on `postgresd-analytics`. **Collector is the only writer and the only service that runs AutoMigrate.** Business/customer/network mirror the model structs they read (read-only repos) — struct definitions must stay byte-identical to collector's.

Conventions: GORM models; every model overrides `TableName()` to the `analytics_*` name below; uuids are `uuid.UUID` (`github.com/ukama/ukama/systems/common/uuid`) stored as `type:uuid`; timestamps `time.Time`; soft delete not used (facts are append-only).

## Foundation (collector-only)

| Table | Columns |
|---|---|
| `analytics_event_logs` | Id uint64 pk auto; RoutingKey string idx; MsgId string uniqueIndex (idempotency key); Payload datatypes.JSON; OccurredAt time idx; CreatedAt |
| `analytics_event_errors` | Id pk auto; RoutingKey string; Reason string; Payload datatypes.JSON; CreatedAt |
| `analytics_refresh_states` | Source string pk (registry/subscriber/dataplan/metrics/node/inventory/billing); Status string (ok/running/failed/stale); Detail string; LastRunAt; LastSuccessAt |
| `analytics_rollup_states` | Rollup string pk (e.g. business_sales_daily); Watermark time; Dirty bool |

## Snapshots (current state; upserted by collector)

| Table | Columns |
|---|---|
| `analytics_network_snapshots` | NetworkId uuid pk; Name; Status; UpdatedAt |
| `analytics_site_snapshots` | SiteId uuid pk; NetworkId uuid idx; Name; Status (online/degraded/offline); Latitude float64; Longitude float64; NodeCount uint32; UpdatedAt |
| `analytics_node_snapshots` | NodeId string pk; SiteId uuid idx; NetworkId uuid idx; Name; Type; Status (online/offline/configuring/needs_attention); Connectivity; LastTelemetryAt *time; UpdatedAt |
| `analytics_customer_snapshots` | CustomerId uuid pk; NetworkId uuid idx; Name; Email idx; Status (active/inactive/expired); PackageId uuid; PackageName; PackageStatus; SimIccid; SimStatus; SiteId uuid idx; LastSeenAt *time; SourceCreatedAt *time; UpdatedAt |
| `analytics_sim_snapshots` | SimId string pk; Iccid string idx; Status (available/assigned/active/suspended/faulty); CustomerId uuid idx; BatchId string idx; AllocatedAt *time; UpdatedAt |
| `analytics_sim_batch_snapshots` | BatchId string pk; Quantity uint32; Assigned uint32; UploadedAt *time; UpdatedAt |
| `analytics_package_snapshots` | PackageId uuid pk; Name; Price float64; Currency; DurationDays uint32; DataQuotaMb float64; Status; ActiveSubscribers uint32; UpdatedAt |
| `analytics_inventory_snapshots` | ComponentId string pk; Type; State (available/deployed/rma); NodeId string idx; UpdatedAt |
| `analytics_billing_snapshots` | Id uint32 pk (always 1, org-level); Balance float64; PaymentMethodStatus; LastInvoiceAt *time; UpdatedAt |
| `analytics_health_report_snapshots` | NodeId string pk; ReportedAt time; Payload datatypes.JSON; UpdatedAt |

## Facts (append-only; written by collector event handlers)

| Table | Columns |
|---|---|
| `analytics_payment_events` | Id uint64 pk auto; ExternalId string uniqueIndex; CustomerId uuid idx; PackageId uuid idx; SiteId uuid idx; NetworkId uuid idx; Amount float64; Currency; Status (success/failed); PaidAt time idx; CreatedAt |
| `analytics_usage_events` | Id pk auto; CustomerId uuid idx; SimId string idx; BytesUsed uint64; StartAt time idx; EndAt time; CreatedAt |
| `analytics_metric_samples` | Id pk auto; Metric string idx; ResourceType string; ResourceId string idx; Value float64; Unit string; SampledAt time idx |
| `analytics_alarm_events` | Id pk auto; AlarmId string uniqueIndex; Severity (critical/warning); State (open/closed); ResourceType (site/node/backhaul/power/radio); ResourceId string idx; Description; CustomersAffected uint32; RevenueAtRisk float64; RecommendedAction; OpenedAt time idx; ClosedAt *time |
| `analytics_node_state_events` | Id pk auto; NodeId string idx; State; OccurredAt time idx |
| `analytics_site_state_events` | Id pk auto; SiteId uuid idx; State; OccurredAt time idx |
| `analytics_customer_events` | Id pk auto; CustomerId uuid idx; Kind (create/update/delete/activation_failed); OccurredAt time idx |
| `analytics_sim_events` | Id pk auto; SimId string idx; Kind (allocate/activate/add_package/active_package/remove_package/delete/upload); OccurredAt time idx |
| `analytics_package_events` | Id pk auto; PackageId uuid idx; Kind (create/update/delete); OccurredAt time idx |
| `analytics_inventory_events` | Id pk auto; ComponentId string idx; Kind (sync); OccurredAt time idx |

## Intervals (derived by collector from state events)

| Table | Columns |
|---|---|
| `analytics_node_state_intervals` | Id pk auto; NodeId string idx; State; StartAt time idx; EndAt *time; DurationSeconds float64 |
| `analytics_site_state_intervals` | Id pk auto; SiteId uuid idx; State; StartAt time idx; EndAt *time; DurationSeconds float64 |
| `analytics_customer_package_intervals` | Id pk auto; CustomerId uuid idx; PackageId uuid; State (active/expired); StartAt; EndAt *time |
| `analytics_sim_state_intervals` | Id pk auto; SimId string idx; State; StartAt; EndAt *time |
| `analytics_maintenance_windows` | Id pk auto; ResourceType; ResourceId string idx; StartAt; EndAt; Reason |

## Rollups (rebuilt by collector; read by owning service)

Business:

| Table | Columns |
|---|---|
| `analytics_business_sales_rollup_daily` | Id pk auto; Day time idx (uniqueIndex day+network_id+site_id); NetworkId uuid; SiteId uuid; Revenue float64; Purchases uint32; PaidCustomers uint32; DataSoldMb float64 |
| `analytics_business_package_rollup_daily` | Id; Day idx (unique day+package_id); PackageId uuid; SoldCount uint32; Revenue float64; DataUsedMb float64 |
| `analytics_business_site_rollup_daily` | Id; Day idx (unique day+site_id); SiteId uuid; Revenue float64; Customers uint32; DataUsedMb float64; UptimePercent float64 |
| `analytics_business_inventory_rollup_daily` | Id; Day uniqueIndex; AvailableSims uint32; ActiveSims uint32; AvailableNodes uint32; DeployedNodes uint32 |
| `analytics_business_billing_rollup_daily` | Id; Day uniqueIndex; InvoicedAmount float64; InvoiceCount uint32 |

Customer:

| Table | Columns |
|---|---|
| `analytics_customer_usage_rollup_daily` | Id; Day idx (unique day+customer_id); CustomerId uuid idx; DataUsedMb float64 |
| `analytics_customer_state_rollup_daily` | Id; Day idx (unique day+network_id); NetworkId uuid; Total uint32; Active uint32; New uint32; Expired uint32; FailedActivations uint32 |

Network:

| Table | Columns |
|---|---|
| `analytics_network_health_rollup_hourly` | Id; Hour idx (unique hour+network_id); NetworkId uuid; SitesOnline uint32; SitesTotal uint32; NodesOnline uint32; NodesTotal uint32; UptimePercent float64; OpenAlarms uint32; CriticalAlarms uint32 |
| `analytics_site_health_rollup_hourly` | Id; Hour idx (unique hour+site_id); SiteId uuid; UptimePercent float64; BackhaulLatencyMs float64; BatteryPercent float64; Status string |
| `analytics_node_health_rollup_hourly` | Id; Hour idx (unique hour+node_id); NodeId string; UptimePercent float64; Status string |
| `analytics_metric_rollup_hourly` | Id; Hour idx (unique hour+metric+resource_id); Metric string; ResourceType string; ResourceId string; Avg float64; Min float64; Max float64; Count uint32 |
| `analytics_alarm_rollup_daily` | Id; Day uniqueIndex; Opened uint32; Closed uint32; Critical uint32; Warning uint32 |
| `analytics_radio_rollup_hourly` | Id; Hour idx (unique hour+node_id); NodeId string; ActiveUes uint32; AttachFailures uint32; DlThroughputMbps float64; UlThroughputMbps float64; SignalDbm float64 |
| `analytics_backhaul_rollup_hourly` | Id; Hour idx (unique hour+site_id); SiteId uuid; LatencyMs float64; DlMbps float64; UlMbps float64; PacketLossPercent float64 |
| `analytics_power_rollup_hourly` | Id; Hour idx (unique hour+site_id); SiteId uuid; BatteryPercent float64; BatteryVoltage float64; LoadWatts float64; SolarWatts float64; TemperatureC float64 |

## Ownership matrix

| Service | Reads | Writes |
|---|---|---|
| collector | all | all |
| business | sales/package/site/inventory/billing rollups; payment_events; package/site/sim/inventory/billing snapshots; event_logs (recent activity); site health rollups (uptime) | — |
| customer | customer/sim/sim_batch/package snapshots; customer_state/usage rollups; customer_package_intervals; usage_events; sim_events; customer_events; event_logs | — |
| network | site/node/network/inventory snapshots; all network rollups; alarm_events; node/site state intervals; metric_samples; event_logs | — |
