# Performance Reports Plan

Goal: resource-detail performance tables (first use case: Package
performance) composed from existing KPIs + entity attributes, defined
declaratively so new resource reports (sites, nodes, subscribers) are a YAML
file, not code.

## Where it lives

**A `performance` module inside the aggregator service** (new RPC on the
existing read API), NOT a fourth pipeline service:

- Aggregator is already the serving layer with the KPI read path; a report is
  a read-time composition (no new state, no writers, no msgbus).
- Entity attributes come from `raw_records` state-as-of (packages, networks,
  sims are already ingested) — aggregator gains a read-only raw-state reader
  (read-only zone crossing, same as analysis has for the raw zone).
- If report generation later grows heavy concerns (PDF/CSV exports,
  scheduled emailed reports), the module lifts out into its own service with
  the same proto — the seam is kept clean (`pkg/performance` package with its
  own interface).

## Declarative report specs (`configs/reports/*.yaml`)

```yaml
report: package_performance
title: Package performance
resource:
  dataset: dataplan.package.getAll     # entity rows = state-as-of latest
  entity_key: package_id
  network_match: network_id            # ""/equal rule (org-level packages)
  attributes:                          # passthrough columns from raw state
    - { name: name,     field: name }
    - { name: price,    field: amount,   format: money }
    - { name: validity, field: duration, format: days }
    - { name: active,   field: active }
row_scope: package_id                  # KPI scope dimension keyed to entity
columns:                               # KPI columns: key + op over the range
  - { name: sold,      kpi: PACKAGE_SALES,     op: SUM }
  - { name: revenue,   kpi: PACKAGE_REVENUE,   op: SUM, format: money }
  - { name: data_used, kpi: PACKAGE_DATA_USED, op: SUM, format: bytes }
status:                                # first matching rule wins
  - { when: "active == false",             label: Inactive }
  - { when: "sold == 0",                   label: No sales }
  - { when: "sold < 25",                   label: Low sales }
  - { default: true,                       label: Active }
sort: { column: revenue, desc: true }
```

Rules engine: minimal comparisons (`==`, `<`, `>`, on columns/attributes) —
implemented as a tiny evaluator, NOT a general expression language. Anything
fancier becomes a registered Go "status resolver" referenced by name (same
philosophy as KPI algos: config wires, code computes).

## Serving — latest available KPI values (decision)

No `from`/`to` params. Report cells read the LATEST available rollup row per
(kpi, scope, span, op) — the same `Latest()` read the `/kpis/values`
endpoint uses, trends included. `span` picks which bucket's latest state:
`daily` = today so far, `weekly` = this ISO week, `monthly` = this calendar
month (each `is_partial`-flagged while current). Freshness = the newest
window the pipeline processed; every cell carries `computed_at`.

```
GetPerformanceReport(report, span, scope_filters, top)
GET /v1/analytics/reports/package_performance?span=weekly&network_id=
```

Response: report metadata + rows, each row = entity id + attributes + KPI
cells (value/formatted/trend/is_partial/computed_at) + derived status label.
`ListReports()` for discovery.

Execution per request (all local reads, no new query machinery):
1. Entity rows: raw-state read of `resource.dataset` (filtered by
   network_match rule when network_id given).
2. KPI cells: per column, latest kpi_rollups rows for (kpi, span, op),
   grouped by `row_scope`.
3. Join on entity key; entities with no KPI rows show zeros (a package with
   no sales still lists); KPI scopes referencing deleted entities are
   dropped (tombstoned in the change-log).
4. Apply status rules, sort, top-N.

**Phase 2 (committed, not now): trailing ranges anchored to request time** —
[now − span, now] computed from base kpi_windows components, giving
"last 7 days ending right now" semantics instead of calendar buckets
(re-anchors on every call). Also elapsed-mode trends ("this week so far vs
same elapsed portion of last week").

## Why compose at read time (not precompute report rows)

Reports are cheap joins over already-materialized rollups; precomputing
would re-introduce a rollup-of-rollups pipeline stage for no freshness gain.
If a report ever gets hot/expensive, add caching at the gateway, not a new
pipeline stage.

## Phases

1. `pkg/performance`: report spec loader+validation (columns reference
   existing KPIs/ops — checked at startup), raw-state reader, composer,
   status rules; proto additions (`ListReports`, `GetPerformanceReport`);
   gateway route. Ship `package_performance.yaml`.
2. `site_performance.yaml` (SITES_ONLINE/DEGRADED + usage by site — needs
   site-scoped KPIs first) and `network_performance.yaml` — YAML-only once
   their KPIs exist.
3. CSV export of any report (gateway serialization of the same response).

## Notes / limits

- Attribute values are "current state" while KPI cells cover the selected
  range — a renamed package shows its current name against historical sales
  (industry-standard behavior for such tables).
- Trends per cell come free from kpi_rollups only when the range is a single
  span bucket; multi-bucket ranges return totals without trend (UI shows
  trend on the single-bucket views like "Today"/"This month").
- The status column in the mock ("Testing") implies a package lifecycle
  attribute that data-plan doesn't expose today; until it does, status is
  derived purely from `active` + sales thresholds.
