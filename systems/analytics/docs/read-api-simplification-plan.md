# Analytics read API — simplification plan

**Status:** Phases 0–3 implemented (2026-08-11); Phase 4 (console-bff + lab cutover,
legacy route deletion) pending · **Author:** Salman + Claude session

> Implementation notes vs the original text: the legacy endpoints became
> adapters over the Query planner in Phase 3 (with Phase 2 they were left
> untouched); rollup storage now holds ONE components row per scope/span
> (`op = 'VAL'`, with a `last` column) and a boot backfill re-materializes
> history from kpi_windows, so no manual migration is needed; `rollup_ops`
> was removed from the KPI specs; trend is computed at read time only.
**Scope:** the read/query surface of `systems/analytics` (api-gateway + aggregator) and the
KPI spec fields that drive it. The write pipeline (ingest → analysis → rollup) is sound and
stays as-is; this plan changes how the data is *asked for*, not how it is computed.

---

## 1. Diagnosis — why the current API feels complicated

A caller of `/v1/analytics/kpis/*` today juggles four knobs — `span`, `op`, `scope`,
`group_by` — across three endpoints (`values`, `timeseries`, `breakdown`), and the valid
combinations depend on rules that live in the server's head, not in the request:

| Knob | What the caller must know today |
|---|---|
| `span` | 6 tokens covering **three different concepts at once**: bucket size (daily/weekly/monthly), time range (last_24h/7d/30d), and calendar-vs-rolling alignment — which also selects between two storage paths (kpi_rollups vs kpi_windows). |
| `op` | 7 values (SUM/AVG/MIN/MAX/COUNT/LAST/DELTA), a **per-KPI whitelist** (`rollup_ops`), a per-KPI default rule (LAST if allowed, else first listed), and the fact that LAST/DELTA don't compose with folding. To ask "how much data did network X use", the caller must know DATA_USAGE is a flow (SUM) while MRR is a gauge (LAST). |
| `scope` | A filter — which now also implies the fold grain, but only for component ops, with a single-row passthrough rule, and not for LAST/DELTA. |
| `group_by` | Explicit fold grain, overrides the implicit one, component ops only. |
| endpoints | `breakdown` is really `group_by` + sort + top-N; `values` vs `timeseries` is really "one point vs many points". |

None of these knobs is wrong individually — each was added for a real need. The complexity
is *emergent*: the knobs are not orthogonal, and their interactions produce special cases
(`LAST` + filter, folded rows without trend, rolling vs calendar trend semantics) that a
caller cannot predict from the request shape.

**Root causes, in order of importance:**

1. **`op` leaks storage internals.** Whether a KPI is summed or read-latest is a property
   of the KPI, not of the question. Every mainstream analytics system (Prometheus
   counter/gauge, Cube.dev measure types, Looker measure aggregation, Amplitude metric
   definitions) defines the aggregation **in the metric definition**, server-side. Ukama
   already has the right place for it: the KPI spec.
2. **`span` overloads range + granularity + alignment.** Industry convention separates a
   *time range* (`from/to` or a relative token) from a *granularity* (bucket size for a
   series). "Rolling vs calendar" is then just what the range token means — not an API mode
   that flips storage paths and trend semantics.
3. **`LAST` and `DELTA` are compensations for a design we just removed.** DELTA existed to
   derive span usage from cumulative snapshots — DATA_USAGE now stores increments, so DELTA
   has no remaining user. LAST exists because gauges (MRR, ACTIVE_CUSTOMERS, uptime) answer
   "current value" with the latest window — that's a *KPI kind*, not a caller decision.
4. **Two grouping triggers.** Implicit (filter grain) and explicit (`group_by`) grouping
   with an override rule is one mechanism too many, forced by back-compat with per-row reads.
5. **Housekeeping debt** amplifies everything: KPI specs are duplicated in
   `analysis/configs/kpis/` and `aggregator/configs/kpis/` with nothing enforcing they
   match, and nothing cross-checks KPI input dataset keys against source specs (the
   PACKAGE_DATA_USED silent-death defect).

---

## 2. Target model — one question shape

### 2.1 KPI spec: declare the kind, drop the op menu

Each KPI declares what it *is*; the server derives how to aggregate it. Two kinds cover
every existing KPI:

