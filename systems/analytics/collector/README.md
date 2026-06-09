# Collector

The collector is the write side of the Ukama analytics system. It is the
**only writer** (and the only service running AutoMigrate) on the shared
`analytics` PostgreSQL database. The business, customer and network analytics
services read the same database through read-only model mirrors.

## What it does

- **Event ingestion** (`pkg/server/event.go`): consumes platform events from
  the message bus and writes them as snapshots, append-only facts and derived
  state intervals. Every event is logged into `analytics_event_logs` with a
  deterministic message id for idempotency; duplicates are skipped, malformed
  payloads are recorded in `analytics_event_errors` and acknowledged.
- **Snapshot refresh** (`pkg/refresh`): on demand (`Refresh` RPC) pulls
  current state from source api-gateways (registry, subscriber, dataplan,
  metrics, node, inventory, billing) and upserts snapshots, tracking per-source
  status in `analytics_refresh_states`.
- **Rollups** (`RebuildRollups` RPC): rebuilds daily/hourly aggregate tables
  from facts via SQL aggregates and maintains rollup watermarks/dirty flags in
  `analytics_rollup_states`.
- **SeedDemo** RPC: debug-mode-only helper for demo data.

## Events consumed

Payments: `payment.success`, `payment.failed`.
Subscribers: `subscriber.create/update/delete`.
Sims: `sim.allocate/activate/addpackage/activepackage/removepackage/delete`,
`sims.upload`.
Packages: `package.create/update/delete`.
Registry: `network.add`, `site.create/update`,
`node.create/update/assign/release`.
Node state: `node.online`, `node.offline`, `node.transition`,
`health.report.store`.
Inventory: `components.sync`. Billing: `invoice.generate`.

Routing keys come from `systems/common/events`; the deployment must configure
`MSGCLIENT_LISTENERROUTES` accordingly.

## How to run

```sh
make gen     # generate pb/gen and mocks (protoc + mockery)
go mod tidy
make build   # builds bin/collector
make test
```

Configuration follows the standard service pattern (env vars / `collector.yaml`
via num30/config). The database name is fixed to `analytics`.
