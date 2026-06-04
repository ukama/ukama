# Customer Analytics Service

gRPC service of the **analytics** system serving customer-facing analytics:
customer base overview KPIs, customer list/search/detail, support diagnosis,
sims and sim pool.

## Data access

This service is **read-only**. It connects to the shared `analytics` database
owned by the `collector` service and never migrates or writes to it. The model
structs in `pkg/db/model.go` are mirrors of the collector's models
(source of truth: `collector/pkg/db/model.go`, contract: `../docs/schema.md`).

It registers with the message bus client for lifecycle only (no listener
routes, no event publishing).

## RPCs

- `GetOverview` — total/active/new/expired/failed-activation KPIs, with deltas
  vs the previous window of equal length.
- `List` / `Search` — paginated customer rows (search is ILIKE on
  name/email/sim iccid).
- `Get` — customer detail with usage KPIs and package history.
- `GetSupport` — support diagnosis: derived signals (sim/package/site
  health/usage/last seen), likely issue, recommended action, escalation flag,
  recent activity.
- `GetSims` / `GetSimPool` — sims and sim pool KPIs, including the
  `low_stock` KPI (available < `SimLowStockThreshold`, default 50).

## Development

```sh
make gen    # generate pb/gen from pb/customer.proto
make test
make build
```