```yaml
# flow: an amount that accrues over time (bytes, sales, cents).
#   over time  -> SUM        across scopes -> SUM
kpi: DATA_USAGE
kind: flow

# gauge: a level that exists at a moment (customers, MRR, uptime ratio).
#   over time  -> latest complete bucket (value), AVG for a series point
#   across scopes -> declared: sum (additive: customers, MRR) | avg (ratios: uptime, weighted by components)
kpi: MRR
kind: gauge
scope_agg: sum

kpi: NETWORK_UPTIME
kind: gauge
scope_agg: avg      # weighted: Σsum/Σcount from components — already exact today
```

`rollup_ops` disappears from the public contract (the engine can still materialize
whatever it likes internally). An optional `agg=` query param remains as an expert
override (`sum|avg|min|max|count`), defaulting from the spec — 99% of callers never
send it. This is exactly the Prometheus counter/gauge and Cube measure-type pattern.

Everything is computable because `kpi_windows` and (since this series) `kpi_rollups`
carry `Sum/Count/Min/Max` components — the spec change is metadata only.

### 2.2 One endpoint, five orthogonal parameters

```
GET /v1/analytics/query
    kpis=DATA_USAGE,REVENUE            # required; csv
    network_id=X&site_id=Y&iccid=Z     # filters — always fold to the answer's grain
    group_by=package_id                # extra dimensions; default none
    range=this_month                   # today|this_week|this_month|last_24h|last_7d|last_30d|from=..&to=..
    granularity=total                  # total|day|week|month  (total = one point per row)
    sort=-value&limit=10               # optional top-N
```

**Invariants — the whole contract in four sentences:**

- *Filters always aggregate.* A filtered read answers at exactly the grain of
  `filter keys + group_by keys`. No implicit/explicit distinction, no passthrough rule,
  no per-row legacy mode.
- *No `group_by` → one row per KPI* (the total for whatever the filter matches).
  `group_by` adds one row per distinct value combination. `breakdown` = `group_by` +
  `sort` + `limit` — not a separate endpoint.
- *`granularity=total` returns one point per row; any other granularity returns a series
  of buckets per row.* `values` vs `timeseries` collapses into this one switch.
- *Trend is always the same query evaluated over the immediately preceding range, at the
  same grain.* One definition, both storage paths, folded or not.

One response shape for everything:

```json
{ "rows": [
    { "kpi": "DATA_USAGE",
      "dims": { "network_id": "X", "package_id": "p1" },
      "points": [ { "t": "2026-08-01", "value": 1.2e9,
                    "trend": { "direction": "up", "change_pct": 12.5, "prev_value": 1.05e9 } } ],
      "unit": "bytes", "symbol": "B", "is_partial": true } ] }
```

### 2.3 What callers write, before vs after

| Question | Today | After |
|---|---|---|
| Usage of network X this month | `values?keys=DATA_USAGE&span=monthly&op=SUM&network_id=X` (must know op=SUM) | `query?kpis=DATA_USAGE&network_id=X&range=this_month` |
| Top-5 packages by usage within X | `breakdown?key=DATA_USAGE&by=package_id&span=monthly&op=SUM&network_id=X&top=5` | `query?kpis=DATA_USAGE&network_id=X&group_by=package_id&sort=-value&limit=5&range=this_month` |
| Daily usage chart, 30 days | `timeseries?keys=DATA_USAGE&span=daily&op=SUM&from=..&to=..&network_id=X` | `query?kpis=DATA_USAGE&network_id=X&range=last_30d&granularity=day` |
| Current MRR per network | `values?keys=MRR&span=daily&op=LAST` (must know LAST) | `query?kpis=MRR&group_by=network_id` |
| Console tile row (6 KPIs) | one call, but per-KPI op defaults differ silently | one call, kinds resolve per KPI — same URL shape as every other question |

### 2.4 Storage is untouched; the planner picks the path

`kpi_windows` (5-min, components) stays the source of truth; `kpi_rollups`
(daily/weekly/monthly materialization) becomes an internal *cache*, not an API concept.
A small query planner chooses:

- `granularity=day|week|month` → fold matching `kpi_rollups` rows (cheap, pre-materialized).
- `granularity=total` with a calendar range → fold the range's rollup rows.
- rolling ranges (`last_24h/7d/30d`) → fold `kpi_windows` components (today's rolling.go,
  generalized).

