# Business Service (Analytics)

Read-only gRPC service serving business KPIs for the Ukama analytics system:
home dashboard, sales overview, package performance, billing summary, sites
and inventory readiness.

## Data access

The service reads the shared `analytics` database (rollups, snapshots and
fact tables). The **collector** service is the only writer and the only
service that runs AutoMigrate; this service connects without migrating and
mirrors the model structs it reads (`pkg/db/model.go`, contract in
`../docs/schema.md`).

## RPCs

- `GetHome` — revenue today, active customers, data sold, network uptime, site summaries, top packages, recent activity
- `GetSalesOverview` — revenue / purchases / avg purchase / paid customers with deltas vs the previous window, revenue trend, revenue by site/package
- `GetPackagePerformance` — MRR, ARPU, top plan revenue/share, per-package rows and revenue mix
- `GetBillingSummary` — current balance, last invoice amount, invoice history
- `GetSites` / `GetSite` — per-site revenue, customers, data, uptime
- `GetInventoryReadiness` — available/active SIMs, available/deployed nodes

KPI formulas live in `pkg/server/formula.go`; window resolution
(today/week/month/custom, prev-window for deltas) in `pkg/server/window.go`.

## Development

```sh
make gen    # generate pb/gen from pb/business.proto
make test   # unit tests
make build  # static linux binary in bin/business
```
