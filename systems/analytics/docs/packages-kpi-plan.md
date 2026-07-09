# Packages KPI Plan

Target: the Packages dashboard (MRR, ARPU, top plan, customers on a plan,
package performance table, top packages, revenue mix).

## New datasets (ingest specs)

| Dataset key | Endpoint | Strategy | Notes |
|---|---|---|---|
| `dataplan.package.getAll` | data-plan `GET /v1/packages` | full_snapshot, entity `package_id` | fields: package_id, name, amount, currency, data_volume, duration, sim_type, active, **network_id** (empty = org-level → valid for every network) |
| `subscriber.sim.list` | subscriber `GET /v1/sim?network_id=` (no status filter) | full_snapshot, for_each networks | all sims, not just active — inactive sims still carry purchase history |
| `subscriber.simpackage.listBySim` | subscriber `GET /v1/sim/{sim_id}/package` | full_snapshot, entity sim_package id, for_each `subscriber.sim.list` (bind sim_id, network_id) | one row per assignment incl. the queue: package_id, start_date, end_date, is_active, as_expired |

Existing datasets reused: `registry.network.getAll`, `subscriber.usage.getBySim`.

## Package↔network resolution rule

A package applies to a network if `package.network_id == network.id` OR
`package.network_id == ""` (org-level). Implemented once as a shared helper
in the algos package; every package KPI uses it.

## KPIs (all scope `network_id`; per-package ones scope `network_id + package_id`)

| KPI | Formula | Mode | Rollup ops |
|---|---|---|---|
| `PACKAGE_SALES` | count of sim-package assignments with `start_date` in window | state→windowed derivation | SUM (daily sold), AVG |
| `PACKAGE_REVENUE` | PACKAGE_SALES × package.amount (org-currency cents) | same | SUM — feeds top-packages + revenue-mix via `breakdown?by=package_id` |
| `MRR` | Σ over currently-active assignments: amount × (30d / package duration) | state gauge | LAST (headline), AVG |
| `ARPU` | MRR ÷ distinct subscribers with ≥1 active assignment | state gauge | LAST |
| `CUSTOMERS_ON_PLAN` | distinct sims (→subscribers) with an active assignment | state gauge | LAST; also breakdown by package_id for "across N plans" |
| `ACTIVE_PLANS` | distinct packages with ≥1 active assignment | state gauge | LAST |
| `PACKAGE_DATA_USED` | window usage per sim (existing dataset) attributed to the sim's active package | windowed join | SUM |
| `PACKAGES_EXPIRING` (optional) | active assignments with `end_date` within next 7d | state gauge | LAST |
| `QUEUED_PACKAGES` (optional) | assignments not yet started (queue depth) | state gauge | LAST |

## UI mapping

- MRR / ARPU / Customers-on-plan cards → `GetKpis(keys=MRR,ARPU,CUSTOMERS_ON_PLAN, op=LAST)`
- Top plan by revenue → `GetKpiBreakdown(key=PACKAGE_REVENUE, by=package_id, top=1)` (name/price/validity joined client-side from data-plan)
- Package performance table → breakdowns of PACKAGE_SALES, PACKAGE_REVENUE, PACKAGE_DATA_USED by package_id + package attributes from data-plan directly (price/validity/status are registry data, not analytics)
- Top packages / revenue mix → PACKAGE_REVENUE breakdown
- Time selector (Today/7d/month/quarter) → span=daily + from/to range on timeseries; quarterly span later (config list)

## Caveats (agreed-to-be-acceptable unless flagged)

1. **Revenue is a sales proxy** (assignments × price), not settled payments — real revenue KPIs come later from billing/payments sources.
2. `PACKAGE_SALES` counts assignments by `start_date` falling in a window; a backdated assignment created later is attributed to its start window only if that window recomputes (dirty path) — otherwise missed.
3. MRR normalization: `amount × 30d/duration` (UI already labels it "estimated").
4. Fan-out grows: networks → all sims → sim-packages per sim. Fine at current scale; the bulk-endpoint question applies if sims grow large (a `/v1/sim-packages?network_id=` bulk endpoint would collapse it).

## Implementation checklist

1. Add 3 pulls (data-plan spec file is new; extend subscriber spec).
2. Algos: `package_sales@v1`, `package_revenue@v1`, `mrr@v1`, `arpu@v1`, `customers_on_plan@v1`, `active_plans@v1`, `package_data_used@v1` + shared package-network matcher.
3. KPI spec yamls ×7 (analysis + aggregator copies).
4. No engine/aggregator/gateway changes expected — multi-key scopes and breakdowns already work.