Both paths already share the same component-fold math; the plan merely routes to them from
one request shape instead of exposing the fork as `span` tokens.

---

## 3. Migration plan — five phases, each shippable alone

**Phase 0 — metadata (no behavior change).**
Add `kind` and `scope_agg` to every KPI spec (14 files); keep `rollup_ops` accepted and
ignored-if-`kind`-present. Extend `ListKpis` to return them so the console can start
rendering pickers from the registry instead of hardcoding.

**Phase 1 — kill the spec duplication (pipeline hygiene, enables everything else).**
One shared `systems/analytics/configs/kpis/` consumed by analysis *and* aggregator
(mount/copy in both images; loader unchanged). Add the startup cross-check of KPI input
dataset keys against loaded source specs — the lint that would have made the
PACKAGE_DATA_USED silent failure a boot error. Same for report specs.

**Phase 2 — implement `Query` (new gRPC RPC + `/v1/analytics/query`).**
A thin planner over the existing repos: resolve kinds → build filter+group grain → pick
windows vs rollups by range/granularity → component fold (reuse `grouping.go` /
`foldWindowAggs`) → trend = same plan over previous range → sort/limit. The three old
endpoints become 20-line adapters that translate their params into a Query call —
their behavior is then *defined* by the new model, and the special-case code
(`resolveOp` defaults, implicit-vs-explicit grouping, per-row passthrough, breakdown
folding) is deleted rather than maintained in parallel.

**Phase 3 — retire LAST/DELTA from the contract.**
DELTA: no spec uses it after the DATA_USAGE cutover — delete the op. LAST: served by
`kind: gauge` semantics; remove from specs and from the public param space (internally,
"latest complete bucket" remains as the gauge time-agg). `rollup_ops` fields deleted from
specs; the rollup engine materializes components only (one row per scope/span instead of
one per op — a storage *reduction*).

**Phase 4 — cut consumers over, then delete.**
console-bff and lab checks move to `/query` (mechanical: the table in §2.3 is the
translation guide). One release later, delete the old handlers and the adapter layer.
The gateway keeps returning 400 with a pointer to `/query` for the old routes for one
more release.

**Phase 5 — future-ready surface (as needed, not speculative).**
The registry (`GET /kpis`) already self-describes keys/units/scopes; with kind +
dimensions it becomes the single source the console uses to build any KPI screen
generically — adding a KPI stays a spec-file-only change end to end. Natural extensions
that fit this shape without new concepts: cursor pagination on `rows`, `compare=previous`
to return both periods explicitly, saved queries (a report column becomes
`{kpi, filter?, group_by?}` — the performance-report composer re-expressed on Query),
and per-org rate limits keyed on planner cost (rows × range / granularity).

---

## 4. What deliberately does NOT change

- The window grid, ledger, change-log ingest, algo Result components, and the
  weighted-fold math — these are the system's strengths; every phase above reuses them.
- The performance reports YAML model (entity rows + KPI columns + status rules) — only the
  column's `op` field goes away (derived from kind).
- DATA_USAGE semantics from the metrics cutover (increments, `state_prev`, filter
  validation with 400s) — the Query endpoint inherits them unchanged.

## 5. Risks and mitigations

- *Gauge-across-scope semantics* (sum vs avg) must be right per KPI — reviewed once, in
  Phase 0, in the spec files, where it is visible; today the same decision is smeared
  across caller-chosen ops.
- *Old-endpoint behavior drift during Phase 2* — mitigated by making the old endpoints
  adapters over Query on day one, so there is exactly one semantics.
- *Console breakage* — the adapter phase means the console keeps working untouched until
  Phase 4; its migration is a URL rewrite, not a data-model change.

## 6. Effort estimate

| Phase | Size | Notes |
|---|---|---|
| 0 | S (½ day) | 14 spec files + ListKpis fields |
| 1 | S (½–1 day) | build/deploy wiring + startup lint |
| 2 | M (2–3 days) | planner + adapters + tests; heaviest reuse of existing fold code |
| 3 | S (1 day) | deletions + spec cleanup + rollup writes shrink |
| 4 | M (1–2 days) | console-bff + lab call-site rewrites |
| 5 | on demand | each item independently small |
